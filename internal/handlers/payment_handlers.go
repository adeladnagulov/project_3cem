package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"project_3sem/internal/middleware"
	"project_3sem/internal/repositories"
	"project_3sem/internal/responses"
	"project_3sem/internal/services"
	"strings"
	"time"

	"github.com/google/uuid"
)

var SitePrise = 199.00
var SiteCurrency = "RUB"

var YooKassaIPs = []string{
	"185.71.76.0/27",
	"185.71.77.0/27",
	"77.75.153.0/25",
	"77.75.156.11",
	"77.75.156.35",
	"77.75.154.128/25",
	"2a02:5180::/32",
}

type PaymentHandler struct {
	RepoPayments       repositories.RepoPayments
	RepoUsers          repositories.RepoUsers
	RepoOrdersPayments repositories.RepoOrdersPayments
	RepoOrders         repositories.RepoOrders
}

func NewPaymentHandler(repoPayments repositories.RepoPayments, repoUsers repositories.PgRepoUsers,
	repoOrdersPayments repositories.RepoOrdersPayments, repoOrders repositories.RepoOrders) *PaymentHandler {
	return &PaymentHandler{
		RepoPayments:       repoPayments,
		RepoUsers:          &repoUsers,
		RepoOrdersPayments: repoOrdersPayments,
		RepoOrders:         repoOrders,
	}
}

type YooKassaRequest struct {
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Capture      bool `json:"capture"`
	Confirmation struct {
		Type       string `json:"type"`
		Return_url string `json:"return_url"`
	} `json:"confirmation"`
	Description string `json:"description"`
}

type YooKassaResponse struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Amount struct {
		Value    string `json:"value"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Description  string `json:"description"`
	Confirmation struct {
		Confirmation_url string `json:"confirmation_url"`
	} `json:"confirmation"`
	Paid bool `json:"paid"`
}

func (h *PaymentHandler) CreatePayments(w http.ResponseWriter, r *http.Request) { //нужен рефакторинг
	log.Println("===start payment===")
	var req struct {
		//Value       float64 `json:"value"`
		//Currency    string  `json:"currency"`
		Return_url  string `json:"return_url"`
		Description string `json:"description"`
		Site_ID     string `json:"site_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("decode error" + err.Error())
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	yooReq := YooKassaRequest{
		Amount: struct {
			Value    string "json:\"value\""
			Currency string "json:\"currency\""
		}{
			Value:    fmt.Sprintf("%.2f", SitePrise), //цена фикс
			Currency: SiteCurrency,
		},
		Capture: true,
		Confirmation: struct {
			Type       string "json:\"type\""
			Return_url string "json:\"return_url\""
		}{
			Type:       "redirect",
			Return_url: req.Return_url,
		},
		Description: req.Description,
	}
	reqBody, err := json.Marshal(yooReq)
	if err != nil {
		http.Error(w, "Marhal error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	yooKassReq, err := http.NewRequest("POST", "https://api.yookassa.ru/v3/payments", bytes.NewBuffer(reqBody))
	if err != nil {
		http.Error(w, "Create request error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	shopId := os.Getenv("SHOP_ID")
	apiKey := os.Getenv("API_KEY")
	auth := base64.StdEncoding.EncodeToString([]byte(shopId + ":" + apiKey))
	yooKassReq.Header.Add("Authorization", "Basic "+auth)
	yooKassReq.Header.Add("Idempotence-Key", uuid.NewString())
	yooKassReq.Header.Add("Content-Type", "application/json")

	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(yooKassReq)
	if err != nil {
		http.Error(w, "Payment service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	log.Printf("yooKassa respone: status=%d, body=%s", resp.StatusCode, bodyBytes)
	//log.Println(shopId + " " + apiKey)

	yooResp := YooKassaResponse{}
	if err := json.Unmarshal(bodyBytes, &yooResp); err != nil {
		log.Println("decode error: " + err.Error())
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	userId := r.Context().Value(middleware.IdKey).(string)
	//log.Println("amount = " + yooResp.Amount.Value)
	if err = h.RepoPayments.SavePayment(
		yooResp.Id,
		yooResp.Status,
		yooResp.Amount.Value,
		yooResp.Amount.Currency,
		yooResp.Description,
		req.Site_ID,
		userId); err != nil {

		log.Println("SavePayment error " + err.Error())
		http.Error(w, "SavePayment error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"payment_id":       yooResp.Id,
		"confirmation_url": yooResp.Confirmation.Confirmation_url,
		"status":           yooResp.Status,
	}
	responses.SendJSONResp(w, response, http.StatusOK)
}

func (h *PaymentHandler) PaymentWebhook(w http.ResponseWriter, r *http.Request) {
	ip := getRealIP(r)
	log.Printf("Webhook received from IP: %s", ip)
	if !IsYooKassaIp(ip) {
		log.Printf("Invalid YooKassa IP: %s", ip)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req struct {
		Object struct {
			Id     string `json:"id"`
			Status string `json:"status"`
		} `json:"object"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Decode error: %s", err.Error())
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Сначала проверяем, является ли это оплата корзины
	orderPayment, err := h.RepoOrdersPayments.GetByYookassaID(req.Object.Id)
	if err == nil {
		log.Printf("Processing order payment webhook: %s, status: %s", req.Object.Id, req.Object.Status)

		if err := h.RepoOrdersPayments.UpdateStatus(req.Object.Id, req.Object.Status); err != nil {
			log.Printf("Update order payment status error: %s", err.Error())
		}

		if req.Object.Status == "succeeded" {
			if err := h.RepoOrders.UpdateOrderStatus(orderPayment.OrderID.String(), "paid"); err != nil {
				log.Printf("Update order status error: %s", err.Error())
			}
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	// Если не найден как оплата заказа, проверяем как оплату сайта (поддомена)
	payment, err := h.RepoPayments.UpdateStatus(req.Object.Id, req.Object.Status)
	if err != nil {
		log.Printf("Update site payment status error: %s", err.Error())
		http.Error(w, "Payment not found", http.StatusNotFound)
		return
	}

	user, err := h.RepoUsers.GetUserByID(payment.User_id)
	if err != nil {
		log.Printf("Get user error: %s", err.Error())
		w.WriteHeader(http.StatusOK)
		return
	}

	message := ""
	if req.Object.Status == "succeeded" {
		message = fmt.Sprintf("Платёж %s на сумму %.2f %s прошел успешно", payment.Yookassa_payment_id, payment.Amount, payment.Currency)
	} else {
		message = fmt.Sprintf("Платёж %s отменён", payment.Yookassa_payment_id)
	}

	if err := services.CreateEmailService().SendCodeToEmail(user.Email, message); err != nil {
		log.Printf("Send email error: %s", err.Error())
	}

	log.Printf("Notification sent: %s to email: %s", message, user.Email)
	w.WriteHeader(http.StatusOK)
}

func getRealIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func IsYooKassaIp(ipstr string) bool {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return false
	}
	for _, yooKAssaip := range YooKassaIPs {
		if strings.Contains(yooKAssaip, "/") {
			_, ipNet, err := net.ParseCIDR(yooKAssaip)
			if err != nil {
				continue
			}
			if ipNet.Contains(ip) {
				return true
			}
		} else {
			certainIp := net.ParseIP(yooKAssaip)
			if certainIp != nil && certainIp.Equal(ip) {
				return true
			}
		}
	}
	return false
}

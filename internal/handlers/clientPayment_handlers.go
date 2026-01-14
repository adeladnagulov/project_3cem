package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"project_3sem/internal/middleware"
	"project_3sem/internal/repositories"
	"project_3sem/internal/responses"
	"project_3sem/internal/services"
)

type OrderHandler struct {
	RepoSites          repositories.PgRepoSites
	RepoOrders         repositories.RepoOrders
	RepoOrdersPayments repositories.RepoOrdersPayments
}

func NewOrderHandler(repoSites repositories.PgRepoSites, repoOrders repositories.RepoOrders,
	repoOrdersPayments repositories.RepoOrdersPayments) *OrderHandler {
	return &OrderHandler{
		RepoSites:          repoSites,
		RepoOrders:         repoOrders,
		RepoOrdersPayments: repoOrdersPayments,
	}
}

func (h *OrderHandler) BasketPayment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Items       map[string]interface{} `json:"items"`
		TotalAmount float64                `json:"total_amount"`
		Return_url  string                 `json:"return_url"`
		Description string                 `json:"description"`
		Email       string                 `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("bad request, error :%s", err.Error())
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	subdomain := r.Context().Value(middleware.SubdomainKey).(string)
	site := h.RepoSites.GetPublishBySubdomain(subdomain)
	orderId, err := h.RepoOrders.AddOrder(site.ID, req.Items, req.TotalAmount)
	if err != nil {
		log.Println("add order error: " + err.Error())
		http.Error(w, "add new order error", http.StatusInternalServerError)
		return
	}

	yooResp, err := services.CreatePayment(w, req.TotalAmount, SiteCurrency, req.Return_url, req.Description)
	if err != nil {
		log.Println("create payment error: " + err.Error())
		http.Error(w, "dont create payment", http.StatusInternalServerError)
		return
	}

	err = h.RepoOrdersPayments.SaveOrderPayment(yooResp.Id, yooResp.Status, yooResp.Amount.Value,
		yooResp.Amount.Currency, yooResp.Description, site.ID, orderId)
	if err != nil {
		log.Println("SaveOrderPayment error: " + err.Error())
		http.Error(w, "Dont save order payment", http.StatusInternalServerError)
		return
	}

	err = services.CreateEmailService().SendCodeToEmail(req.Email, "Ваш номер заказа: "+orderId)
	if err != nil {
		log.Println("SendCodeToEmail err0r: " + err.Error())
		http.Error(w, "Dont send orderId", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"order_id":         orderId,
		"confirmation_url": yooResp.Confirmation.Confirmation_url,
		"amount":           yooResp.Amount.Value,
		"currency":         yooResp.Amount.Currency,
		"status":           yooResp.Status,
		"description":      yooResp.Description,
	}
	responses.SendJSONResp(w, resp, http.StatusOK)
}

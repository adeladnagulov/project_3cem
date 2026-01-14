package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

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

func CreatePayment(w http.ResponseWriter, value float64, currency, return_url, description string) (*YooKassaResponse, error) {
	yooReq := YooKassaRequest{
		Amount: struct {
			Value    string "json:\"value\""
			Currency string "json:\"currency\""
		}{
			Value:    fmt.Sprintf("%.2f", value),
			Currency: currency,
		},
		Capture: true,
		Confirmation: struct {
			Type       string "json:\"type\""
			Return_url string "json:\"return_url\""
		}{
			Type:       "redirect",
			Return_url: return_url,
		},
		Description: description,
	}
	reqBody, err := json.Marshal(yooReq)
	if err != nil {
		http.Error(w, "Marhal error: "+err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	yooKassReq, err := http.NewRequest("POST", "https://api.yookassa.ru/v3/payments", bytes.NewBuffer(reqBody))
	if err != nil {
		http.Error(w, "Create request error: "+err.Error(), http.StatusInternalServerError)
		return nil, err
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
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	log.Printf("yooKassa respone: status=%d, body=%s", resp.StatusCode, bodyBytes)
	//log.Println(shopId + " " + apiKey)

	yooResp := YooKassaResponse{}
	if err := json.Unmarshal(bodyBytes, &yooResp); err != nil {
		log.Println("decode error: " + err.Error())
		http.Error(w, "bad request: "+err.Error(), http.StatusBadRequest)
		return nil, err
	}
	return &yooResp, nil
}

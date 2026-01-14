package models

import "github.com/google/uuid"

type OrderPayment struct {
	ID                uuid.UUID
	OrderID           uuid.UUID
	SiteID            uuid.UUID
	YookassaPaymentID string
	Status            string
	Amount            float64
	Currency          string
	Description       string
}

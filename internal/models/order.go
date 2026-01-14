package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID          uuid.UUID `json:"id"`
	SiteID      uuid.UUID `json:"site_id"`
	Items       []byte    `json:"items"` // Хранит JSONB данные из базы
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

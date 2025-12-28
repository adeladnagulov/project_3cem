package models

import (
	"database/sql"
	"time"
)

type Site struct {
	ID         string                 `json:"id"`
	UserID     string                 `json:"user_id"`
	Subdomain  string                 `json:"subdomain"`
	Pattern    string                 `json:"pattern"`
	Config     map[string]interface{} `json:"config"`
	Status     string                 `json:"status"`
	CreatedAt  time.Time              `json:"created_at"`
	PublishdAt sql.NullTime           `json:"pablish_at"`
}

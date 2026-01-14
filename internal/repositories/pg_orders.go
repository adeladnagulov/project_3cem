package repositories

import (
	"database/sql"
	"encoding/json"
)

type RepoOrders interface {
	AddOrder(siteId string, items map[string]interface{}, totalAmount float64) (string, error)
}

type PgRepoOrders struct {
	db *sql.DB
}

func NewPgRepoOrders(db *sql.DB) *PgRepoOrders {
	return &PgRepoOrders{
		db: db,
	}
}

func (r *PgRepoOrders) AddOrder(siteId string, items map[string]interface{}, totalAmount float64) (string, error) {
	id := ""
	itemsJSON, err := json.Marshal(items)
	if err != nil {
		return "", err
	}
	err = r.db.QueryRow(`
	INSERT INTO orders (site_id, items, total_amount, status)
	VALUES ($1, $2, $3, $4)
	RETURNING id`, siteId, itemsJSON, totalAmount, "awaiting payment").Scan(&id)
	return id, err
}

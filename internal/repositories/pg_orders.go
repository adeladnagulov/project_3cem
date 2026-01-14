package repositories

import (
	"database/sql"
	"encoding/json"
	"project_3sem/internal/models"
)

type RepoOrders interface {
	AddOrder(siteId string, items map[string]interface{}, totalAmount float64) (string, error)
	GetOrdersBySiteID(siteID string) ([]models.Order, error)
	GetPaidOrdersAmountBySiteID(siteID string) (float64, error)
	UpdateOrderStatus(orderID string, status string) error
}

type PgRepoOrders struct {
	db *sql.DB
}

func NewPgRepoOrders(db *sql.DB) *PgRepoOrders {
	return &PgRepoOrders{
		db: db,
	}
}

func (r *PgRepoOrders) UpdateOrderStatus(orderID string, status string) error {
	_, err := r.db.Exec(`
        UPDATE orders	
        SET status = $1
        WHERE id = $2::uuid`, status, orderID)
	return err
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

func (r *PgRepoOrders) GetOrdersBySiteID(siteID string) ([]models.Order, error) {
	var orders []models.Order
	rows, err := r.db.Query(`
        SELECT id, site_id, items, total_amount, status, created_at 
        FROM orders 
        WHERE site_id = $1 
        ORDER BY created_at DESC`, siteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		var items []byte
		err := rows.Scan(&order.ID, &order.SiteID, &items, &order.TotalAmount, &order.Status, &order.CreatedAt)
		if err != nil {
			return nil, err
		}
		order.Items = items
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *PgRepoOrders) GetPaidOrdersAmountBySiteID(siteID string) (float64, error) {
	var totalAmount float64
	err := r.db.QueryRow(`
        SELECT COALESCE(SUM(total_amount), 0) 
        FROM orders 
        WHERE site_id = $1 AND status = 'paid'`, siteID).Scan(&totalAmount)

	return totalAmount, err
}

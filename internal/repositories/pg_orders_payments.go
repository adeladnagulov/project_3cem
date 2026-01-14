package repositories

import (
	"database/sql"
	"errors"
	"project_3sem/internal/models"
	"strconv"
)

type RepoOrdersPayments interface {
	SaveOrderPayment(yookassaID, status, amountStr, currency, description, site_id, order_id string) error
	UpdateStatus(yookassaID, status string) error
	GetByYookassaID(yookassaID string) (*models.OrderPayment, error)
	UpdateOrderPaymentStatus(yookassaID, status string) error
}

type PgRepoOrdersPayments struct {
	db *sql.DB
}

func NewPgRepoOrdersPayments(db *sql.DB) *PgRepoOrdersPayments {
	return &PgRepoOrdersPayments{
		db: db,
	}
}

func (r *PgRepoOrdersPayments) SaveOrderPayment(yookassaID, status, amountStr, currency, description, site_id, order_id string) error {
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return errors.New("invalid amount")
	}
	_, err = r.db.Exec(`
	INSERT INTO orders_payments (yookassa_payment_id, status, amount, currency, description, sites_id, order_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`, yookassaID, status, amount, currency, description, site_id, order_id)
	return err
}

func (r *PgRepoOrdersPayments) UpdateStatus(yookassaID, status string) error {
	_, err := r.db.Exec(`
        UPDATE orders_payments
        SET status = $1
        WHERE yookassa_payment_id = $2`, status, yookassaID)
	return err
}

func (r *PgRepoOrdersPayments) GetByYookassaID(yookassaID string) (*models.OrderPayment, error) {
	payment := &models.OrderPayment{}
	err := r.db.QueryRow(`
        SELECT id, order_id, sites_id, yookassa_payment_id, status, amount, currency, description 
        FROM orders_payments 
        WHERE yookassa_payment_id = $1`, yookassaID).Scan(
		&payment.ID, &payment.OrderID, &payment.SiteID, &payment.YookassaPaymentID,
		&payment.Status, &payment.Amount, &payment.Currency, &payment.Description)

	if err == sql.ErrNoRows {
		return nil, errors.New("payment not found")
	}
	return payment, err
}

func (r *PgRepoOrdersPayments) UpdateOrderPaymentStatus(yookassaID, status string) error {
	_, err := r.db.Exec(`
        UPDATE orders_payments
        SET status = $1
        WHERE yookassa_payment_id = $2`, status, yookassaID)
	return err
}

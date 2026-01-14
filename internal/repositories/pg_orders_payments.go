package repositories

import (
	"database/sql"
	"errors"
	"strconv"
)

type RepoOrdersPayments interface {
	SaveOrderPayment(yookassaID, status, amountStr, currency, description, site_id string) error
}

type PgRepoOrdersPayments struct {
	db *sql.DB
}

func NewPgRepoOrdersPayments(db *sql.DB) *PgRepoPayments {
	return &PgRepoPayments{
		db: db,
	}
}

func (r *PgRepoPayments) SaveOrderPayment(yookassaID, status, amountStr, currency, description, site_id string) error {
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return errors.New("invalid amount")
	}
	_, err = r.db.Exec(`
	INSERT INTO orders_payments (yookassa_payment_id, status, amount, currency, description, sites_id)
	VALUES ($1, $2, $3, $4, $5, $6)`, yookassaID, status, amount, currency, description, site_id)
	return err
}

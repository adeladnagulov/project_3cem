package repositories

import (
	"database/sql"
	"errors"
	"strconv"
)

type PgPayments interface {
	SavePayment(yookassaID, status, amountStr, currency, description, site_id, user_id string) error
}

type PgRepoPayments struct {
	db *sql.DB
}

func NewPgRepoPayments(db *sql.DB) *PgRepoPayments {
	return &PgRepoPayments{
		db: db,
	}
}

func (r *PgRepoPayments) SavePayment(yookassaID, status, amountStr, currency, description, site_id, user_id string) error {
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return errors.New("invalid amount")
	}
	_, err = r.db.Exec(`
	INSERT INTO payments (yookassa_payment_id, status, amount, currency, description, user_id, sites_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7)`, yookassaID, status, amount, currency, description, user_id, site_id)
	return err
}

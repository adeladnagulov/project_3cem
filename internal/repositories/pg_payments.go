package repositories

import (
	"database/sql"
	"errors"
	"project_3sem/internal/models"
	"strconv"
)

type RepoPayments interface {
	SavePayment(yookassaID, status, amountStr, currency, description, site_id, user_id string) error
	UpdateStatus(yookassaID, status string) (*models.Payment, error)
	CheckStatus(siteId string) (string, error)
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

func (r *PgRepoPayments) UpdateStatus(yookassaID, status string) (*models.Payment, error) {
	payment := models.Payment{}
	err := r.db.QueryRow(`
	UPDATE payments
	SET status = $1
	WHERE yookassa_payment_id = $2
	RETURNING id, user_id, sites_id, yookassa_payment_id, status, amount, currency, description
	`, status, yookassaID).Scan(&payment.Id, &payment.User_id, &payment.Site_id, &payment.Yookassa_payment_id, &payment.Status,
		&payment.Amount, &payment.Currency, &payment.Description)
	return &payment, err
}

func (r *PgRepoPayments) CheckStatus(siteId string) (string, error) {
	status := ""
	err := r.db.QueryRow(`
	SELECT status
	FROM payments WHERE sites_id = $1`, siteId).Scan(&status)
	return status, err
}

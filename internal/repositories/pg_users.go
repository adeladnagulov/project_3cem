package repositories

import (
	"database/sql"
	"project_3sem/internal/models"
)

type PgRepoUsers struct {
	db *sql.DB
}

func NewPgRepoUsers(db *sql.DB) *PgRepoUsers {
	return &PgRepoUsers{
		db: db,
	}
}

func (r *PgRepoUsers) Authorization(email string) (*models.User, error) {
	u := models.User{}

	err := r.db.QueryRow(`
	INSERT INTO users (email) VALUES ($1) 
	ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email
	RETURNING id, email`, email).Scan(&u.ID, &u.Email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *PgRepoUsers) GetUserByID(id string) (*models.User, error) {
	u := models.User{}

	err := r.db.QueryRow(`
	SELECT id, email FROM users WHERE id = $1`, id).Scan(&u.ID, &u.Email)

	return &u, err
}

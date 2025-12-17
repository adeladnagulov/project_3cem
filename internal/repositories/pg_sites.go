package repositories

import (
	"database/sql"
	"errors"
	"log"
	"project_3sem/internal/models"
)

type PgRepoSites struct {
	db *sql.DB
}

func NewPgRepoSites(db *sql.DB) *PgRepoSites {
	return &PgRepoSites{
		db: db,
	}
}

func (r *PgRepoSites) AddDrawtToRepo(subdomain, pattern, userId string, config map[string]interface{}) (*models.Site, error) {
	if pattern == "" {
		return nil, errors.New("invalid pattern")
	}
	if userId == "" {
		return nil, errors.New("invalid userId")
	}

	if subdomain == "" {
		return nil, errors.New("invalid subdomain")
	} else if !r.CheckSubdomainInFree(subdomain) {
		return nil, errors.New("subdomain already taken")
	}
	site := models.Site{}

	err := r.db.QueryRow(`
	INSERT INTO users (user_id, subdomain, pattern, config)
	VALUES ($1, $2, $3, $4)
	RETURNING id, user_id, subdomain, pattern, config, status_site, created_at`, userId, subdomain, pattern, config,
	).Scan(&site.ID, &site.UserID, &site.Subdomain, &site.Pattern, &site.Config, &site.Status, &site.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &site, nil
}

func (r *PgRepoSites) CheckSubdomainInFree(subdomain string) bool {
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(
	SELECT 1 FROM sites
	WHERE subdomain = $1 AND status_site = 'published')`, subdomain).Scan(&exists)

	if err != nil {
		log.Printf("Error checking subdomain: %v", err)
		return false
	}
	return !exists
}

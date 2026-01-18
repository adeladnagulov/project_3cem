package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"project_3sem/internal/models"
	"time"
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

	cnfJson, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	site := models.Site{}
	var cnfBytes []byte
	err = r.db.QueryRow(`
	INSERT INTO sites (user_id, subdomain, pattern, config)
	VALUES ($1, $2, $3, $4)
	RETURNING id, user_id, subdomain, pattern, config, status_site, created_at`, userId, subdomain, pattern, cnfJson,
	).Scan(&site.ID, &site.UserID, &site.Subdomain, &site.Pattern, &cnfBytes, &site.Status, &site.CreatedAt)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(cnfBytes, &site.Config); err != nil {
		return nil, err
	}

	log.Printf("===add drawt in PG===")
	return &site, nil
}

func (r *PgRepoSites) PublishSite(siteId string) (*models.Site, error) {
	site := models.Site{}
	var cnfBytes []byte
	err := r.db.QueryRow(`
	SELECT id, user_id, subdomain, pattern, config, status_site, created_at
    FROM sites WHERE id = $1`, siteId).Scan(&site.ID, &site.UserID, &site.Subdomain, &site.Pattern, &cnfBytes, &site.Status, &site.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("site not found")
	}
	if err != nil {
		log.Printf("error find site")
		return nil, err
	}

	if err = json.Unmarshal(cnfBytes, &site.Config); err != nil {
		return nil, err
	}

	if !r.CheckSubdomainInFree(site.Subdomain) {
		return nil, errors.New("subdomain already taken")
	}

	err = r.db.QueryRow(`
	UPDATE sites
	SET status_site = 'published', published_at = $1
	WHERE id = $2
	RETURNING status_site, published_at
	`, time.Now(), siteId).Scan(&site.Status, &site.PublishdAt)
	if err != nil {
		log.Printf("error update site")
		return nil, err
	}
	return &site, nil
}

func (r *PgRepoSites) GetPublishBySubdomain(subdomain string) *models.Site {
	site := models.Site{}
	var cnfBytes []byte
	err := r.db.QueryRow(`
	SELECT id, user_id, subdomain, pattern, config, status_site, created_at, published_at
    FROM sites WHERE subdomain = $1 AND status_site = 'published'`, subdomain).
		Scan(&site.ID, &site.UserID, &site.Subdomain, &site.Pattern, &cnfBytes, &site.Status, &site.CreatedAt, &site.PublishdAt)
	if err == sql.ErrNoRows {
		log.Println(err.Error())
		return nil
	}
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	if err = json.Unmarshal(cnfBytes, &site.Config); err != nil {
		log.Println(err.Error())
		return nil
	}
	return &site
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

func (r *PgRepoSites) GetUserSites(userId string) ([]models.Site, error) {
	var sites []models.Site
	rows, err := r.db.Query(`
        SELECT id, subdomain, pattern, config, status_site, created_at, published_at 
        FROM sites 
        WHERE user_id = $1 
        ORDER BY created_at DESC`, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var site models.Site
		var config []byte
		var publishedAt sql.NullTime

		err := rows.Scan(&site.ID, &site.Subdomain, &site.Pattern, &config, &site.Status, &site.CreatedAt, &publishedAt)
		if err != nil {
			return nil, err
		}

		site.PublishdAt = publishedAt

		if config != nil {
			json.Unmarshal(config, &site.Config)
		} else {
			site.Config = make(map[string]interface{})
		}

		sites = append(sites, site)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return sites, nil
}

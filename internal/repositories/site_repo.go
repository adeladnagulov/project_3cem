package repositories

import (
	"errors"
	"project_3sem/internal/models"
	"sync"
	"time"

	"github.com/google/uuid"
)

type RepoSite interface {
	AddDrawtToRepo(subdomain, pattern, userId string, config map[string]interface{}) (*models.Site, error)
}

type MemoryRepoSites struct {
	mu    sync.RWMutex
	Sites map[string]*models.Site
}

func NewMemoryRepoSites() *MemoryRepoSites {
	return &MemoryRepoSites{
		Sites: make(map[string]*models.Site),
	}
}

func (r *MemoryRepoSites) AddDrawtToRepo(subdomain, pattern, userId string, config map[string]interface{}) (*models.Site, error) {
	if pattern == "" {
		return nil, errors.New("invalid pattern")
	}
	if userId == "" {
		return nil, errors.New("invalid userId")
	}

	//добавить проверку поддомена

	site := &models.Site{
		ID:        uuid.NewString(),
		UserID:    userId,
		Subdomain: subdomain,
		Pattern:   pattern,
		Config:    config,
		Status:    "draft",
		CreatedAt: time.Now(),
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Sites[site.ID] = site
	return site, nil
}

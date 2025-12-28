package repositories

import (
	"database/sql"
	"errors"
	"project_3sem/internal/models"
	"sync"
	"time"

	"github.com/google/uuid"
)

type RepoSite interface {
	AddDrawtToRepo(subdomain, pattern, userId string, config map[string]interface{}) (*models.Site, error)
	PublishSite(siteId string) (*models.Site, error)
	GetPublishBySubdomain(subdomain string) *models.Site
	CheckSubdomainInFree(subdomain string) bool
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

	if subdomain == "" {
		return nil, errors.New("invalid subdomain")
	} else if !r.CheckSubdomainInFree(subdomain) {
		return nil, errors.New("subdomain already taken")
	}

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

func (r *MemoryRepoSites) PublishSite(siteId string) (*models.Site, error) {
	r.mu.RLock()
	site, ok := r.Sites[siteId]
	if !ok {
		return nil, errors.New("site not sound")
	}
	r.mu.RUnlock()
	if !r.CheckSubdomainInFree(site.Subdomain) {
		return nil, errors.New("subdomain already taken")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	site.Status = "published"
	site.PublishdAt = sql.NullTime{Time: time.Now(), Valid: true}
	return site, nil
}

func (r *MemoryRepoSites) GetPublishBySubdomain(subdomain string) *models.Site {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, s := range r.Sites {
		if s.Subdomain == subdomain && s.Status == "published" {
			rez := s
			return rez
		}
	}
	return nil
}

func (r *MemoryRepoSites) CheckSubdomainInFree(subdomain string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, s := range r.Sites {
		if s.Subdomain == subdomain && s.Status == "published" {
			return false
		}
	}
	return true
}

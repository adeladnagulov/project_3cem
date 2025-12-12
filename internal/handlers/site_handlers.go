package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"project_3sem/internal/middleware"
	"project_3sem/internal/repositories"
	"project_3sem/internal/responses"

	"github.com/gorilla/mux"
)

type SiteHandle struct {
	RepoSite repositories.RepoSite
}

func NewSiteHandler(repoSite repositories.RepoSite) *SiteHandle {
	return &SiteHandle{
		RepoSite: repoSite,
	}
}

func (h *SiteHandle) SaveDraft(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Subdomain string                 `json:"subdomain"`
		Pattern   string                 `json:"pattern"`
		Config    map[string]interface{} `json:"config"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("bad request, error :%s", err.Error())
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	userId := r.Context().Value(middleware.IdKey).(string)

	site, err := h.RepoSite.AddDrawtToRepo(req.Subdomain, req.Pattern, userId, req.Config)
	if err != nil {
		log.Printf("bad request, error :%s", err.Error())
		http.Error(w, "invalid params, err: "+err.Error(), http.StatusBadRequest)
		return
	}
	log.Printf("Save Drawt done")

	resp := map[string]interface{}{
		"id":        site.ID,
		"subdomain": site.Subdomain,
		"config":    site.Config,
		"status":    site.Status,
	}
	responses.SendJSONResp(w, resp, http.StatusCreated)
}

func (h *SiteHandle) Publish(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	siteId := vars["id"]
	log.Printf("site id: %s", siteId)

	site, err := h.RepoSite.PublishSite(siteId)
	if err != nil {
		log.Printf("publish error: %s", err.Error())
		http.Error(w, "cannot publish", http.StatusBadRequest)
		return
	}
	log.Printf("site status: %s", site.Status)
	resp := map[string]interface{}{
		"id":        site.ID,
		"subdomain": site.Subdomain,
		"status":    site.Status,
	}
	responses.SendJSONResp(w, resp, http.StatusOK)
}

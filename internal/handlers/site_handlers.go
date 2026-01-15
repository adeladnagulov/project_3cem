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
	RepoSite    repositories.RepoSite
	RepoPayment repositories.RepoPayments
}

func NewSiteHandler(repoSite repositories.RepoSite, repoPayment repositories.RepoPayments) *SiteHandle {
	return &SiteHandle{
		RepoSite:    repoSite,
		RepoPayment: repoPayment,
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
	paymentStatus, err := h.RepoPayment.CheckStatus(siteId)
	if err != nil || paymentStatus != "succeeded" {
		log.Printf("payment status = %s, error = %s", paymentStatus, err)
		http.Error(w, "access is restricted", http.StatusForbidden)
		return
	}

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

func (h *SiteHandle) RenderSite(w http.ResponseWriter, r *http.Request) {
	log.Printf("render start!")
	subI := r.Context().Value(middleware.SubdomainKey)
	if subI == nil {
		log.Printf("no subdomain provided")
		http.Error(w, "no subdomain provided", http.StatusBadRequest)
		return
	}
	subD := subI.(string)
	log.Printf("current subdomain: %s", subD)
	if subD == "" {
		log.Printf("empty subdomain")
		http.Error(w, "no subdomain provided", http.StatusBadRequest)
		return
	}

	site := h.RepoSite.GetPublishBySubdomain(subD)
	if site == nil {
		log.Printf("not found site")
		http.Error(w, "not found site", http.StatusNotFound)
		return
	}
	if site.Status != "published" {
		log.Println("site not published")
		http.Error(w, "site not published", http.StatusNotFound)
		return
	}

	resp := map[string]interface{}{
		"id":        site.ID,
		"subdomain": site.Subdomain,
		"pattern":   site.Pattern,
		"config":    site.Config,
	}
	log.Printf("render done!")
	responses.SendJSONResp(w, resp, http.StatusOK)
}

func (h *SiteHandle) GetUserSites(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.IdKey).(string)

	sites, err := h.RepoSite.GetUserSites(userId)
	if err != nil {
		log.Printf("Ошибка получения сайтов пользователя: %v", err)
		http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"sites": sites,
	}

	responses.SendJSONResp(w, resp, http.StatusOK)
}

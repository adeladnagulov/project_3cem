package handlers

import (
	"net/http"
	"project_3sem/internal/middleware"
	"project_3sem/internal/responses"
)

func (h *UserHandle) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(middleware.IdKey).(string)
	email := r.Context().Value(middleware.EmailKey).(string)

	resp := map[string]interface{}{
		"id":    id,
		"email": email,
	}
	responses.SendJSONResp(w, resp, http.StatusOK)
}

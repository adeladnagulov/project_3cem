package handlers

import (
	"encoding/json"
	"net/http"
	"project_3sem/internal/middleware"
)

func (h *UserHandle) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(middleware.IdKey).(string)
	email := r.Context().Value(middleware.EmailKey).(string)

	resp := map[string]interface{}{
		"id":    id,
		"email": email,
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

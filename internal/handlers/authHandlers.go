package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"project_3sem/internal/repositories"
)

type UserHandle struct {
	Repo repositories.RepoMemUs
}

func NewUserHandler(repo repositories.RepoMemUs) *UserHandle {
	return &UserHandle{Repo: repo}
}

func (h *UserHandle) Authorization(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Email \"%s\" request err: %s", req.Email, err)
		http.Error(w, "login failed", http.StatusInternalServerError)
		return
	}
	//отправляет код на почту и возращает верный ли ответ
	h.Repo.Authorization(req.Email) //если да то входит в аккаунт или создает нового и входит
	w.WriteHeader(http.StatusOK)
}

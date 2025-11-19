package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"project_3sem/internal/repositories"
	"project_3sem/internal/services"
	"strconv"
)

type UserHandle struct {
	RepoUsers repositories.RepoMemUs
	RepoCodes repositories.RepoMemCode
}

func NewUserHandler(repoUs repositories.RepoMemUs, repoCode repositories.RepoMemCode) *UserHandle {
	return &UserHandle{
		RepoUsers: repoUs,
		RepoCodes: repoCode,
	}
}

func (h *UserHandle) SendAuthCode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Email \"%s\" request err: %s", req.Email, err)
		http.Error(w, "Uncurrect email, err: "+err.Error(), http.StatusBadRequest)
		return
	}
	code := services.CreateCode()
	err := services.SendCodeToEmail(req.Email, strconv.Itoa(code))
	if err != nil {
		log.Printf("Send code error: %s", err)
		http.Error(w, "Send code error", http.StatusInternalServerError)
		return
	}

	h.RepoCodes.AddNewCode(code, req.Email)

	log.Printf("Send code Done")
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandle) Authorization(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		Code  int    `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Email: \"%s\", code: %d request err: %s", req.Email, req.Code, err)
		http.Error(w, "Uncurrect email or code err: "+err.Error(), http.StatusBadRequest)
		return
	}
	//реализовать хранение кода
	//проверка кода
	//авторизация
	//вернуть ответ
}

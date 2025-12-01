package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"project_3sem/internal/models"
	"project_3sem/internal/repositories"
	"project_3sem/internal/services"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserHandle struct {
	RepoUsers repositories.RepoUsers
	RepoCodes repositories.RepoMemCode
}

func NewUserHandler(repoUs repositories.RepoUsers, repoCode repositories.RepoMemCode) *UserHandle {
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

	ok, err := h.RepoCodes.ValidateCode(req.Code, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		log.Printf("Uncurrected code: %d", req.Code)
		http.Error(w, "Uncurrected code", http.StatusBadRequest)
		return
	}
	log.Printf("Code accepted")

	u := h.RepoUsers.Authorization(req.Email)
	token, err := GenerateJWTToken(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(token)
}

var jwtSecret = []byte("My_super_secret_key")

func GenerateJWTToken(u *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user": map[string]string{
			"id":    u.ID,
			"email": u.Email,
		},
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

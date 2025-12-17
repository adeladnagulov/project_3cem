package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"project_3sem/internal/repositories"
	"project_3sem/internal/responses"
	"project_3sem/internal/services"
)

type UserHandle struct {
	RepoUsers    repositories.RepoUsers
	RepoCodes    repositories.RepoMemCode
	EmailService services.EmailService
	TokenService services.TokenService
}

func NewUserHandler(repoUs repositories.RepoUsers, repoCode repositories.RepoMemCode, emailSer services.MyEmailService, tokenSer services.TokenService) *UserHandle {
	return &UserHandle{
		RepoUsers:    repoUs,
		RepoCodes:    repoCode,
		EmailService: &emailSer,
		TokenService: tokenSer,
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
	code := h.EmailService.CreateCode()
	if err := h.EmailService.ValidateEmail(req.Email); err != nil {
		log.Printf("Email \"%s\" request err: %s", req.Email, err)
		http.Error(w, "unacceptable email, err: "+err.Error(), http.StatusUnauthorized)
		return
	}
	err := h.EmailService.SendCodeToEmail(req.Email, code)
	if err != nil {
		log.Printf("Send code error: %s", err)
		http.Error(w, "Send code error", http.StatusTooManyRequests)
		return
	}

	h.RepoCodes.AddNewCode(code, req.Email)

	log.Printf("Send code Done")
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandle) Authorization(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Email: \"%s\", code: %s request err: %s", req.Email, req.Code, err)
		http.Error(w, "Uncurrect email or code err: "+err.Error(), http.StatusBadRequest)
		return
	}

	ok, err := h.RepoCodes.ValidateCode(req.Code, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		log.Printf("Uncurrected code: %s", req.Code)
		http.Error(w, "Uncurrected code", http.StatusBadRequest)
		return
	}
	log.Printf("Code accepted")

	u, err := h.RepoUsers.Authorization(req.Email)
	if err != nil {
		log.Printf("Authorization error: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken, err := h.TokenService.GenerateAccessToken(u)
	if err != nil {
		log.Printf("GenerateAccessToken errpr: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	refreshToken := h.TokenService.GenerateRefreshToken(u)

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
	})

	resp := map[string]interface{}{
		"accessToken": accessToken,
		"token_type":  "Bearer",
		"expires_in":  900,
	}
	responses.SendJSONResp(w, resp, http.StatusOK)
}

func (h *UserHandle) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		log.Printf("get token from cookie error: %s", err)
		http.Error(w, "No refresh token in cookie", http.StatusUnauthorized)
		return
	}

	id, ok := h.TokenService.ValidateRefreshToken(cookie.Value)
	if !ok {
		log.Printf("invalid refresh token: %s", id)
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	user := h.RepoUsers.GetUserByID(id)
	if user == nil {
		log.Printf("user == nil")
		http.Error(w, "invalid user", http.StatusUnauthorized)
		return
	}
	accessToken, err := h.TokenService.GenerateAccessToken(user)
	if err != nil {
		log.Printf("GenerateAccessToken error: %s", id)
		http.Error(w, "Access Token error", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"accessToken": accessToken,
		"token_type":  "Bearer",
	}
	responses.SendJSONResp(w, resp, http.StatusOK)
}

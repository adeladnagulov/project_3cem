package middleware

import (
	"context"
	"net/http"
	"project_3sem/internal/services"
	"strings"
)

type ctxKey string

const (
	IdKey    ctxKey = "id"
	EmailKey ctxKey = "email"
)

func AuthMiddlewera(ts *services.TokenServiceRepo, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "No token", http.StatusUnauthorized)
			return
		}
		partsToken := strings.Split(token, " ")
		if len(partsToken) != 2 || partsToken[0] != "Bearer" {
			http.Error(w, "Invalid format", http.StatusUnauthorized)
			return
		}
		u, err := ts.ValidateAccessToken(partsToken[1])
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), IdKey, u.ID)
		ctx = context.WithValue(ctx, EmailKey, u.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

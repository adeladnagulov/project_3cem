package services

import (
	"errors"
	"project_3sem/internal/models"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService interface {
	GenerateAccessToken(u *models.User) (string, error)
	ValidateAccessToken(tokenString string) (*models.User, error)
	GenerateRefreshToken(u *models.User) string
	ValidateRefreshToken(token string) (string, bool)
}

type TokenServiceSrtuct struct {
	jwtSecret     []byte
	mu            sync.RWMutex
	refreshTokens map[string]string
}

func NewTokenService(jwtSecret string) *TokenServiceSrtuct {
	return &TokenServiceSrtuct{
		jwtSecret:     []byte(jwtSecret),
		refreshTokens: make(map[string]string),
	}
}

func (ts *TokenServiceSrtuct) GenerateAccessToken(u *models.User) (string, error) {
	claims := jwt.MapClaims{
		"id":    u.ID,
		"email": u.Email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
		"type":  "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ts.jwtSecret)
}

func (ts *TokenServiceSrtuct) ValidateAccessToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return ts.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid claims format")
	}

	id, _ := claims["id"].(string)
	email, _ := claims["email"].(string)
	if id == "" || email == "" {
		return nil, errors.New("missing required claims")
	}

	u := models.User{
		ID:    id,
		Email: email,
	}
	return &u, nil
}

func (ts *TokenServiceSrtuct) GenerateRefreshToken(u *models.User) string {
	token := uuid.New().String()
	ts.mu.Lock()
	ts.refreshTokens[token] = u.ID
	ts.mu.Unlock()
	return token
}

func (ts *TokenServiceSrtuct) ValidateRefreshToken(token string) (string, bool) {
	ts.mu.RLock()
	id, ok := ts.refreshTokens[token]
	ts.mu.RUnlock()
	return id, ok
}

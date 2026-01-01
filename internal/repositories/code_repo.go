package repositories

import (
	"errors"
	"log"
	"project_3sem/internal/models"
	"sync"
	"time"
)

type RepoMemCode interface {
	AddNewCode(code string, creator string) error
	ValidateCode(code string, email string) (bool, error)
}

type MemoryRepoCodes struct {
	mu    sync.Mutex
	Codes map[string]*models.Code
}

func NewMemoryRepoCodes() *MemoryRepoCodes {
	return &MemoryRepoCodes{
		Codes: make(map[string]*models.Code),
	}
}

func (r *MemoryRepoCodes) AddNewCode(code string, creator string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.Codes[creator]; ok {
		delete(r.Codes, creator)
		log.Printf("Delete old code to email: %s", creator)
	}

	c := models.Code{
		Value:     code,
		ExpiresAt: time.Now().Add(time.Minute * 5),
	}
	r.Codes[creator] = &c
	log.Printf("AddCode: %s with creator: %s", code, creator)
	return nil
}

func (r *MemoryRepoCodes) ValidateCode(code string, email string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.Codes[email]; !ok {
		log.Printf("Code missing in repo")
		return false, errors.New("code no exist")
	}
	if time.Now().After(r.Codes[email].ExpiresAt) {
		return false, nil
	}
	return code == r.Codes[email].Value, nil
}

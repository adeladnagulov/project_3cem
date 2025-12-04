package repositories

import (
	"log"
	"project_3sem/internal/models"
	"sync"

	"github.com/google/uuid"
)

type RepoUsers interface {
	Authorization(email string) *models.User
	GetUserByID(id string) *models.User
}

type MemoryRepoUsers struct {
	mu    sync.Mutex
	Users map[string]*models.User
}

func NewMemoryRepoUsers() *MemoryRepoUsers {
	return &MemoryRepoUsers{
		Users: make(map[string]*models.User),
	}
}

func (r *MemoryRepoUsers) Authorization(email string) *models.User {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.Users[email]; !ok {
		u := &models.User{
			ID:    uuid.NewString(),
			Email: email,
		}
		r.Users[email] = u
		log.Printf("Add to memoryRepo new user to email: %s", email)
		return u
	} else {
		log.Printf("Put in memoryRepo user to email: %s", email)
		return r.Users[email]
	}
}

func (r *MemoryRepoUsers) GetUserByID(id string) *models.User {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range r.Users {
		if u.ID == id {
			return u
		}
	}
	return nil
}

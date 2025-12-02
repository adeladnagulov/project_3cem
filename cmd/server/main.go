package main

import (
	"fmt"
	"log"
	"net/http"
	"project_3sem/internal/handlers"
	"project_3sem/internal/repositories"
	"project_3sem/internal/services"

	"github.com/gorilla/mux"
)

func main() {
	handleUsers := handlers.NewUserHandler(repositories.NewMemoryRepoUsers(), repositories.NewMemoryRepoCodes(), *services.CreateEmailService())
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/auth/login", handleUsers.SendAuthCode).Methods("POST")
	r.HandleFunc("/api/v1/auth/confirm", handleUsers.Authorization).Methods("POST")

	fmt.Println("Start server to: 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

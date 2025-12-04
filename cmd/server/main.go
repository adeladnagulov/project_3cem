package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"project_3sem/internal/handlers"
	"project_3sem/internal/repositories"
	"project_3sem/internal/services"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		return
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	handleUsers := handlers.NewUserHandler(
		repositories.NewMemoryRepoUsers(),
		repositories.NewMemoryRepoCodes(),
		*services.CreateEmailService(),
		services.NewTokenService(jwtSecret),
	)
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/auth/login", handleUsers.SendAuthCode).Methods("POST")
	r.HandleFunc("/api/v1/auth/confirm", handleUsers.Authorization).Methods("POST")
	r.HandleFunc("/api/v1/auth/refresh", handleUsers.RefreshHandler).Methods("POST")

	fmt.Println("Start server to: 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

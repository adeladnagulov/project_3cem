package main

import (
	"fmt"
	"log"
	"net/http"
	"project_3sem/internal/handlers"
	"project_3sem/internal/repositories"

	"github.com/gorilla/mux"
)

func main() {
	handleUsers := handlers.NewUserHandler(repositories.NewMemoryRepoUsers())
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/auth/login/", handleUsers.Authorization).Methods("POST")

	fmt.Print("Start server to: 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

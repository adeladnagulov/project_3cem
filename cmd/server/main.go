package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"project_3sem/internal/handlers"
	"project_3sem/internal/middleware"
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
	tokenService := services.NewTokenService(jwtSecret)
	dbPostgres, err := repositories.NewPostgresDB()
	if err != nil {
		log.Fatalf("postgress error: %s", err.Error())
		return
	}

	redisClient, err := repositories.NewRedisDb(repositories.RedisCnf{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if err != nil {
		log.Fatalf("redis error: %s", err.Error())
		return
	}

	userRepo := repositories.NewPgRepoUsers(dbPostgres)
	siteRepo := repositories.NewPgRepoSites(dbPostgres)
	paymentRepo := repositories.NewPgRepoPayments(dbPostgres)
	handleUsers := handlers.NewUserHandler(
		userRepo,
		repositories.NewRedesRepoCodes(redisClient),
		*services.CreateEmailService(),
		tokenService,
	)
	handleSite := handlers.NewSiteHandler(
		siteRepo,
		paymentRepo,
	)
	handlePayment := handlers.NewPaymentHandler(
		paymentRepo,
		*userRepo,
	)
	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return middleware.SubdomainMiddlewera(next)
	})
	r.Use(func(next http.Handler) http.Handler {
		return middleware.CORSmiddlewera(next)
	})

	r.HandleFunc("/api/v1/auth/login", handleUsers.SendAuthCode).Methods("POST")
	r.HandleFunc("/api/v1/auth/confirm", handleUsers.Authorization).Methods("POST")
	r.HandleFunc("/api/v1/auth/refresh", handleUsers.RefreshHandler).Methods("POST")

	r.HandleFunc("/api/v1/payment/webhook", handlePayment.PaymentWebhook).Methods("POST")

	protected := r.PathPrefix("/api/v1/me").Subrouter()
	protected.Use(func(next http.Handler) http.Handler {
		return middleware.AuthMiddlewera(tokenService, next)
	})
	protected.HandleFunc("/dashboard", handleUsers.DashboardHandler).Methods("POST")
	protected.HandleFunc("/sites/save", handleSite.SaveDraft).Methods("POST")
	protected.HandleFunc("/sites/{id}/publish", handleSite.Publish).Methods("POST")

	protected.HandleFunc("/payment/create", handlePayment.CreatePayments).Methods("POST")

	r.HandleFunc("/site-config", handleSite.RenderSite).Methods("GET")

	fmt.Println("Start server to: 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

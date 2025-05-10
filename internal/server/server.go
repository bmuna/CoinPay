package server

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	internal "github.com/bmuna/CoinPay/backend/internal/auth"
	"github.com/bmuna/CoinPay/backend/internal/config"
	"github.com/bmuna/CoinPay/backend/internal/controller"
	"github.com/bmuna/CoinPay/backend/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	_ "github.com/lib/pq"
)

func NewServer() http.Handler {
	internal.NewAuth()

	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		log.Fatal("DB is not found in the enviroment")
	}

	conn, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("Failed to connect with the DB", err)
	}
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	apiCfg := config.ApiConfig{
		DB: database.New(conn),
	}

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Get("/healthz", controller.HandlerReadiness)
	router.Get("/auth/{provider}", controller.BiginAuthProviderCallback)
	router.Get("/auth/{provider}/callback", controller.GetAuthCallBackFuction(&apiCfg))
	router.Post("/api/signin", controller.Signin(&apiCfg))
	router.Post("/api/signup", controller.Signup(&apiCfg))

	return router

}

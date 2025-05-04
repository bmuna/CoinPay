package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func NewServer() http.Handler {

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Get("/healthz", handlerReadiness)

	return router

}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	responseWithJSON(w, 200, struct{}{})
}

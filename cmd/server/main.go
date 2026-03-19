package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/fyreex/physicalc/api"
)

// @title           Physicalc API
// @version         1.0
// @description     Calculadora de Física e Cálculo em Go
// @host            localhost:8080
// @BasePath        /api
func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}))

	api.RegisterRoutes(r)

	fmt.Println("🚀 Physicalc rodando em http://localhost:8080")
	fmt.Println("📄 Swagger UI em  http://localhost:8080/swagger/")

	log.Fatal(http.ListenAndServe(":8080", r))
}

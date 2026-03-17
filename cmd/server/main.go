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

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}))

	// Registra todas as rotas
	api.RegisterRoutes(r)

	fmt.Println("🚀 Physicalc rodando em http://localhost:8080")
	fmt.Println("📖 Swagger UI em  http://localhost:8080/swagger/index.html")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}

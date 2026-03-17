package api

import (
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registra todas as rotas da aplicação no router.
// Cada grupo de rotas corresponde a um módulo matemático/físico.
func RegisterRoutes(r *chi.Mux) {
	r.Route("/api", func(r chi.Router) {

		// --- Cálculo ---
		r.Route("/calculus", func(r chi.Router) {
			r.Post("/integrate", IntegrateHandler)   // integração numérica
			r.Post("/derivative", DerivativeHandler) // derivada em um ponto
			r.Post("/limit", LimitHandler)           // limite numérico
		})

		// --- Álgebra Linear ---
		r.Route("/algebra", func(r chi.Router) {
			r.Post("/vector", VectorHandler) // operações com vetores
			r.Post("/matrix", MatrixHandler) // operações com matrizes
		})

		// --- Física ---
		r.Route("/physics", func(r chi.Router) {
			r.Post("/kinematics", KinematicsHandler) // cinemática
			r.Post("/dynamics", DynamicsHandler)     // dinâmica
		})

		// --- Equações Diferenciais ---
		r.Route("/ode", func(r chi.Router) {
			r.Post("/solve", ODESolveHandler) // resolver EDO numericamente
		})

		// Health check
		r.Get("/health", HealthHandler)
	})
}

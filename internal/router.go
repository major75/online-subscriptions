package router

import (
	"github.com/go-chi/chi/v5"

	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/major75/online-subscriptions/database"
	"github.com/major75/online-subscriptions/internal/handlers"
	"github.com/major75/online-subscriptions/internal/repository/subscriptions"
	"github.com/major75/online-subscriptions/pkg/logger"
)

func NewRouter(log logger.Logger, db *database.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Mount("/docs", httpSwagger.WrapHandler)

	userRepo := subscriptions.NewUserSubscriptionRepository(db, log)

	healthHandler := handlers.NewHealthHandler(log)
	userHandler := handlers.NewUserHandler(log, userRepo)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", healthHandler.HealthCheck)
		r.Route("/subscriptions", func(r chi.Router) {
			// Create user subscription
			r.Post("/", userHandler.CreateUserSubscription)

			// Get subscription details
			r.Get("/{id}", userHandler.GetSubscription)

			// Delete subscription
			r.Delete("/{id}", userHandler.DeleteSubscription)

			// Update subscription
			r.Put("/{id}", userHandler.UpdateSubscription)

			// Patch subscription
			r.Patch("/{id}", userHandler.PatchSubscription)

			// Get subscriptions report total
			r.Get("/total", userHandler.GetSubscriptionsTotal)
		})

		r.Route("/user", func(r chi.Router) {
			r.Get("/{userID}/subscriptions", userHandler.GetUserSubscriptions)
		})
	})

	return r
}

package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Post("/api/payment-intent", app.GetPaymentIntent)

	mux.Get("/api/widget/{id}", app.GetWidgetByID)

	mux.Post("/api/create-customer-and-subscribe-to-plan", app.CreateCustomerAndSubscribeToPlan)

	mux.Post("/api/authenticate", app.CreateAuthToken)

	mux.Post("/api/is-authenticated", app.CheckAuthentication)

	mux.Post("/api/forgot-password", app.SendPasswordResetEmail)

	mux.Route("/api/admin", func(mux chi.Router) {
		mux.Use(app.Auth)

		mux.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Got in"))
		})

		mux.Post("/virtual-terminal-succeeded", app.VirtualTerminalPaymentSucceeded)
	})

	return mux
}

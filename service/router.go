package service

import (
	"net/http"

	"github.com/go-chi/chi"
)

// RouterBasePath Variable
var RouterBasePath string

// Router Variable
var Router *chi.Mux

// routerInit Function
func routerInit() {
	// Initialize Router
	Router = chi.NewRouter()

	// Set Router Middleware
	Router.Use(routerCORS)
	Router.Use(routerRealIP)
	Router.Use(routerLogs)
	Router.Use(routerEntitySize)

	// Set Router Handler
	Router.NotFound(handlerNotFound)
	Router.MethodNotAllowed(handlerMethodNotAllowed)
	Router.Get("/favicon.ico", handlerFavIcon)
}

// HealthCheck Function
func HealthCheck(w http.ResponseWriter) {
	// Return Success
	ResponseSuccess(w, "")
}

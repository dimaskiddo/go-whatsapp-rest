package router

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/dimaskiddo/go-whatsapp-rest/pkg/server"
)

// RouterBasePath Variable
var RouterBasePath string

// Router Variable
var Router *chi.Mux

// Initialize Function in Router
func init() {
	// Initialize Router
	Router = chi.NewRouter()
	RouterBasePath = server.Config.GetString("ROUTER_BASE_PATH")

	// Set Router CORS Configuration
	routerCORSCfg.Origins = server.Config.GetString("CORS_ALLOWED_ORIGIN")
	routerCORSCfg.Methods = server.Config.GetString("CORS_ALLOWED_METHOD")
	routerCORSCfg.Headers = server.Config.GetString("CORS_ALLOWED_HEADER")

	// Set Router Middleware
	Router.Use(routerCORS)
	Router.Use(routerRealIP)
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

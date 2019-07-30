package router

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/dimaskiddo/go-whatsapp-rest/hlp"
)

// RouterBasePath Variable
var RouterBasePath string

// Router Variable
var Router *chi.Mux

// Initialize Function in Router
func init() {
	// Initialize Router
	Router = chi.NewRouter()
	RouterBasePath = hlp.Config.GetString("ROUTER_BASE_PATH")

	// Set Router CORS Configuration
	routerCORSCfg.Origins = hlp.Config.GetString("CORS_ALLOWED_ORIGIN")
	routerCORSCfg.Methods = hlp.Config.GetString("CORS_ALLOWED_METHOD")
	routerCORSCfg.Headers = hlp.Config.GetString("CORS_ALLOWED_HEADER")

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

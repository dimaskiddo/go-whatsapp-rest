package service

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

// ResSuccess Struct
type ResSuccess struct {
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ResError Struct
type ResError struct {
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

// Router CORS Configuration Struct
type routerCORSConfig struct {
	Origins string
	Methods string
	Headers string
}

// Router CORS Configuration Variable
var routerCORSCfg routerCORSConfig

// RouterBasePath Variable
var RouterBasePath string

// Router Variable
var Router *chi.Mux

// routerInit Function
func routerInit() {
	// Initialize Router
	Router = chi.NewRouter()

	// Set Router Entity Size
	Router.Use(routerEntitySize)

	// Set Router CORS
	Router.Use(routerCORS)

	// Set Router Logging
	Router.Use(routerLogs)

	// Set Handler for /favicon.ico
	Router.Get("/favicon.ico", handlerFavIcon)

	// Set Handler for Not Found
	Router.NotFound(handlerNotFound)

	// Set Handler for Method Not Allowed
	Router.MethodNotAllowed(handlerMethodNotAllowed)
}

// RouterEntitySize Function
func routerEntitySize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate Entity Size
		r.Body = http.MaxBytesReader(w, r.Body, Config.GetInt64("SERVER_UPLOAD_LIMIT"))
		next.ServeHTTP(w, r)
	})
}

// RouterCORS Function
func routerCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add Header for CORS
		w.Header().Set("Access-Control-Allow-Origin", routerCORSCfg.Origins)
		w.Header().Set("Access-Control-Allow-Methods", routerCORSCfg.Methods)
		w.Header().Set("Access-Control-Allow-Headers", routerCORSCfg.Headers)
		next.ServeHTTP(w, r)
	})
}

// routerLogs Function
func routerLogs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log HTTP Access if Not Acessing /favicon.ico
		if r.RequestURI != "/favicon.ico" {
			Log("info", "http-access", "access method "+r.Method+" at URI "+r.RequestURI)
		}
		next.ServeHTTP(w, r)
	})
}

// HandlerNotFound Function
func handlerNotFound(w http.ResponseWriter, r *http.Request) {
	Log("warn", "http-access", "not found method "+r.Method+" at URI "+r.RequestURI)
	ResponseNotFound(w, "not found method "+r.Method+" at URI "+r.RequestURI)
}

// HandlerMethodNotAllowed Function
func handlerMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	Log("warn", "http-access", "not allowed method "+r.Method+" at URI "+r.RequestURI)
	ResponseMethodNotAllowed(w, "not allowed method "+r.Method+" at URI "+r.RequestURI)
}

// HandlerFavIcon Function
func handlerFavIcon(w http.ResponseWriter, r *http.Request) {
	ResponseNoContent(w)
}

// HealthCheck Function
func HealthCheck(w http.ResponseWriter) {
	// Return Success
	ResponseSuccess(w, "")
}

// ResponseWrite Function
func ResponseWrite(w http.ResponseWriter, responseCode int, responseData interface{}) {
	// Write Response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode)

	// Write JSON to Response
	json.NewEncoder(w).Encode(responseData)
}

// ResponseSuccess Function
func ResponseSuccess(w http.ResponseWriter, message string) {
	var response ResSuccess

	// Set Default Message
	if len(message) == 0 {
		message = "Success"
	}

	// Set Response Data
	response.Status = true
	response.Code = http.StatusOK
	response.Message = message

	// Set Response Data to HTTP
	ResponseWrite(w, response.Code, response)
}

// ResponseCreated Function
func ResponseCreated(w http.ResponseWriter) {
	var response ResSuccess

	// Set Response Data
	response.Status = true
	response.Code = http.StatusCreated
	response.Message = "Created"

	// Set Response Data to HTTP
	ResponseWrite(w, response.Code, response)
}

// ResponseUpdated Function
func ResponseUpdated(w http.ResponseWriter) {
	var response ResSuccess

	// Set Response Data
	response.Status = true
	response.Code = http.StatusOK
	response.Message = "Updated"

	// Set Response Data to HTTP
	ResponseWrite(w, response.Code, response)
}

// ResponseNoContent Function
func ResponseNoContent(w http.ResponseWriter) {
	w.WriteHeader(204)
}

// ResponseNotFound Function
func ResponseNotFound(w http.ResponseWriter, message string) {
	var response ResError

	// Set Default Message
	if len(message) == 0 {
		message = "Not Found"
	}

	// Set Response Data
	response.Status = false
	response.Code = http.StatusNotFound
	response.Message = "Not Found"
	response.Error = message

	// Set Response Data to HTTP
	ResponseWrite(w, response.Code, response)
}

// ResponseMethodNotAllowed Function
func ResponseMethodNotAllowed(w http.ResponseWriter, message string) {
	var response ResError

	// Set Default Message
	if len(message) == 0 {
		message = "Method Not Allowed"
	}

	// Set Response Data
	response.Status = false
	response.Code = http.StatusMethodNotAllowed
	response.Message = "Method Not Allowed"
	response.Error = message

	// Set Response Data to HTTP
	ResponseWrite(w, response.Code, response)
}

// ResponseBadRequest Function
func ResponseBadRequest(w http.ResponseWriter, message string) {
	var response ResError

	// Set Default Message
	if len(message) == 0 {
		message = "Bad Request"
	}

	// Set Response Data
	response.Status = false
	response.Code = http.StatusBadRequest
	response.Message = "Bad Request"
	response.Error = message

	// Logging Error
	Log("error", "http-access", strings.ToLower(message))

	// Set Response Data to HTTP
	ResponseWrite(w, response.Code, response)
}

// ResponseInternalError Function
func ResponseInternalError(w http.ResponseWriter, message string) {
	var response ResError

	// Set Default Message
	if len(message) == 0 {
		message = "Internal Server Error"
	}

	// Set Response Data
	response.Status = false
	response.Code = http.StatusInternalServerError
	response.Message = "Internal Server Error"
	response.Error = message

	// Logging Error
	Log("error", "http-access", strings.ToLower(message))

	// Set Response Data to HTTP
	ResponseWrite(w, response.Code, response)
}

// ResponseUnauthorized Function
func ResponseUnauthorized(w http.ResponseWriter) {
	var response ResError

	// Set Response Data
	response.Status = false
	response.Code = http.StatusUnauthorized
	response.Message = "Unauthorized"
	response.Error = "Unaothorized"

	// Set Response Data to HTTP
	ResponseWrite(w, response.Code, response)
}

// ResponseAuthenticate Function
func ResponseAuthenticate(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Authorization Required"`)
	ResponseUnauthorized(w)
}

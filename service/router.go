package service

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// FormatSuccess Struct
type FormatSuccess struct {
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// FormatError Struct
type FormatError struct {
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

// Router CORS Configuration Struct
type routerCORSConfig struct {
	Headers []string
	Origins []string
	Methods []string
}

// Router CORS Configuration Variable
var routerCORSCfg routerCORSConfig

// RouterBasePath Variable
var RouterBasePath string

// RouterHandler Variable
var RouterHandler http.Handler

// Router Variable
var Router *mux.Router

// InitRouter Function
func initRouter() {
	// Initialize Router
	Router = mux.NewRouter().StrictSlash(true)

	// Set Router Handler with Logging & CORS Support
	RouterHandler = handlers.LoggingHandler(os.Stdout, handlers.CORS(
		handlers.AllowedHeaders(routerCORSCfg.Headers),
		handlers.AllowedOrigins(routerCORSCfg.Origins),
		handlers.AllowedMethods(routerCORSCfg.Methods))(Router))

	// Set Router Default Not Found Handler
	Router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ResponseNotFound(w, "Not Found Method "+r.Method+" at URI "+r.RequestURI)
	})
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
	var response FormatSuccess

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
	var response FormatSuccess

	// Set Response Data
	response.Status = true
	response.Code = http.StatusCreated
	response.Message = "Created"

	// Set Response Data to HTTP
	ResponseWrite(w, response.Code, response)
}

// ResponseUpdated Function
func ResponseUpdated(w http.ResponseWriter) {
	var response FormatSuccess

	// Set Response Data
	response.Status = true
	response.Code = http.StatusOK
	response.Message = "Updated"

	// Set Response Data to HTTP
	ResponseWrite(w, response.Code, response)
}

// ResponseNotFound Function
func ResponseNotFound(w http.ResponseWriter, message string) {
	var response FormatError

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

// ResponseBadRequest Function
func ResponseBadRequest(w http.ResponseWriter, message string) {
	var response FormatError

	// Set Default Message
	if len(message) == 0 {
		message = "Bad Request"
	}

	// Set Response Data
	response.Status = false
	response.Code = http.StatusBadRequest
	response.Message = "Bad Request"
	response.Error = message

	// Set Response Data to HTTP
	ResponseWrite(w, response.Code, response)
}

// ResponseInternalError Function
func ResponseInternalError(w http.ResponseWriter, message string) {
	var response FormatError

	// Set Default Message
	if len(message) == 0 {
		message = "Internal Server Error"
	}

	// Set Response Data
	response.Status = false
	response.Code = http.StatusInternalServerError
	response.Message = "Internal Server Error"
	response.Error = message

	// Set Response Data to HTTP
	ResponseWrite(w, response.Code, response)
}

// ResponseUnauthorized Function
func ResponseUnauthorized(w http.ResponseWriter) {
	var response FormatError

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

package router

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/dimaskiddo/go-whatsapp-rest/hlp"
)

// ResSuccess Struct
type ResSuccess struct {
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ResSuccessWithData Struct
type ResSuccessWithData struct {
	Status  bool        `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ResError Struct
type ResError struct {
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error"`
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

// ResponseSuccessWithData Function
func ResponseSuccessWithData(w http.ResponseWriter, message string, data interface{}) {
	var response ResSuccessWithData

	// Set Default Message
	if len(message) == 0 {
		message = "Success"
	}

	// Set Response Data
	response.Status = true
	response.Code = http.StatusOK
	response.Message = message
	response.Data = data

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
	hlp.LogPrintln(hlp.LogLevelError, "http-access", strings.ToLower(message))

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
	hlp.LogPrintln(hlp.LogLevelError, "http-access", strings.ToLower(message))

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
	response.Error = "Unauthorized"

	// Set Response Data to HTTP
	ResponseWrite(w, response.Code, response)
}

// ResponseAuthenticate Function
func ResponseAuthenticate(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="Authorization Required"`)
	ResponseUnauthorized(w)
}

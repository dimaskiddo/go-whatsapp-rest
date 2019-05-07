package service

import (
	"net/http"
)

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

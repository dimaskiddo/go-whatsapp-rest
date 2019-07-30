package router

import (
	"net/http"
	"strings"

	"github.com/dimaskiddo/go-whatsapp-rest/hlp"
)

// Router CORS Configuration Struct
type routerCORSConfig struct {
	Origins string
	Methods string
	Headers string
}

// Router CORS Configuration Variable
var routerCORSCfg routerCORSConfig

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

// RouterRealIP Function
func routerRealIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Real IP from Cannoical Header
		if XForwardedFor := r.Header.Get(http.CanonicalHeaderKey("X-Forwarded-For")); XForwardedFor != "" {
			dataIndex := strings.Index(XForwardedFor, ", ")
			if dataIndex == -1 {
				dataIndex = len(XForwardedFor)
			}
			r.RemoteAddr = XForwardedFor[:dataIndex]
		} else if XRealIP := r.Header.Get(http.CanonicalHeaderKey("X-Real-IP")); XRealIP != "" {
			r.RemoteAddr = XRealIP
		}
		next.ServeHTTP(w, r)
	})
}

// RouterLogs Function
func routerLogs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log HTTP Access if Not Acessing /favicon.ico
		if r.RequestURI != "/favicon.ico" {
			hlp.LogPrintln(hlp.LogLevelInfo, "http-access", "access method "+r.Method+" at URI "+r.RequestURI)
		}
		next.ServeHTTP(w, r)
	})
}

// RouterEntitySize Function
func routerEntitySize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate Entity Size
		r.Body = http.MaxBytesReader(w, r.Body, hlp.Config.GetInt64("SERVER_UPLOAD_LIMIT"))
		next.ServeHTTP(w, r)
	})
}

package index

import (
	"net/http"

	"github.com/dimaskiddo/go-whatsapp-rest/pkg/router"
)

// GetIndex Function to Show API Information
func GetIndex(w http.ResponseWriter, r *http.Request) {
	router.ResponseSuccess(w, "Go WhatsApp REST is running")
}

// GetHealth Function to Show Health Check Status
func GetHealth(w http.ResponseWriter, r *http.Request) {
	router.HealthCheck(w)
}

package controller

import (
	"net/http"

	svc "github.com/dimaskiddo/go-whatsapp-rest/service"
)

// GetIndex Function to Show API Information
func GetIndex(w http.ResponseWriter, r *http.Request) {
	svc.ResponseSuccess(w, "WhatsApp Go Service is running")
}

// GetHealth Function to Show Health Check Status
func GetHealth(w http.ResponseWriter, r *http.Request) {
	svc.HealthCheck(w)
}

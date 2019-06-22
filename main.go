package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	svc "github.com/dimaskiddo/go-whatsapp-rest/service"
)

// Server Variable
var svr *svc.Server

// Init Function
func init() {
	// Initialize Routes
	routesInit()

	// Initialize Server
	svr = svc.NewServer(svc.Router)
}

// Main Function
func main() {
	// Starting Server
	svr.Start()

	// Make Channel for OS Signal
	sig := make(chan os.Signal, 1)

	// Notify Any Signal to OS Signal Channel
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	// Return OS Signal Channel
	// As Exit Sign
	<-sig

	// Log Break Line
	fmt.Println("")

	// Stopping Server
	defer svr.Stop()
}

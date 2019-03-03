package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	svc "github.com/dimaskiddo/go-whatsapp-rest/service"
)

// Main Server Variable
var mainServer *svc.Server

// Init Function
func init() {
	// Initialize service
	svc.Initialize()

	// Initialize Routes
	log.Println("Initialize - Routes")
	initRoutes()

	// Initialize Server
	log.Println("Initialize - Server")
	mainServer = svc.NewServer(svc.RouterHandler)
}

// Main Function
func main() {
	// Starting Server
	mainServer.Start()

	// Make Channel to Catch OS Signal
	osSignal := make(chan os.Signal, 1)

	// Catch OS Signal from Channel
	signal.Notify(osSignal, os.Interrupt)
	signal.Notify(osSignal, syscall.SIGTERM)

	// Return OS Signal as Exit Code
	<-osSignal

	// Termination Symbol Log Line
	fmt.Println("")

	// Stopping Server
	defer mainServer.Stop()
}

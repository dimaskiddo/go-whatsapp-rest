package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dimaskiddo/go-whatsapp-rest/pkg/router"
	"github.com/dimaskiddo/go-whatsapp-rest/pkg/server"

	"github.com/dimaskiddo/go-whatsapp-rest/internal"
)

// Server Variable
var svr *server.Server

// Init Function
func init() {
	// Set Go Log Flags
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	// Load Routes
	internal.LoadRoutes()

	// Initialize Server
	svr = server.NewServer(router.Router)
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

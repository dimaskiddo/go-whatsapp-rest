package service

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

// Server Struct
type Server struct {
	srv *http.Server
	wg  sync.WaitGroup
}

// Server Configuration Struct
type serverConfig struct {
	IP   string
	Port string
}

// Server Configuration Variable
var serverCfg serverConfig

// NewServer Function to Create a New Server Handler
func NewServer(handler http.Handler) *Server {
	// Initialize New Server
	return &Server{
		srv: &http.Server{
			Addr:    net.JoinHostPort(serverCfg.IP, serverCfg.Port),
			Handler: handler,
		},
	}
}

// Start Method for Server
func (s *Server) Start() {
	// Initialize Context Handler Without Timeout
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Add to The WaitGroup for The Listener GoRoutine
	// And Wait for 1 Routine to be Done
	s.wg.Add(1)

	// Start The Server
	go func() {
		log.Println("Server - Starting")
		log.Println("Server - Started at " + net.JoinHostPort(serverCfg.IP, serverCfg.Port))
		s.srv.ListenAndServe()

		s.wg.Done()
	}()
}

// Stop Method for Server
func (s *Server) Stop() {
	// Initialize Timeout
	timeout := 5 * time.Second

	// Initialize Context Handler With Timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Hanlde Any Error While Stopping Server
	if err := s.srv.Shutdown(ctx); err != nil {
		if err = s.srv.Close(); err != nil {
			log.Println(err.Error())
			return
		}
	}
	s.wg.Wait()
	log.Println("Server - Stopping")
	log.Println("Server - Stopped from " + net.JoinHostPort(serverCfg.IP, serverCfg.Port))
}

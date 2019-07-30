package hlp

import (
	"context"
	"net"
	"net/http"
	"os"
	"strconv"
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
	LogPrintln(LogLevelInfo, "http-server", "server master started at PID "+strconv.Itoa(os.Getpid()))
	go func() {
		LogPrintln(LogLevelInfo, "http-server", "server worker started at PID "+strconv.Itoa(os.Getpid())+" listening on "+net.JoinHostPort(serverCfg.IP, serverCfg.Port))
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
			LogPrintln("error", "http-server", err.Error())
			return
		}
	}
	s.wg.Wait()
}

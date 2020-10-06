package server

import (
	"context"
	"fmt"
	"log"
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
var ServerCfg serverConfig

// NewServer Function to Create a New Server Handler
func NewServer(handler http.Handler) *Server {
	// Initialize New Server
	return &Server{
		srv: &http.Server{
			Addr:    net.JoinHostPort(ServerCfg.IP, ServerCfg.Port),
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
	log.Println("{\"label\":\"server-http\",\"level\":\"info\",\"msg\":\"server master started at pid " + strconv.Itoa(os.Getpid()) + "\",\"service\":\"" + Config.GetString("SERVER_NAME") + "\",\"time\":" + fmt.Sprint(time.Now().Format(time.RFC3339Nano)) + "\"}")
	go func() {
		log.Println("{\"label\":\"server-http\",\"level\":\"info\",\"msg\":\"server worker started at pid " + strconv.Itoa(os.Getpid()) + " listening on " + net.JoinHostPort(ServerCfg.IP, ServerCfg.Port) + "\",\"service\":\"" + Config.GetString("SERVER_NAME") + "\",\"time\":" + fmt.Sprint(time.Now().Format(time.RFC3339Nano)) + "\"}")
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
			log.Fatalln("{\"label\":\"server-http\",\"level\":\"error\",\"msg\":\"" + err.Error() + "\",\"service\":\"" + Config.GetString("SERVER_NAME") + "\",\"time\":" + fmt.Sprint(time.Now().Format(time.RFC3339Nano)) + "\"}")
			return
		}
	}
	s.wg.Wait()
}

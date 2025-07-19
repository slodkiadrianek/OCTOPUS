package api

import (
	"net/http"
	"time"
)

type Config struct {
	Port   string
	Logger *logger
}

type Server struct {
	config *Config
	server *http.Server
}

func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:         ":" + s.config.Port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return s.server.ListenAndServe()
}

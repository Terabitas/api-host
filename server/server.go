package server

import (
	"net/http"

	log "github.com/nildev/api-host/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/nildev/api-host/config"
	"github.com/nildev/api-host/endpoints"
	"github.com/nildev/api-host/version"
)

// Server type
type Server struct {
	stop    chan bool
	cfg     config.Config
	handler http.Handler
}

// New type
func New(cfg config.Config) (*Server, error) {
	srv := Server{
		cfg:     cfg,
		stop:    nil,
		handler: endpoints.Router(cfg),
	}
	return &srv, nil
}

// Run server
func (s *Server) Run() {
	ctxLog := log.WithField("version", version.Version).WithField("git-hash", version.GitHash).WithField("build-time", version.BuiltTimestamp)

	ctxLog.Infof("Starting api-host service [%s:%s]", s.cfg.IP, s.cfg.Port)
	s.stop = make(chan bool)

	go func() {
		ctxLog.Infof("Starting HTTP server ...")
		if err := http.ListenAndServe(s.cfg.IP+":"+s.cfg.Port, s.handler); err != nil {
			ctxLog.Fatalf("Unable to create listener, %s", err)
		}
	}()
}

// Stop server
func (s *Server) Stop() {
	close(s.stop)
}

// Purge server
func (s *Server) Purge() {
}

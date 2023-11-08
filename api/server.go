package api

import (
	"log/slog"
	"net/http"

	"github.com/R894/go-blog/models"
)

type Server struct {
	listenAddr string
	logger     *slog.Logger
	db         models.Database
}

func NewServer(listenAddr string, logger *slog.Logger, db models.Database) *Server {
	return &Server{
		listenAddr: listenAddr,
		logger:     logger,
		db:         db,
	}
}

func (s *Server) Start() error {
	s.logger.Info("Server started", "Port", s.listenAddr)
	return http.ListenAndServe("localhost:"+s.listenAddr, s.routes())
}

package api

import (
	"log/slog"
	"net/http"

	"github.com/R894/go-blog/api/interfaces"
	"github.com/R894/go-blog/models"
)

type Server struct {
	listenAddr string
	logger     *slog.Logger
	db         models.Database
}

func NewServer(listenAddr string, logger *slog.Logger, db models.Database) interfaces.ServerInterface {
	return &Server{
		listenAddr: listenAddr,
		logger:     logger,
		db:         db,
	}
}

func (s *Server) Start() error {
	s.logger.Info("Server started", "Port", s.listenAddr)
	return http.ListenAndServe("0.0.0.0:"+s.listenAddr, s.routes())
}

func (s *Server) GetDB() models.Database {
	return s.db
}

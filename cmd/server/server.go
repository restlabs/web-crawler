package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/we-are-discussing-rest/web-crawler/internal/logger"
	"github.com/we-are-discussing-rest/web-crawler/internal/repository"
	"net/http"
)

type Server struct {
	store repository.Repository
	*chi.Mux
	*logger.Logger
}

func NewServer(store repository.Repository, logger *logger.Logger) *Server {
	s := new(Server)
	s.store = store

	router := chi.NewRouter()

	router.Use(logger.LoggerMiddleware)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)

	router.Get("/healthcheck", s.handleHealthcheck)

	for _, v := range router.Routes() {
		if v.SubRoutes != nil {
			for _, v := range v.SubRoutes.Routes() {
				logger.Info("registering", "route", v.Pattern)
			}
		}
		logger.Info("registering", "route", v.Pattern)
	}

	s.Logger = logger
	s.Mux = router

	return s
}

func (s *Server) handleHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

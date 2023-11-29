package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/we-are-discussing-rest/web-crawler/logger"
	"github.com/we-are-discussing-rest/web-crawler/repository"
	"net/http"
)

type SeedUrlsDto struct {
	SeedUrls []string `json:"seedUrls"`
}

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
	router.Route("/api/v1", func(r chi.Router) {
		r.Post("/seed", s.seedHandler)
	})

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

func (s *Server) seedHandler(w http.ResponseWriter, r *http.Request) {
	var seedUrl SeedUrlsDto
	if err := json.NewDecoder(r.Body).Decode(&seedUrl); err != nil {
		s.Logger.Error("error parsing body", "error", err)
		fmt.Fprintf(w, "%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, v := range seedUrl.SeedUrls {
		if insertErr := s.store.Insert(v); insertErr != nil {
			s.Logger.Error("error seeding urls", "error", insertErr, "url", v)
			fmt.Fprintf(w, "%v", insertErr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

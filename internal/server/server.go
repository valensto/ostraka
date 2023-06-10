package server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/valensto/ostraka/internal/config/env"
	"github.com/valensto/ostraka/internal/logger"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
)

type Server struct {
	Router *chi.Mux
	Port   string
	Host   string
}

func New(config *env.Config) *Server {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)

	return &Server{
		Router: mux,
		Port:   config.Port,
		Host:   config.Host,
	}
}

func (s *Server) Run() error {
	h2s := &http2.Server{}
	server := &http.Server{
		Addr:    ":" + s.Port,
		Handler: h2c.NewHandler(s.Router, h2s),
	}
	return server.ListenAndServe()
}

func (s *Server) Respond(w http.ResponseWriter, _ *http.Request, status int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if data == nil {
		return
	}

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		logger.Get().Warn().Msgf("cannot format response json. err=%v\n", err)
	}
}

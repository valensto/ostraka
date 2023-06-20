package http

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	mid "github.com/go-chi/chi/v5/middleware"
	"github.com/valensto/ostraka/internal/env"
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/middleware"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
)

type Server struct {
	Router *chi.Mux
	Port   string
	Host   string

	Middlewares *middleware.Middlewares
}

func New(config *env.Config) *Server {
	mux := chi.NewRouter()
	mux.Use(mid.Recoverer)
	mux.Use(mid.Logger)

	return &Server{
		Router: mux,
		Port:   config.Port,
		Host:   config.Host,
	}
}

func (s *Server) Serve() error {
	h2s := &http2.Server{}
	server := &http.Server{
		Addr:    ":" + s.Port,
		Handler: h2c.NewHandler(s.Router, h2s),
	}
	return server.ListenAndServe()
}

func (s *Server) Respond(w http.ResponseWriter, _ *http.Request, status int, data any) {
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

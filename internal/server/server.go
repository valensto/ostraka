package server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/valensto/ostraka/internal/logger"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
)

type Server struct {
	Router *chi.Mux
	port   string
}

func New(port string) *Server {
	return &Server{
		Router: chi.NewRouter(),
		port:   port,
	}
}

func (s *Server) Serve() error {
	s.Router.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("webui/dist/assets"))))
	s.Router.Get("/dashboard", s.webui())

	h2s := &http2.Server{}
	server := &http.Server{
		Addr:    ":" + s.port,
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

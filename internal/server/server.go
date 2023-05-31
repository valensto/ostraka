package server

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/valensto/ostraka/internal/logger"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"html/template"
	"net/http"
)

type Server struct {
	Router *chi.Mux
	port   string
}

func New(port string) *Server {
	s := &Server{
		Router: chi.NewRouter(),
		port:   port,
	}

	s.initializeRouter()
	return s
}

func (s *Server) initializeRouter() {
	s.Router.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("webui/dist/assets"))))
	s.Router.Get("/dashboard", s.webui())
}

func (s *Server) Serve() error {
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

func (s *Server) webui() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("webui/dist/index.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			s.Respond(w, r, http.StatusInternalServerError, nil)
			return
		}

		s.Respond(w, r, http.StatusOK, nil)
	}
}

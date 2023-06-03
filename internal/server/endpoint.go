package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"html/template"
	"net/http"
)

type Method string

const (
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	PATCH  Method = "PATCH"
	DELETE Method = "DELETE"
	OPTION Method = "OPTION"
)

func (m Method) String() string {
	return string(m)
}

func (m Method) validate() error {
	switch m {
	case GET, POST, PUT, PATCH, DELETE, OPTION:
		return nil
	default:
		return fmt.Errorf("invalid method %s", m)
	}
}

type Endpoint struct {
	Method      Method
	Path        string
	Cors        *cors.Cors
	Handler     func(w http.ResponseWriter, r *http.Request)
	Middlewares []func(handlerFunc http.HandlerFunc) http.HandlerFunc
}

func (ep Endpoint) validate() error {
	if err := ep.Method.validate(); err != nil {
		return err
	}

	if ep.Path == "" {
		return fmt.Errorf("empty path for endpoint %s", ep.Path)
	}

	if ep.Handler == nil {
		return fmt.Errorf("empty handler for endpoint %s", ep.Path)
	}

	return nil
}

func (s *Server) AddSubRouter(endpoint Endpoint) error {
	if err := endpoint.validate(); err != nil {
		return err
	}

	subRouter := chi.NewRouter()

	handler := endpoint.Handler
	for i := len(endpoint.Middlewares) - 1; i >= 0; i-- {
		handler = endpoint.Middlewares[i](handler)
	}

	if endpoint.Cors != nil {
		subRouter.Use(endpoint.Cors.Handler)
	}

	subRouter.MethodFunc(endpoint.Method.String(), "/", handler)
	s.Router.Mount(endpoint.Path, subRouter)
	return nil
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

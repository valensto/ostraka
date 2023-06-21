package http

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/valensto/ostraka/internal/middleware"
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
	Cors        *middleware.CORS
	Auth        middleware.Authenticator
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

	if endpoint.Cors != nil {
		subRouter.Use(endpoint.Cors.Init().Handler)
	}

	if endpoint.Auth != nil {
		subRouter.Use(endpoint.Auth.Register)
	}

	handler := endpoint.Handler
	for i := len(endpoint.Middlewares) - 1; i >= 0; i-- {
		handler = endpoint.Middlewares[i](handler)
	}

	subRouter.MethodFunc(endpoint.Method.String(), "/", handler)
	s.Router.Mount(endpoint.Path, subRouter)
	return nil
}

package webui

import (
	"fmt"
	"github.com/valensto/ostraka/internal/config/env"
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/provider/sse"
	"github.com/valensto/ostraka/internal/server"
	"github.com/valensto/ostraka/internal/workflow"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"strings"
)

type Webui struct {
	server *server.Server
	events chan []byte
	config env.Webui
}

func New(config env.Webui, server *server.Server, workflows []*workflow.Workflow) (*Webui, error) {
	webui := &Webui{
		server: server,
		events: make(chan []byte),
		config: config,
	}

	server.Router.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("webui/dist/assets"))))
	server.Router.Get("/admin/dashboard", webui.basicAuth(webui.dashboard()))
	server.Router.Get("/workflows", webui.workflows(workflows))

	output := workflow.WebUIOutput()
	p, err := sse.NewPublisher(output, server)
	if err != nil {
		return nil, fmt.Errorf("error getting publisher: %w", err)
	}

	logger.Get().Info().Msgf("webui running on %s:%s/admin/dashboard", webui.server.Host, webui.server.Port)
	return webui, p.Publish(webui.events)
}

func (webui *Webui) Consume(bytes []byte) {
	webui.events <- bytes
}

func (webui *Webui) basicAuth(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !webui.config.Enabled {
			handler.ServeHTTP(w, r)
			return
		}

		u, p, ok := r.BasicAuth()
		if !ok || len(strings.TrimSpace(u)) < 1 || len(strings.TrimSpace(p)) < 1 {
			unauthorised(w)
			return
		}

		logger.Get().Debug().Msgf("basic auth credentials for user %s", u)
		logger.Get().Debug().Msgf("basic auth credentials for pwd %s", p)
		if !webui.isAuth(u, p) {
			logger.Get().Error().Msgf("invalid basic auth credentials for user %s", u)
			unauthorised(w)
			return
		}

		handler.ServeHTTP(w, r)
	}
}

func (webui *Webui) isAuth(user, password string) bool {
	if user != webui.config.BasicAuthUsername {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(webui.config.BasicAuthPwd), []byte(password))
	if err != nil {
		return false
	}

	return true
}

func unauthorised(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
	w.WriteHeader(http.StatusUnauthorized)
}

func (webui *Webui) workflows(workflows []*workflow.Workflow) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		webui.server.Respond(w, r, http.StatusOK, mapWorkflowToDTO(workflows))
	}
}

func (webui *Webui) dashboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("webui/dist/index.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			webui.server.Respond(w, r, http.StatusInternalServerError, nil)
			return
		}

		webui.server.Respond(w, r, http.StatusOK, nil)
	}
}

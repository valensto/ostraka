package webui

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/valensto/ostraka/internal/config/env"
	ostraHTTP "github.com/valensto/ostraka/internal/http"
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/middleware"
	"github.com/valensto/ostraka/internal/provider/sse"
	"github.com/valensto/ostraka/internal/workflow"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"strings"
)

type Webui struct {
	server *ostraHTTP.Server
	config env.Webui

	publisher *sse.Publisher
}

func New(config env.Webui, server *ostraHTTP.Server, workflows []*workflow.Workflow) (*Webui, error) {
	mux := chi.NewRouter()
	publisher, err := sse.WebUIPublisher(config, server)
	if err != nil {
		return nil, fmt.Errorf("cannot create webui publisher: %w", err)
	}

	webui := &Webui{
		server:    server,
		config:    config,
		publisher: publisher,
	}

	cors := &middleware.CORS{
		AllowedOrigins: config.AllowedOrigins,
		AllowedMethods: []string{"GET", "POST"},
	}
	mux.Use(cors.Init().Handler)
	server.Router.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("ui/dist/assets"))))
	mux.Get("/dashboard", webui.basicAuth(webui.dashboard()))
	mux.Get("/workflows", webui.workflows(workflows))

	server.Router.Mount("/webui", mux)

	logger.Get().Info().Msgf("views running on %s:%s/webui/dashboard", webui.server.Host, webui.server.Port)
	return webui, nil
}

func (webui *Webui) Consume(entry workflow.Entry) {
	b, err := entry.JSONEncode()
	if err != nil {
		logger.Get().Error().Msgf("error encoding entry %s got: %s", entry.Id, err.Error())
		return
	}

	webui.publisher.Publish(b)
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
		tmpl := template.Must(template.ParseFiles("ui/dist/index.html"))
		err := tmpl.Execute(w, nil)
		if err != nil {
			webui.server.Respond(w, r, http.StatusInternalServerError, nil)
			return
		}

		webui.server.Respond(w, r, http.StatusOK, nil)
	}
}

package webui

import (
	"github.com/valensto/ostraka/internal/collector"
	"github.com/valensto/ostraka/internal/config/env"
	"github.com/valensto/ostraka/internal/event"
	ostraHTTP "github.com/valensto/ostraka/internal/http"
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/workflow"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"strings"
)

type Webui struct {
	server *ostraHTTP.Server
	events chan event.Payload
	config env.Webui
}

func New(config env.Webui, server *ostraHTTP.Server, workflows []*workflow.Workflow) (*Webui, error) {
	webui := &Webui{
		server: server,
		events: make(chan event.Payload),
		config: config,
	}

	server.Router.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("ui/dist/assets"))))
	server.Router.Get("/webui/dashboard", webui.basicAuth(webui.dashboard()))
	server.Router.Get("/webui/workflows", webui.workflows(workflows))

	logger.Get().Info().Msgf("views running on %s:%s/ui/dashboard", webui.server.Host, webui.server.Port)

	/*return webui, sse.WebUIPublisher(config).Publish(webui.events, server)*/
	return webui, nil
}

func (webui *Webui) Consume(event collector.Event) {
	webui.events <- event.ToMap()
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

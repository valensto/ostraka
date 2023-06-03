package server

import (
	"encoding/json"
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/workflow"
	"html/template"
	"net/http"
)

type Notifier interface {
	FullName() string
}

type workflowDTO struct {
	Name      string `json:"name"`
	NbInputs  int    `json:"nb_inputs"`
	NbOutputs int    `json:"nb_outputs"`
}

func mapWorkflowToDTO(workflows []*workflow.Workflow) []workflowDTO {
	var dtos []workflowDTO
	for _, wf := range workflows {
		dtos = append(dtos, workflowDTO{
			Name:      wf.Name,
			NbInputs:  len(wf.Inputs),
			NbOutputs: len(wf.Outputs),
		})
	}
	return dtos
}

func (s *Server) serveWebui(workflows []*workflow.Workflow) {
	// TODO: add basic auth from config or env var

	s.Router.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("webui/dist/assets"))))
	s.Router.Get("/dashboard", s.webui())
	s.Router.Get("/workflows", s.getWorkflows(workflows))

	logger.Get().Info().Msgf("webui running on %s:%s/dashboard", s.host, s.port)
}

func (s *Server) getWorkflows(workflows []*workflow.Workflow) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Respond(w, r, http.StatusOK, mapWorkflowToDTO(workflows))
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

func (s *Server) Notifications() <-chan []byte {
	return s.notifications
}

func (s *Server) NotifyWebUI(workflowName string, notifier Notifier, event []byte, err error) {
	nStatus := succeed
	if err != nil {
		nStatus = failed
	}

	var nAction action
	switch notifier.(type) {
	case *workflow.Output:
		nAction = sent
	case *workflow.Input:
		nAction = received
	default:
		logger.Get().Error().Msgf("unknown notifier: %+v cannot notify", notifier)
		return
	}

	n := notification{
		Workflow: workflowName,
		Action:   nAction,
		Notifier: notifier.FullName(),
		Event:    event,
		State:    nStatus,
		Message:  "message",
	}

	s.notifications <- n.marshall()
}

type (
	state  string
	action string
)

func (s state) String() string {
	return string(s)
}

func (a action) String() string {
	return string(a)
}

const (
	succeed state = "succeed"
	failed  state = "failed"

	received action = "received"
	sent     action = "sent"
)

type notification struct {
	Workflow string `json:"workflow"`
	Action   action `json:"action"`
	Notifier string `json:"notifier"`
	Event    []byte `json:"event"`
	State    state  `json:"state"`
	Message  string `json:"message"`
}

func (n notification) marshall() []byte {
	marshal, err := json.Marshal(n)
	if err != nil {
		logger.Get().Error().Msgf("error %s marshalling event: %+v", err.Error(), n)
		return nil
	}

	return marshal
}

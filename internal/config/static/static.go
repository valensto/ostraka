package static

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/valensto/ostraka/internal/event"
	"github.com/valensto/ostraka/internal/http"
	"github.com/valensto/ostraka/internal/middleware"
	provider "github.com/valensto/ostraka/internal/provider"
	"github.com/valensto/ostraka/internal/workflow"
	"gopkg.in/yaml.v3"
)

type ContentFile map[string][]byte

func BuildWorkflows(contentFile ContentFile, server *http.Server) ([]*workflow.Workflow, error) {
	var wfs []*workflow.Workflow
	for fn, content := range contentFile {
		var sw Workflow
		err := yaml.Unmarshal(content, &sw)
		if err != nil {
			return nil, fmt.Errorf("error parsing YAML wf: %w in file %s", err, fn)
		}

		err = validator.New().Struct(sw)
		if err != nil {
			return nil, fmt.Errorf("error: validating wf: %w in file %s", err, fn)
		}

		wf, err := sw.toWorkflow(server)
		if err != nil {
			return nil, fmt.Errorf("error: converting workflow: %w. In file: %s", err, fn)
		}

		wfs = append(wfs, wf)
	}

	return wfs, nil
}

func (sw Workflow) toWorkflow(server *http.Server) (*workflow.Workflow, error) {
	eventType, err := sw.EventType.toEventType()
	if err != nil {
		return nil, err
	}

	middlewares, err := sw.Middlewares.toMiddleware()
	if err != nil {
		return nil, err
	}

	opts := provider.Options{
		Middlewares: middlewares,
		Server:      server,
	}

	inputs := make([]*workflow.Input, len(sw.Inputs))
	for i, si := range sw.Inputs {
		inputs[i], err = si.toInput(eventType, opts)
		if err != nil {
			return nil, err
		}
	}

	outputs := make([]*workflow.Output, len(sw.Outputs))
	for i, so := range sw.Outputs {
		outputs[i], err = so.toOutput(opts)
		if err != nil {
			return nil, err
		}
	}

	return workflow.New(sw.Name, inputs, outputs)
}

func (ms Middlewares) toMiddleware() (*middleware.Middlewares, error) {
	middlewares := &middleware.Middlewares{
		HTTP: middleware.HTTP{
			CORS:           make(map[string]middleware.CORS, len(ms.CORS)),
			Authenticators: make(map[string]middleware.Authenticator, len(ms.Auth)),
		},
	}

	for _, ma := range ms.Auth {
		a, err := middleware.NewAuthentication(middleware.Auth{
			Type:   ma.Type,
			Params: ma.Params,
		})
		if err != nil {
			return nil, fmt.Errorf("error creating authenticator %s: %w", ma.Name, err)
		}

		middlewares.HTTP.Authenticators[ma.Name] = a
	}

	for _, mc := range ms.CORS {
		c, err := middleware.NewCORS(
			mc.AllowedOrigins,
			mc.AllowedMethods,
			mc.AllowedHeaders,
			mc.AllowCredentials,
			mc.MaxAge,
		)
		if err != nil {
			return nil, fmt.Errorf("error creating cors %s: %w", mc.Name, err)
		}

		middlewares.HTTP.CORS[mc.Name] = *c
	}

	return middlewares, nil
}

func (sc Condition) toCondition() (*workflow.Condition, error) {
	cs := make([]*workflow.Condition, len(sc.Conditions))
	for i, c := range sc.Conditions {
		nc, err := c.toCondition()
		if err != nil {
			return nil, fmt.Errorf("error converting condition yaml: %w", err)
		}
		cs[i] = nc
	}

	return workflow.NewCondition(sc.Field, sc.Operator, sc.Value, cs...)
}

func (se EventType) toEventType() (*event.Type, error) {
	fields := make([]event.Field, len(se.Fields))
	for i, sf := range se.Fields {
		f, err := event.UnmarshallField(sf.Name, sf.DataType, sf.Required)
		if err != nil {
			return nil, fmt.Errorf("error converting field yaml: %w", err)
		}

		fields[i] = f
	}

	return event.UnmarshallType(se.Format, fields...)
}

func (si Input) toInput(t *event.Type, opts provider.Options) (*workflow.Input, error) {
	mappers := make([]event.Mapper, len(si.Decoder.Mappers))
	for _, sm := range si.Decoder.Mappers {
		mappers = append(mappers, event.Mapper{
			Source: sm.Source,
			Target: sm.Target,
		})
	}

	decoder, err := event.UnmarshallDecoder(si.Decoder.Format, mappers, t)
	if err != nil {
		return nil, fmt.Errorf("error converting decoder yaml: %w", err)
	}

	input, err := workflow.UnmarshallInput(si.Name, si.Source, decoder, si.Params, opts)
	if err != nil {
		return nil, fmt.Errorf("error converting input yaml: %w", err)
	}

	return input, nil
}

func (so Output) toOutput(opts provider.Options) (*workflow.Output, error) {
	var condition *workflow.Condition
	if so.Condition != nil {
		c, err := so.Condition.toCondition()
		if err != nil {
			return nil, fmt.Errorf("error converting condition yaml: %w", err)
		}

		condition = c
	}

	encoder, err := event.UnmarshalEncoder()
	if err != nil {
		return nil, fmt.Errorf("error converting encoder yaml: %w", err)
	}

	output, err := workflow.UnmarshallOutput(so.Name, so.Destination, condition, encoder, so.Params, opts)
	if err != nil {
		return nil, fmt.Errorf("error converting output yaml: %w", err)
	}

	return output, nil
}

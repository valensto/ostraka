package static

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/valensto/ostraka/internal/event"
	middleware2 "github.com/valensto/ostraka/internal/middleware"
	"github.com/valensto/ostraka/internal/workflow"
	"github.com/valensto/ostraka/internal/workflow/provider"
	"gopkg.in/yaml.v3"
)

func BuildWorkflows(contentFile ContentFile) ([]*workflow.Workflow, error) {
	var wfs []*workflow.Workflow
	for fn, content := range contentFile {
		var sw workflowModel
		err := yaml.Unmarshal(content, &sw)
		if err != nil {
			return nil, fmt.Errorf("error parsing YAML wf: %w in file %s", err, fn)
		}

		err = validator.New().Struct(sw)
		if err != nil {
			return nil, fmt.Errorf("error: validating wf: %w in file %s", err, fn)
		}

		wf, err := sw.toWorkflow()
		if err != nil {
			return nil, fmt.Errorf("error: converting workflow: %w. In file: %s", err, fn)
		}

		wfs = append(wfs, wf)
	}

	return wfs, nil
}

func (sw workflowModel) toWorkflow() (*workflow.Workflow, error) {
	event, err := sw.EventType.toEvent()
	if err != nil {
		return nil, err
	}

	middlewares, err := sw.Middlewares.toMiddleware()
	if err != nil {
		return nil, err
	}

	subscribers := make([]workflow.Subscriber, len(sw.Inputs))
	for i, si := range sw.Inputs {
		subscribers[i], err = si.toSubscriber(event, middlewares)
		if err != nil {
			return nil, err
		}
	}

	publishers := make([]workflow.Publisher, len(sw.Outputs))
	for i, so := range sw.Outputs {
		publishers[i], err = so.toPublisher(middlewares)
		if err != nil {
			return nil, err
		}
	}

	return workflow.New(sw.Name, subscribers, publishers)
}

func (ms middlewares) toMiddleware() (*middleware2.Middlewares, error) {
	middlewares := &middleware2.Middlewares{
		HTTP: middleware2.HTTP{
			CORS:           make(map[string]middleware2.CORS, len(ms.CORS)),
			Authenticators: make(map[string]middleware2.Authenticator, len(ms.Auth)),
		},
	}

	for _, ma := range ms.Auth {
		a, err := middleware2.NewAuthentication(middleware2.Auth{
			Type:   ma.Type,
			Params: ma.Params,
		})
		if err != nil {
			return nil, fmt.Errorf("error creating authenticator %s: %w", ma.Name, err)
		}

		middlewares.HTTP.Authenticators[ma.Name] = a
	}

	for _, mc := range ms.CORS {
		c, err := middleware2.NewCORS(
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

func (sc condition) toCondition() (*workflow.Condition, error) {
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

func (se event) toEvent() (*event.Type, error) {
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

func (si input) toSubscriber(event *event.Type, middlewares *middleware2.Middlewares) (workflow.Subscriber, error) {
	mappers := make([]event.Mapper, len(si.Decoder.Mappers))
	for _, sm := range si.Decoder.Mappers {
		mappers = append(mappers, event.Mapper{
			Source: sm.Source,
			Target: sm.Target,
		})
	}

	decoder, err := event.UnmarshallDecoder(si.Decoder.Format, mappers, event)
	if err != nil {
		return nil, fmt.Errorf("error converting decoder yaml: %w", err)
	}

	input, err := workflow.UnmarshallInput(si.Name, si.Source, *decoder)
	if err != nil {
		return nil, fmt.Errorf("error converting input yaml: %w", err)
	}

	return provider.NewSubscriber(input, si.Params, middlewares)
}

func (so output) toPublisher(middlewares *middleware2.Middlewares) (workflow.Publisher, error) {
	var condition *workflow.Condition
	if so.Condition != nil {
		c, err := so.Condition.toCondition()
		if err != nil {
			return nil, fmt.Errorf("error converting condition yaml: %w", err)
		}

		condition = c
	}

	output, err := workflow.UnmarshallOutput(so.Name, so.Destination, condition)
	if err != nil {
		return nil, fmt.Errorf("error converting output yaml: %w", err)
	}

	return provider.NewPublisher(output, so.Params, middlewares)
}

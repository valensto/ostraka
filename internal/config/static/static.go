package static

import (
	"fmt"
	"github.com/valensto/ostraka/internal/workflow/provider"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"

	"github.com/valensto/ostraka/internal/workflow"
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
			return nil, fmt.Errorf("error converting workflowModel to workflow: %w in file %s", err, fn)
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

	subscribers := make([]workflow.Subscriber, len(sw.Inputs))
	for i, si := range sw.Inputs {
		subscribers[i], err = si.toSubscriber(event)
		if err != nil {
			return nil, err
		}
	}

	publishers := make([]workflow.Publisher, len(sw.Outputs))
	for i, so := range sw.Outputs {
		publishers[i], err = so.toPublisher()
		if err != nil {
			return nil, err
		}
	}

	return workflow.New(sw.Name, subscribers, publishers)
}

func (sc conditionModel) toCondition() (*workflow.Condition, error) {
	cs := make([]*workflow.Condition, len(sc.Conditions))
	for i, c := range sc.Conditions {
		nc, err := c.toCondition()
		if err != nil {
			return nil, err
		}
		cs[i] = nc
	}

	return workflow.NewCondition(sc.Field, sc.Operator, sc.Value, cs...)
}

func (se eventTypeModel) toEvent() (*workflow.EventType, error) {
	fields := make([]workflow.Field, len(se.Fields))
	for i, sf := range se.Fields {
		f, err := workflow.UnmarshallField(sf.Name, sf.DataType, sf.Required)
		if err != nil {
			return nil, err
		}

		fields[i] = f
	}

	return workflow.UnmarshallEventType(se.Format, fields...)
}

func (si inputModel) toSubscriber(event *workflow.EventType) (workflow.Subscriber, error) {
	mappers := make([]workflow.Mapper, len(si.Decoder.Mappers))
	for _, sm := range si.Decoder.Mappers {
		mappers = append(mappers, workflow.Mapper{
			Source: sm.Source,
			Target: sm.Target,
		})
	}

	decoder, err := workflow.UnmarshallDecoder(si.Decoder.Format, mappers)
	if err != nil {
		return nil, err
	}

	input, err := workflow.UnmarshallInput(si.Name, si.Source, *decoder, event)
	if err != nil {
		return nil, err
	}

	return provider.NewSubscriber(input, si.Params)
}

func (so outputModel) toPublisher() (workflow.Publisher, error) {
	var condition *workflow.Condition
	if so.Condition != nil {
		c, err := so.Condition.toCondition()
		if err != nil {
			return nil, err
		}

		condition = c
	}

	output, err := workflow.UnmarshallOutput(so.Name, so.Destination, condition)
	if err != nil {
		return nil, err
	}

	return provider.NewPublisher(output, so.Params)
}

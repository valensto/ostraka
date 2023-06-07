package static

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"

	"github.com/valensto/ostraka/internal/workflow"
)

func BuildWorkflows(contentFile ContentFile) ([]*workflow.Workflow, error) {
	var wfs []*workflow.Workflow
	for fname, content := range contentFile {
		var sw workflowModel
		err := yaml.Unmarshal(content, &sw)
		if err != nil {
			return nil, fmt.Errorf("error parsing YAML wf: %w in file %s", err, fname)
		}

		err = validator.New().Struct(sw)
		if err != nil {
			return nil, fmt.Errorf("error: validating wf: %w in file %s", err, fname)
		}

		wf, err := sw.toWorkflow()
		if err != nil {
			return nil, fmt.Errorf("error converting workflowModel to workflow: %w in file %s", err, fname)
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

	inputs := make([]*workflow.Input, len(sw.Inputs))
	for i, si := range sw.Inputs {
		inputs[i], err = si.toInput(event)
		if err != nil {
			return nil, err
		}
	}

	outputs := make([]*workflow.Output, len(sw.Outputs))
	for i, so := range sw.Outputs {
		outputs[i], err = so.toOutput()
		if err != nil {
			return nil, err
		}
	}

	return workflow.New(sw.Name, inputs, outputs)
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

func (si inputModel) toInput(event *workflow.EventType) (*workflow.Input, error) {
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

	return workflow.UnmarshallInput(si.Name, si.Source, *decoder, si.Params, event)
}

func (so outputModel) toOutput() (*workflow.Output, error) {
	var condition *workflow.Condition
	if so.Condition != nil {
		c, err := so.Condition.toCondition()
		if err != nil {
			return nil, err
		}

		condition = c
	}

	return workflow.UnmarshallOutput(so.Name, so.Destination, condition, so.Params)
}

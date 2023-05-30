package workflow

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Workflow struct {
	Inputs  map[string]Input
	Outputs map[string]Output
}

func New(inputs []*Input, outputs []*Output) (*Workflow, error) {
	wf := Workflow{
		Inputs:  make(map[string]Input),
		Outputs: make(map[string]Output),
	}

	for _, input := range inputs {
		wf.Inputs[input.Name] = *input
	}

	for _, output := range outputs {
		wf.Outputs[output.Name] = *output
	}

	return &wf, nil
}

type parameter interface {
	validate() error
}

func unmarshalParams(marshalled []byte, params interface{}) (err error) {
	t := reflect.TypeOf(params)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("type %T is not a pointer", params)
	}

	err = json.Unmarshal(marshalled, params)
	if err != nil {
		return fmt.Errorf("error unmarshalling params to type %T got: %w ", params, err)
	}

	return nil
}

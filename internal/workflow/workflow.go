package workflow

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"

	"github.com/valensto/ostraka/logger"
)

type Workflows []Workflow

type Workflow struct {
	Event   Event    `yaml:"event" validate:"required,dive,required"`
	Inputs  []Input  `yaml:"inputs" validate:"required,dive,required"`
	Outputs []Output `yaml:"outputs" validate:"required,dive,required"`
}

type Event struct {
	Type   string  `yaml:"type" validate:"required"`
	Fields []Field `yaml:"fields" validate:"required,dive,required"`
}

type Field struct {
	Name     string `yaml:"name" validate:"required"`
	Type     string `yaml:"type" validate:"required"`
	Required bool   `yaml:"required"`
}

func Build() (Workflows, error) {
	directory := ".ostraka/workflows"
	dir, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("error reading resources directory: %w", err)
	}

	var workflows Workflows
	for _, file := range dir {
		ext := filepath.Ext(file.Name())
		if ext != ".yaml" && ext != ".yml" {
			logger.Get().Warn().Msgf(`unable to find .yaml or .yml file. '%s' will be skipped`, file.Name())
			continue
		}

		wf, err := extractWorkflow(filepath.Join(directory, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("error extracting workflow: %w", err)
		}
		workflows = append(workflows, wf)
	}

	return workflows, nil
}

func extractWorkflow(filename string) (Workflow, error) {
	f, err := os.Open(filename)
	if err != nil {
		return Workflow{}, fmt.Errorf("error opening file %s: %w", filename, err)
	}
	defer f.Close()

	workflow, err := readConfig(f)
	if err != nil {
		return Workflow{}, fmt.Errorf("error parsing f: %w", err)
	}

	return workflow, nil
}

func readConfig(r io.Reader) (Workflow, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return Workflow{}, fmt.Errorf("error reading ostrakaflow wf: %w", err)
	}

	var wf Workflow
	err = yaml.Unmarshal(b, &wf)
	if err != nil {
		return Workflow{}, fmt.Errorf("error parsing YAML wf: %w", err)
	}

	err = wf.setInputs()
	if err != nil {
		return Workflow{}, fmt.Errorf("unable to set inputs: %w", err)
	}

	err = wf.setOutputs()
	if err != nil {
		return Workflow{}, fmt.Errorf("unable to set outputs: %w", err)
	}

	err = validator.New().Struct(wf)
	if err != nil {
		return Workflow{}, fmt.Errorf("error validating wf: %w", err)
	}

	return wf, nil
}

func unmarshalParams(marshalled []byte, params interface{}) (err error) {
	t := reflect.TypeOf(params)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("type %T is not a pointer", params)
	}

	err = yaml.Unmarshal(marshalled, params)
	if err != nil {
		return fmt.Errorf("error unmarshalling params to type %T got: %w ", params, err)
	}

	err = validator.New().Struct(params)
	if err != nil {
		return fmt.Errorf("error validating params: %w", err)
	}

	return nil
}

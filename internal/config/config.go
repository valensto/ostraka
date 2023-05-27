package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var validate = validator.New()

type Workflows []Workflow

type Workflow struct {
	Inputs  []Input  `yaml:"inputs" validate:"required,dive,required"`
	Event   Event    `yaml:"event" validate:"required,dive,required"`
	Outputs []Output `yaml:"outputs" validate:"required,dive,required"`
}

type Field struct {
	Name     string `yaml:"name" validate:"required"`
	Type     string `yaml:"type" validate:"required"`
	Required bool   `yaml:"required"`
}

func NewWorkflow() (Workflows, error) {
	directory := "resources"
	dir, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("error reading resources directory: %w", err)
	}

	var workflows Workflows
	for _, file := range dir {
		ext := filepath.Ext(file.Name())
		if ext != ".yaml" && ext != ".yml" {
			log.Warningf(`unable to find .yaml or .yml file. "%s" will be skipped`, file.Name())
			continue
		}

		f, err := os.Open(filepath.Join(directory, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("error opening file %s: %w", file, err)
		}
		defer f.Close()

		workflow, err := readConfig(f)
		if err != nil {
			return nil, fmt.Errorf("error parsing f: %w", err)
		}

		workflows = append(workflows, *workflow)
	}

	return workflows, nil
}

func readConfig(r io.Reader) (*Workflow, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("error reading ostrakaflow file: %w", err)
	}

	var file Workflow
	err = yaml.Unmarshal(b, &file)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML file: %w", err)
	}

	err = file.setInputs()
	if err != nil {
		return nil, fmt.Errorf("unable to set inputs: %w", err)
	}

	err = file.populateOutputs()
	if err != nil {
		return nil, fmt.Errorf("unable to set outputs: %w", err)
	}

	err = validate.Struct(file)
	if err != nil {
		return nil, fmt.Errorf("error validating workflow: %w", err)
	}

	return &file, nil
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

	return validate.Struct(params)
}

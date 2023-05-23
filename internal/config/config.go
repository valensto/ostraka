package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"reflect"
)

var validate = validator.New()

type Config []File

type File struct {
	Inputs  []Input  `yaml:"inputs" validate:"required,dive,required"`
	Event   Event    `yaml:"event" validate:"required,dive,required"`
	Outputs []Output `yaml:"outputs" validate:"required,dive,required"`
}

type Field struct {
	Name     string `yaml:"name" validate:"required"`
	Type     string `yaml:"type" validate:"required"`
	Required bool   `yaml:"required"`
}

func LoadConfig() (Config, error) {
	directory := "resources"

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	var config []File
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".yaml" {
			filePath := filepath.Join(directory, f.Name())

			parsedFile, err := parseFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("error parsing f: %w", err)
			}

			config = append(config, *parsedFile)
		}
	}

	return config, nil
}

func parseFile(path string) (*File, error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %w", err)
	}

	var file File
	err = yaml.Unmarshal(yamlFile, &file)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML file: %w", err)
	}

	err = file.populateInputs()
	if err != nil {
		return nil, fmt.Errorf("error populating inputs: %w", err)
	}

	err = file.populateOutputs()
	if err != nil {
		return nil, fmt.Errorf("error populating outputs: %w", err)
	}

	return &file, validate.Struct(file)
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

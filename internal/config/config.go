package config

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
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
	directory := "internal/config/resources"

	var config []File
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".yaml") || strings.HasSuffix(info.Name(), ".yml")) {
			filePath := filepath.Join(directory, info.Name())

			parsedFile, err := parseFile(filePath)
			if err != nil {
				return fmt.Errorf("error parsing f: %w", err)
			}

			config = append(config, *parsedFile)
		}

		return nil
	})
	if err != nil {
		fmt.Println("Error reading directory:", err)
	}

	if len(config) == 0 {
		return nil, fmt.Errorf("no config files found in %s", directory)
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

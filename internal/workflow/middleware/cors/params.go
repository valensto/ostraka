package cors

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type Params struct {
	Name             string   `yaml:"name"`
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

func (c *Params) validate() error {
	if c.Name == "" {
		return fmt.Errorf("cors name is empty")
	}

	if c.AllowedOrigins == nil {
		return fmt.Errorf("cors allowed_origins is empty")
	}

	if c.AllowedMethods == nil {
		return fmt.Errorf("cors allowed_methods is empty")
	}

	if c.AllowedHeaders == nil {
		return fmt.Errorf("cors allowed_headers is empty")
	}

	return nil
}

func UnmarshalCORSConfig(marshalled []byte) (Params, error) {
	c := Params{}
	err := yaml.Unmarshal(marshalled, &c)
	if err != nil {
		return Params{}, fmt.Errorf("error unmarshalling params to type CORSConfig got: %w ", err)
	}

	err = c.validate()
	if err != nil {
		return Params{}, err
	}

	return c, nil
}

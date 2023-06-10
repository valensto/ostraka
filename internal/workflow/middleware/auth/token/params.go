package token

import (
	"fmt"
	"gopkg.in/yaml.v3"
)

type Params struct {
	Token      string `yaml:"token"`
	QueryParam string `yaml:"query_param"`
}

func (t *Params) validate() error {
	if t.Token == "" {
		return fmt.Errorf("token is empty")
	}

	if t.QueryParam == "" {
		return fmt.Errorf("query_param is empty")
	}

	return nil
}

func (t *Params) Unmarshal(marshalled []byte) (err error) {
	err = yaml.Unmarshal(marshalled, t)
	if err != nil {
		return fmt.Errorf("error unmarshalling params to type Token got: %w ", err)
	}

	return t.validate()
}

func UnmarshalTokenParams(marshalled []byte) (Params, error) {
	t := Params{}
	err := t.Unmarshal(marshalled)
	if err != nil {
		return Params{}, err
	}

	return t, nil
}

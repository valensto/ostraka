package smtp

import (
	"encoding/json"
	"fmt"
)

type Params struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`

	BaseURL        string `json:"base_url"`
	Host           string `json:"host"`
	Port           string `json:"port"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	EnableStartTLS bool   `json:"enable_starttls"`
}

func (p *Params) validate() error {
	if p.From == "" {
		return fmt.Errorf("smtp no_reply is empty")
	}

	if p.BaseURL == "" {
		return fmt.Errorf("smtp base_url is empty")
	}

	if p.Host == "" {
		return fmt.Errorf("smtp host is empty")
	}

	if p.Port == "" {
		return fmt.Errorf("smtp port is empty")
	}

	if p.Username == "" {
		return fmt.Errorf("smtp username is empty")
	}

	if p.Password == "" {
		return fmt.Errorf("smtp password is empty")
	}

	return nil
}

func unmarshalParams(bytes []byte) (*Params, error) {
	sse := Params{}
	err := json.Unmarshal(bytes, &sse)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling params to type SMTP got: %w ", err)
	}

	err = sse.validate()
	if err != nil {
		return nil, err
	}

	return &sse, nil
}

package jwt

import (
	"encoding/json"
	"fmt"
)

type Params struct {
	Header           string  `json:"header"`
	Secret           string  `json:"secret"`
	Algorithm        string  `json:"algorithm"`
	VerifyExpiration bool    `json:"verify_expiration"`
	MaxAge           int     `json:"max_age"`
	Payload          []Field `json:"payload"`
}

type Field struct {
	Name     string `json:"name"`
	DataType string `json:"data_type"`
	Required bool   `json:"required"`
}

func (jwt *Params) validate() error {
	if jwt.Header == "" {
		return fmt.Errorf("jwt header is empty")
	}

	if jwt.Secret == "" {
		return fmt.Errorf("jwt secret is empty")
	}

	if jwt.Algorithm == "" {
		return fmt.Errorf("jwt algorithm is empty")
	}

	if jwt.MaxAge == 0 {
		return fmt.Errorf("jwt max_age is empty")
	}

	if jwt.Payload == nil {
		return fmt.Errorf("jwt payload is empty")
	}

	return nil
}

func UnmarshalJWT(marshalled []byte) (Params, error) {
	jwt := Params{}
	err := json.Unmarshal(marshalled, &jwt)
	if err != nil {
		return Params{}, fmt.Errorf("error unmarshalling params to type JWT got: %w ", err)
	}

	err = jwt.validate()
	if err != nil {
		return Params{}, err
	}

	return jwt, nil
}

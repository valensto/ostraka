package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Token struct {
	Token      string `json:"token" yaml:"token" validate:"required"`
	QueryParam string `json:"query_param" yaml:"query_param" validate:"required"`
}

func (config Token) validate() error {
	if config.Token == "" {
		return fmt.Errorf("token is empty")
	}

	if config.QueryParam == "" {
		return fmt.Errorf("query_param is empty")
	}

	return nil
}

func unmarshalToken(bytes []byte) (*Token, error) {
	config := Token{}
	err := json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling params to type JWT got: %w ", err)
	}

	err = config.validate()
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (config Token) Register(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			token := r.URL.Query().Get(config.QueryParam)
			if token == "" {
				http.Error(w, fmt.Sprintf("token query param %s is empty", config.QueryParam), http.StatusUnauthorized)
				return
			}

			if token != config.Token {
				http.Error(w, fmt.Sprintf("token query param %s is invalid", config.QueryParam), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		},
	)
}

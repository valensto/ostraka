package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

type JWT struct {
	Header           string  `json:"header"`
	Secret           string  `json:"secret"`
	Algorithm        string  `json:"algorithm"`
	VerifyExpiration bool    `json:"verify_expiration"`
	MaxAge           int     `json:"max_age"`
	Payload          []field `json:"payload"`
}

type field struct {
	Name     string `json:"name"`
	DataType string `json:"data_type"`
	Required bool   `json:"required"`
}

func (config JWT) validate() error {
	if config.Header == "" {
		return fmt.Errorf("jwt header is empty")
	}

	if config.Secret == "" {
		return fmt.Errorf("jwt secret is empty")
	}

	if config.Algorithm == "" {
		return fmt.Errorf("jwt algorithm is empty")
	}

	if config.MaxAge == 0 {
		return fmt.Errorf("jwt max_age is empty")
	}

	if config.Payload == nil {
		return fmt.Errorf("jwt payload is empty")
	}

	return nil
}

func unmarshalJWT(bytes []byte) (*JWT, error) {
	config := JWT{}
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

func (config JWT) Register(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get(config.Header)
		if authorizationHeader == "" {
			http.Error(w, fmt.Sprintf("authorization %s is empty", config.Header), http.StatusUnauthorized)
			return
		}

		split := strings.Split(authorizationHeader, "Bearer ")
		if len(split) != 2 {
			http.Error(w, "authorization header is invalid missing bearer", http.StatusUnauthorized)
			return
		}

		str := strings.TrimSpace(split[1])
		token, err := jwt.Parse(str, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("signing method is invalid")
			}

			return config.Secret, nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "token is invalid", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

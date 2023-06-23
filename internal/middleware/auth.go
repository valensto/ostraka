package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	JWTType   = "jwt"
	TokenType = "token"
)

type Auth struct {
	Type   string `json:"type" yaml:"type" validate:"required"`
	Params any    `json:"params" yaml:"params" validate:"required"`

	Authenticator Authenticator `json:"-" yaml:"-"`
}

type Authenticator interface {
	Register(next http.Handler) http.Handler
}

func loadAuthenticator(auth Auth) (Authenticator, error) {
	b, err := json.Marshal(auth.Params)
	if err != nil {
		return nil, fmt.Errorf("error marshalling authenticator params: %w", err)
	}

	switch auth.Type {
	case JWTType:
		return unmarshalJWT(b)

	case TokenType:
		return unmarshalToken(b)

	default:
		return nil, fmt.Errorf("unknown authenticator type: %s", auth.Type)
	}
}

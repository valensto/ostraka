package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	JWT   = "jwt"
	Token = "token"
)

type Auth struct {
	Type   string `json:"type"`
	Params any    `json:"params"`
}

type Authenticator interface {
	Register(next http.Handler) http.Handler
}

func NewAuthentication(auth Auth) (Authenticator, error) {
	b, err := json.Marshal(auth.Params)
	if err != nil {
		return nil, fmt.Errorf("error marshalling authenticator params: %w", err)
	}

	switch auth.Type {
	case JWT:
		return unmarshalJWT(b)

	case Token:
		return unmarshalToken(b)

	default:
		return nil, fmt.Errorf("unknown authenticator type: %s", auth.Type)
	}
}

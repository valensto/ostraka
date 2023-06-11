package middleware

import (
	"fmt"
)

type Middlewares struct {
	Web Web
}

type Web struct {
	CORS           map[string]CORS
	Authenticators map[string]Authenticator
}

func (w Web) GetAuthenticator(name string) (Authenticator, error) {
	if a, ok := w.Authenticators[name]; ok {
		return a, nil
	}

	return nil, fmt.Errorf("unknown authenticator: %s", name)
}

func (w Web) GetCORS(name string) (*CORS, error) {
	if c, ok := w.CORS[name]; ok {
		return &c, nil
	}

	return nil, fmt.Errorf("unknown cors: %s", name)
}

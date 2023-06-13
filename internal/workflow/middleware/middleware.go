package middleware

import (
	"fmt"
)

type Middlewares struct {
	HTTP HTTP
}

type HTTP struct {
	Authenticators map[string]Authenticator
	CORS           map[string]CORS
}

func (w HTTP) Authenticator(name string) (Authenticator, error) {
	if a, ok := w.Authenticators[name]; ok {
		return a, nil
	}

	return nil, fmt.Errorf("unknown authenticator: %s", name)
}

func (w HTTP) Cors(name string) (*CORS, error) {
	if c, ok := w.CORS[name]; ok {
		return &c, nil
	}

	return nil, fmt.Errorf("unknown cors: %s", name)
}

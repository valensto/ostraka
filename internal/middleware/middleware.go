package middleware

import "fmt"

type Middlewares struct {
	CORS map[string]*CORS `json:"cors" yaml:"cors"`
	Auth map[string]*Auth `json:"auth" yaml:"auth"`
}

func (m *Middlewares) LoadAuthenticators() error {
	for name, auth := range m.Auth {
		a, err := loadAuthenticator(*auth)
		if err != nil {
			return fmt.Errorf("error loading authenticator %s: %w", name, err)
		}

		m.Auth[name].Authenticator = a
	}

	return nil
}

func (m *Middlewares) Authenticator(name string) (Authenticator, error) {
	if a, ok := m.Auth[name]; ok {
		return a.Authenticator, nil
	}

	return nil, fmt.Errorf("unknown authenticator: %s", name)
}

func (m *Middlewares) Cors(name string) (*CORS, error) {
	if c, ok := m.CORS[name]; ok {
		return c, nil
	}

	return nil, fmt.Errorf("unknown cors: %s", name)
}

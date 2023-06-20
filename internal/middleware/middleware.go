package middleware

type Middlewares struct {
	CORS map[string]*CORS `json:"cors" yaml:"cors"`
	Auth map[string]*Auth `json:"auth" yaml:"auth"`
}

/*type HTTP struct {
	CORS map[string]*CORS `json:"cors" yaml:"cors"`
	Auth map[string]*Auth `json:"auth" yaml:"auth"`
}
*/
/*func (m Middlewares) Authenticator(name string) (Authenticator, error) {
	if a, ok := m.Authenticators[name]; ok {
		return a, nil
	}

	return nil, fmt.Errorf("unknown authenticator: %s", name)
}

func (m Middlewares) Cors(name string) (*CORS, error) {
	if c, ok := m.CORS[name]; ok {
		return &c, nil
	}

	return nil, fmt.Errorf("unknown cors: %s", name)
}

*/

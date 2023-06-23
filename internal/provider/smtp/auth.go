package smtp

import (
	"fmt"
	"net/smtp"
)

func authenticate(username string, password string, host string) smtp.Auth {
	if username == "" && password == "" {
		return nil
	}
	return AgnosticAuth("", username, password, host)
}

type loginAuth struct {
	username, password string
}

// LoginAuth returns an Auth that implements the LOGIN authentication which is still used by some SMTP server
// https://datatracker.ietf.org/doc/html/draft-murchison-sasl-login-00
func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(_ *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		}
	}
	return nil, nil
}

type agnosticAuth struct {
	identity string
	username string
	password string
	host     string
	auth     smtp.Auth
}

func (a *agnosticAuth) createAuth(mode string) smtp.Auth {
	switch mode {
	case "LOGIN":
		return LoginAuth(a.username, a.password)
	case "PLAIN":
		return smtp.PlainAuth(a.identity, a.username, a.password, a.host)
	case "CRAM-MD5":
		return smtp.CRAMMD5Auth(a.username, a.password)
	default:
		return nil
	}
}

func AgnosticAuth(identity, username, password, host string) smtp.Auth {
	return &agnosticAuth{identity, username, password, host, nil}
}

func (a *agnosticAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	for _, auth := range server.Auth {
		a.auth = a.createAuth(auth)
		if a.auth == nil {
			continue
		}

		proto, toServer, err := a.auth.Start(server)
		if err == nil {
			return proto, toServer, err
		}
	}

	return "", nil, fmt.Errorf("no supported authentication method found")
}

func (a *agnosticAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	return a.auth.Next(fromServer, more)
}

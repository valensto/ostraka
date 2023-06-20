package middleware

import (
	"fmt"
	"github.com/go-chi/cors"
)

type CORS struct {
	AllowedOrigins   []string `json:"allowed_origins" yaml:"allowed_origins"`
	AllowedMethods   []string `json:"allowed_methods" yaml:"allowed_methods"`
	AllowedHeaders   []string `json:"allowed_headers" yaml:"allowed_headers"`
	AllowCredentials bool     `json:"allow_credentials" yaml:"allow_credentials"`
	MaxAge           int      `json:"max_age" yaml:"max_age"`
}

func NewCORS(allowedOrigins, allowedMethods, allowedHeaders []string, allowCredentials bool, maxAge int) (*CORS, error) {
	c := &CORS{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   allowedMethods,
		AllowedHeaders:   allowedHeaders,
		AllowCredentials: allowCredentials,
		MaxAge:           maxAge,
	}

	err := c.validate()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *CORS) validate() error {
	if c.AllowedOrigins == nil {
		return fmt.Errorf("cors allowed_origins is empty")
	}

	if c.AllowedMethods == nil {
		return fmt.Errorf("cors allowed_methods is empty")
	}

	if c.AllowedHeaders == nil {
		return fmt.Errorf("cors allowed_headers is empty")
	}

	return nil
}

func (c *CORS) Init() *cors.Cors {
	if c.AllowedOrigins == nil {
		c.AllowedOrigins = []string{"*"}
	}

	if c.AllowedMethods == nil {
		c.AllowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	}

	if c.AllowedHeaders == nil {
		c.AllowedHeaders = []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}
	}

	if c.MaxAge == 0 {
		c.MaxAge = 300
	}

	return cors.New(cors.Options{
		AllowedOrigins:   c.AllowedOrigins,
		AllowedMethods:   c.AllowedMethods,
		AllowedHeaders:   c.AllowedHeaders,
		AllowCredentials: c.AllowCredentials,
		MaxAge:           c.MaxAge,
	})
}

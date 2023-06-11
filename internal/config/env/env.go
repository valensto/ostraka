package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Host  string
	Port  string
	Webui Webui
}

type Webui struct {
	Enabled           bool
	BasicAuthPwd      string
	BasicAuthUsername string
	AllowedOrigins    []string
	AuthToken         string
}

func Load() (*Config, error) {
	host := os.Getenv("HOST")
	if host == "" {
		return nil, fmt.Errorf("HOST env var is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		return nil, fmt.Errorf("PORT env var is required")
	}

	webui, err := loadWebui()
	if err != nil {
		return nil, err
	}

	return &Config{
		Host:  host,
		Port:  port,
		Webui: webui,
	}, nil
}

func loadWebui() (Webui, error) {
	enabled, err := strconv.ParseBool(os.Getenv("WEBUI_ENABLED"))
	if err != nil {
		return Webui{}, fmt.Errorf("cannot parse WEBUI_ENABLED env var to boolean: %w", err)
	}

	if !enabled {
		return Webui{}, nil
	}

	allowedOrigins := os.Getenv("WEBUI_SSE_ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		return Webui{}, fmt.Errorf("WEBUI_SSE_ALLOWED_ORIGINS env var is required")
	}
	aos := strings.Split(allowedOrigins, ",")
	if len(aos) == 0 {
		return Webui{}, fmt.Errorf("WEBUI_SSE_ALLOWED_ORIGINS env var is required")
	}

	token := os.Getenv("WEBUI_SSE_TOKEN")
	if token == "" {
		return Webui{}, fmt.Errorf("WEBUI_SSE_TOKEN env var is required")
	}

	webui := Webui{
		Enabled:        enabled,
		AuthToken:      token,
		AllowedOrigins: aos,
	}

	basicAuth := os.Getenv("WEBUI_BASIC_AUTH")
	if basicAuth != "" {
		parsed := strings.Split(basicAuth, ":")
		if len(parsed) != 2 {
			return Webui{}, fmt.Errorf("WEBUI_BASIC_AUTH env var is not valid")
		}

		webui.BasicAuthUsername = parsed[0]
		webui.BasicAuthPwd = parsed[1]
	}

	return webui, nil
}

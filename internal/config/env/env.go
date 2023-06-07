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

	webui := Webui{
		Enabled:           enabled,
		BasicAuthPwd:      "",
		BasicAuthUsername: "",
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

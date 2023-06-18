package env

import (
	"github.com/valensto/ostraka/internal/logger"
	"os"
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

func Load() *Config {
	return &Config{
		Host:  mustGet("HOST"),
		Port:  mustGet("PORT"),
		Webui: loadWebui(),
	}
}

func loadWebui() Webui {
	enabled := mustGet("WEBUI_ENABLED") == "true"
	if !enabled {
		return Webui{}
	}

	allowedOrigins := mustGet("WEBUI_SSE_ALLOWED_ORIGINS")
	aos := strings.Split(allowedOrigins, ",")
	if len(aos) == 0 {
		logger.Get().Fatal().Msgf("WEBUI_SSE_ALLOWED_ORIGINS env var must be a comma separated list")
	}

	webui := Webui{
		Enabled:        enabled,
		AuthToken:      mustGet("WEBUI_SSE_TOKEN"),
		AllowedOrigins: aos,
	}

	basicAuth := mustGet("WEBUI_BASIC_AUTH")
	if basicAuth != "" {
		parsed := strings.Split(basicAuth, ":")
		if len(parsed) != 2 {
			logger.Get().Fatal().Msgf("WEBUI_BASIC_AUTH env var must be a colon separated string")
		}

		webui.BasicAuthUsername = parsed[0]
		webui.BasicAuthPwd = parsed[1]
	}

	return webui
}

func mustGet(key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Get().Fatal().Msgf("env var %s is required", key)
	}

	return value
}

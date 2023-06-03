package cors

import (
	"github.com/go-chi/cors"
)

func New(allowedOrigins, allowedMethods, allowedHeaders, exposedHeaders []string, allowCredentials bool, maxAge int) *cors.Cors {
	if allowedOrigins == nil {
		allowedOrigins = []string{"*"}
	}
	if allowedMethods == nil {
		allowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	}
	if allowedHeaders == nil {
		allowedHeaders = []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}
	}
	if exposedHeaders == nil {
		exposedHeaders = []string{"Link"}
	}
	if maxAge == 0 {
		maxAge = 300
	}

	return cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   allowedMethods,
		AllowedHeaders:   allowedHeaders,
		ExposedHeaders:   exposedHeaders,
		AllowCredentials: allowCredentials,
		MaxAge:           maxAge,
	})
}

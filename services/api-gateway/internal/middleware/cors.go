package middleware

import (
	"net/http"

	"github.com/go-chi/cors"
)

func CORS(origins []string) func(next http.Handler) http.Handler {
	opts := cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}
	return cors.New(opts).Handler
}

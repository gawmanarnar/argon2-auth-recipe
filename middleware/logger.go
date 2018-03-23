package middleware

import (
	"net/http"

	"github.com/kevinburke/handlers"
)

// Logger - enables logging middleware
func Logger(next http.Handler) http.Handler {
	return handlers.Log(next)
}

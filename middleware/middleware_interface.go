package middleware

import "net/http"

type MiddlewareInterface interface {
	Apply(next http.Handler) http.Handler
}

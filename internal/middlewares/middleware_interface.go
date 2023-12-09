package middlewares

import "net/http"

type MiddlewareInterface interface {
	Apply(next http.Handler) http.Handler
}

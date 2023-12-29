package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"

	"github.com/ranta0/rest-and-go/domain/auth"
	"github.com/ranta0/rest-and-go/error"
	"github.com/ranta0/rest-and-go/response"
)

type JWTMiddleware struct {
	tokenService *auth.RevokedJWTTokenService
}

func NewJWTMiddleware(tokenService *auth.RevokedJWTTokenService) *JWTMiddleware {
	return &JWTMiddleware{
		tokenService: tokenService,
	}
}

// JWT validation middleware for protected routes
func (m *JWTMiddleware) validateTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			response.Error(w, r, http.StatusUnauthorized, error.ErrUnauthorized.Error())
			return
		}

		// Strip the "Bearer " prefix if present
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		token, err := m.tokenService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			response.Error(w, r, http.StatusUnauthorized, error.ErrTokenInvalid.Error())
			return
		}

		if m.tokenService.IsTokenRevoked(tokenString) {
			response.Error(w, r, http.StatusUnauthorized, error.ErrTokenRevoked.Error())
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Error(w, r, http.StatusUnauthorized, error.ErrTokenInvalid.Error())
			return
		}

		// Check the token type
		tokenType, ok := claims["token_type"].(string)
		if !ok || tokenType != m.tokenService.Strings()["access"] {
			response.Error(w, r, http.StatusUnauthorized, error.ErrTokenType.Error())
			return
		}

		// Token is valid, proceed with the next middleware
		next.ServeHTTP(w, r)
	})
}

// JWT validation middleware Wrapper
func (m *JWTMiddleware) Apply(next http.Handler) http.Handler {
	return m.validateTokenMiddleware(next)
}

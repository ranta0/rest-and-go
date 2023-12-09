package middlewares

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"github.com/ranta0/rest-and-go/internal/domain/auth"
	httpErrors "github.com/ranta0/rest-and-go/internal/errors"
	"github.com/ranta0/rest-and-go/internal/utils"
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
			utils.ErrorJsonResponse(w, r, http.StatusUnauthorized, httpErrors.ErrUnauthorized.Error())
			return
		}

		// Strip the "Bearer " prefix if present
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		token, err := m.tokenService.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			utils.ErrorJsonResponse(w, r, http.StatusUnauthorized, httpErrors.ErrTokenInvalid.Error())
			return
		}

		if m.tokenService.IsTokenRevoked(tokenString) {
			utils.ErrorJsonResponse(w, r, http.StatusUnauthorized, httpErrors.ErrTokenRevoked.Error())
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.ErrorJsonResponse(w, r, http.StatusUnauthorized, httpErrors.ErrTokenInvalid.Error())
			return
		}

		// Check the token type
		tokenType, ok := claims["token_type"].(string)
		if !ok || tokenType != m.tokenService.Strings()["access"] {
			utils.ErrorJsonResponse(w, r, http.StatusUnauthorized, httpErrors.ErrTokenType.Error())
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

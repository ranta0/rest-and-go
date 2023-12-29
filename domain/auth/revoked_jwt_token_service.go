package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/ranta0/rest-and-go/config"
)

type RevokedJWTTokenService struct {
	tokenRepo             *RevokedJWTTokenRepository
	signingKey            []byte
	strings               map[string]string
	expirationTime        time.Duration
	expirationTimeRefresh time.Duration
}

func NewRevokedJWTTokenService(tokenRepo *RevokedJWTTokenRepository, cfg *config.Config) *RevokedJWTTokenService {
	return &RevokedJWTTokenService{
		tokenRepo:             tokenRepo,
		signingKey:            []byte(cfg.JWTSigningKey),
		expirationTime:        cfg.JWTExpiration,
		expirationTimeRefresh: cfg.JWTExpirationRefresh,
		strings: map[string]string{
			"access":  "access_token_type",
			"refresh": "refresh_token_type",
		},
	}
}

func (r *RevokedJWTTokenService) GenerateToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(r.signingKey)
}

func (r *RevokedJWTTokenService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return r.signingKey, nil
	})
}

func (s *RevokedJWTTokenService) IsTokenRevoked(tokenString string) bool {
	return s.tokenRepo.IsTokenRevoked(tokenString)
}

func (s *RevokedJWTTokenService) RevokeToken(tokenString string) error {
	err := s.tokenRepo.RevokeToken(tokenString)
	if err != nil {
		return err
	}

	return nil
}

func (r *RevokedJWTTokenService) Strings() map[string]string {
	return r.strings
}

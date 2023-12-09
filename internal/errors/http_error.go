package error

import (
	"errors"
)

var (
	ErrNotFound            = errors.New("not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrBadRequest          = errors.New("bad request")
	ErrContentType         = errors.New("unsupported content-type")
	ErrPayload             = errors.New("invalid request payload")
	ErrCredentials         = errors.New("invalid credentials")
	ErrTokenCreate         = errors.New("failed to create token")
	ErrTokenInvalid        = errors.New("invalid or expired token")
	ErrTokenRevoked        = errors.New("revoked token")
	ErrTokenRevokedFailure = errors.New("failed revoking token")
	ErrTokenType           = errors.New("wrong token type")
)

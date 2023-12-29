package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/ranta0/rest-and-go/domain/user"
	"github.com/ranta0/rest-and-go/error"
	"github.com/ranta0/rest-and-go/form"
	"github.com/ranta0/rest-and-go/request"
	"github.com/ranta0/rest-and-go/response"
)

type JWTController struct {
	userService  *user.UserService
	tokenService *RevokedJWTTokenService
}

func NewJWTController(userService *user.UserService, tokenService *RevokedJWTTokenService) *JWTController {
	return &JWTController{
		userService:  userService,
		tokenService: tokenService,
	}
}

func (c *JWTController) Register(w http.ResponseWriter, r *http.Request) {
	auth := &form.Auth{}

	err := request.Handler(r, auth)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = request.Validator(auth)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = c.userService.Register(auth.Username, auth.Password)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	stub := response.NewOK(nil, nil)
	response.OK(w, r, http.StatusCreated, stub)
}

func (c *JWTController) Login(w http.ResponseWriter, r *http.Request) {
	auth := &form.Auth{}

	err := request.Handler(r, auth)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = request.Validator(auth)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.userService.Login(auth.Username, auth.Password)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, error.ErrCredentials.Error())
		return
	}

	expirationTime := time.Now().Add(c.tokenService.expirationTime)
	expirationTimeRefresh := time.Now().Add(time.Hour * c.tokenService.expirationTimeRefresh)

	claims := jwt.MapClaims{
		"username": user.Username,
		"id":       user.PublicID,
	}

	claims["token_type"] = c.tokenService.strings["access"]
	claims["exp"] = expirationTime.Unix()
	accessToken, err := c.tokenService.GenerateToken(claims)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, error.ErrTokenCreate.Error())
		return
	}

	claims["token_type"] = c.tokenService.strings["refresh"]
	claims["exp"] = expirationTimeRefresh.Unix()
	refreshToken, err := c.tokenService.GenerateToken(claims)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, error.ErrTokenCreate.Error())
		return
	}

	response.JsonJWT(w, r, http.StatusOK, &response.JWT{
		Status:         "success",
		Message:        "user athenticated successfully",
		Type:           "bearer",
		AccessToken:    accessToken,
		ExpirationTime: expirationTime.Format(time.RFC3339),
		RefreshToken:   refreshToken,
	})
}

func (c *JWTController) Refresh(w http.ResponseWriter, r *http.Request) {
	auth := &form.AuthRefresh{}

	err := request.Handler(r, auth)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = request.Validator(auth)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	token, err := c.tokenService.ValidateToken(auth.RefreshToken)
	if err != nil || !token.Valid {
		response.Error(w, r, http.StatusUnauthorized, error.ErrTokenInvalid.Error())
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		response.Error(w, r, http.StatusUnauthorized, error.ErrTokenInvalid.Error())
		return
	}

	// Check the token type
	if claims["token_type"] != c.tokenService.strings["refresh"] {
		response.Error(w, r, http.StatusUnauthorized, error.ErrTokenType.Error())
		return
	}

	if c.tokenService.IsTokenRevoked(auth.RefreshToken) {
		response.Error(w, r, http.StatusUnauthorized, error.ErrTokenRevoked.Error())
		return
	}

	// Revoke previous token
	err = c.tokenService.RevokeToken(auth.RefreshToken)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, error.ErrTokenRevokedFailure.Error())
		return
	}

	expirationTime := time.Now().Add(c.tokenService.expirationTime)
	expirationTimeRefresh := time.Now().Add(time.Hour * c.tokenService.expirationTimeRefresh)

	// Generate a new access token
	claims["token_type"] = c.tokenService.strings["access"]
	claims["exp"] = expirationTime.Unix()
	newAccessToken, err := c.tokenService.GenerateToken(claims)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, error.ErrTokenCreate.Error())
		return
	}

	// Generate a new refresh token
	claims["token_type"] = c.tokenService.strings["refresh"]
	claims["exp"] = expirationTimeRefresh.Unix()
	newRefreshToken, err := c.tokenService.GenerateToken(claims)
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, error.ErrTokenCreate.Error())
		return
	}

	response.JsonJWT(w, r, http.StatusOK, &response.JWT{
		Status:         "success",
		Message:        "token refreshed successfully",
		Type:           "bearer",
		AccessToken:    newAccessToken,
		ExpirationTime: expirationTime.Format(time.RFC3339),
		RefreshToken:   newRefreshToken,
	})
}

func (c *JWTController) Logout(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		response.Error(w, r, http.StatusUnauthorized, error.ErrUnauthorized.Error())
		return
	}

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, err := c.tokenService.ValidateToken(tokenString)
	if err != nil || !token.Valid {
		response.Error(w, r, http.StatusUnauthorized, error.ErrTokenInvalid.Error())
		return
	}

	err = c.tokenService.RevokeToken(tokenString)
	if err != nil {
		response.Error(w, r, http.StatusUnauthorized, error.ErrTokenRevokedFailure.Error())
		return
	}

	stub := response.NewOK(nil, nil)
	response.OK(w, r, http.StatusOK, stub)
}

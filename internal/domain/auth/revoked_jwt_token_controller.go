package auth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/ranta0/rest-and-go/internal/domain/user"
	"github.com/ranta0/rest-and-go/internal/utils"
	httpErrors "github.com/ranta0/rest-and-go/internal/errors"
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

func (h *JWTController) Register(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Username string `json:"username" form:"username" validate:"required"`
		Password string `json:"password" form:"password" validate:"required"`
	}

	err := utils.HandlePayload(r, &request)
	if err != nil {
		utils.ErrorJsonResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = utils.ValidatePayload(&request)
	if err != nil {
		utils.ErrorJsonResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = h.userService.Register(request.Username, request.Password)
	if err != nil {
		utils.ErrorJsonResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessJsonResponse(w, r, http.StatusCreated, nil)
}

func (c *JWTController) Login(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Username string `json:"username" form:"username" validate:"required"`
		Password string `json:"password" form:"password" validate:"required"`
	}

	err := utils.HandlePayload(r, &request)
	if err != nil {
		utils.ErrorJsonResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = utils.ValidatePayload(&request)
	if err != nil {
		utils.ErrorJsonResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.userService.Login(request.Username, request.Password)
	if err != nil {
		utils.ErrorJsonResponse(w, r, http.StatusBadRequest, httpErrors.ErrCredentials.Error())
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
		utils.ErrorJsonResponse(w, r, http.StatusInternalServerError, httpErrors.ErrTokenCreate.Error())
		return
	}

	claims["token_type"] = c.tokenService.strings["refresh"]
	claims["exp"] = expirationTimeRefresh.Unix()
	refreshToken, err := c.tokenService.GenerateToken(claims)
	if err != nil {
		utils.ErrorJsonResponse(w, r, http.StatusInternalServerError, httpErrors.ErrTokenCreate.Error())
		return
	}

	utils.JWTJsonResponse(w, r, http.StatusOK, &utils.JWTResponse{
		Status:         "success",
		Message:        "user athenticated successfully",
		Type:           "bearer",
		AccessToken:    accessToken,
		ExpirationTime: expirationTime.Format(time.RFC3339),
		RefreshToken:   refreshToken,
	})
}

func (c *JWTController) Refresh(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RefreshToken string `json:"refresh_token" form:"refresh_token" validate:"required"`
	}

	err := utils.HandlePayload(r, &request)
	if err != nil {
		utils.ErrorJsonResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = utils.ValidatePayload(&request)
	if err != nil {
		utils.ErrorJsonResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	token, err := c.tokenService.ValidateToken(request.RefreshToken)
	if err != nil || !token.Valid {
		utils.ErrorJsonResponse(w, r, http.StatusUnauthorized, httpErrors.ErrTokenInvalid.Error())
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		utils.ErrorJsonResponse(w, r, http.StatusUnauthorized, httpErrors.ErrTokenInvalid.Error())
		return
	}

	// Check the token type
	if claims["token_type"] != c.tokenService.strings["refresh"] {
		utils.ErrorJsonResponse(w, r, http.StatusUnauthorized, httpErrors.ErrTokenType.Error())
		return
	}

	if c.tokenService.IsTokenRevoked(request.RefreshToken) {
		utils.ErrorJsonResponse(w, r, http.StatusUnauthorized, httpErrors.ErrTokenRevoked.Error())
		return
	}

	// Revoke previous token
	err = c.tokenService.RevokeToken(request.RefreshToken)
	if err != nil {
		utils.ErrorJsonResponse(w, r, http.StatusUnauthorized, httpErrors.ErrTokenRevokedFailure.Error())
		return
	}

	expirationTime := time.Now().Add(c.tokenService.expirationTime)
	expirationTimeRefresh := time.Now().Add(time.Hour * c.tokenService.expirationTimeRefresh)


	// Generate a new access token
	claims["token_type"] = c.tokenService.strings["access"]
	claims["exp"] = expirationTime.Unix()
	newAccessToken, err := c.tokenService.GenerateToken(claims)
	if err != nil {
		utils.ErrorJsonResponse(w, r, http.StatusInternalServerError, httpErrors.ErrTokenCreate.Error())
		return
	}

	// Generate a new refresh token
	claims["token_type"] = c.tokenService.strings["refresh"]
	claims["exp"] = expirationTimeRefresh.Unix()
	newRefreshToken, err := c.tokenService.GenerateToken(claims)
	if err != nil {
		utils.ErrorJsonResponse(w, r, http.StatusInternalServerError, httpErrors.ErrTokenCreate.Error())
		return
	}

	utils.JWTJsonResponse(w, r, http.StatusOK, &utils.JWTResponse{
		Status:         "success",
		Message:        "token refreshed successfully",
		Type:           "bearer",
		AccessToken:    newAccessToken,
		ExpirationTime: expirationTime.Format(time.RFC3339),
		RefreshToken:   newRefreshToken,
	})
}

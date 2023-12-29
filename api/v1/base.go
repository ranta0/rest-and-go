package api

import (
	"github.com/ranta0/rest-and-go/api/v1/route"
	"github.com/ranta0/rest-and-go/app"
	"github.com/ranta0/rest-and-go/domain/auth"
	"github.com/ranta0/rest-and-go/domain/user"
	"github.com/ranta0/rest-and-go/middleware"
)

func InitAPI(app *app.App) {
	userRepo := user.NewUserRepository(app.DB)
	tokenRepo := auth.NewRevokedJWTTokenRepository(app.DB)

	userService := user.NewUserService(userRepo)
	tokenService := auth.NewRevokedJWTTokenService(tokenRepo, app.Config)

	userController := user.NewUserController(userService)
	jwtController := auth.NewJWTController(userService, tokenService)

	// Api specs
	apiEndpoint := "/api/v1"

	// Middleware definition
	logger := middleware.NewLoggerMiddleware(app.LoggerAPI.GetLogger())
	app.Router.Use(logger.Apply)
	jwtMiddleware := middleware.NewJWTMiddleware(tokenService)

	// Routes
	route.BindAuth(app.Router, apiEndpoint, jwtController, jwtMiddleware)
	route.BindUser(app.Router, apiEndpoint, userController, jwtMiddleware)
}
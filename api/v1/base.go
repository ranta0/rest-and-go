package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ranta0/rest-and-go/api/v1/route"
	"github.com/ranta0/rest-and-go/app"
	"github.com/ranta0/rest-and-go/domain/auth"
	"github.com/ranta0/rest-and-go/domain/role"
	"github.com/ranta0/rest-and-go/domain/user"
	"github.com/ranta0/rest-and-go/logging"
	"github.com/ranta0/rest-and-go/middleware"
)

func InitAPI(app *app.App) {
	userRepo := user.NewUserRepository(app.DB)
	tokenRepo := auth.NewRevokedJWTTokenRepository(app.DB)
	roleRepo := role.NewRoleRepository(app.DB)

	userService := user.NewUserService(userRepo)
	tokenService := auth.NewRevokedJWTTokenService(tokenRepo, app.Config)
	roleService := role.NewRoleService(roleRepo)

	userController := user.NewUserController(userService, roleService)
	jwtController := auth.NewJWTController(userService, tokenService)
	roleController := role.NewRoleController(roleService)

	// Api specs
	apiEndpoint := "/api/v1"

	// Middleware definition
	logger := middleware.NewLoggerMiddleware(app.LoggerAPI.GetLogger())
	app.Router.Use(logger.Apply)
	jwtMiddleware := middleware.NewJWTMiddleware(tokenService)

	// Routes
	route.BindAuth(app.Router, apiEndpoint, jwtController, jwtMiddleware)
	route.BindUser(app.Router, apiEndpoint, userController, jwtMiddleware)
	route.BindRole(app.Router, apiEndpoint, roleController, jwtMiddleware)

	// Route List
	printRoutes(app.Router, app.Logger)
}

func printRoutes(r chi.Router, logger *logging.LoggerFile) {
	var routeInfo []string
	walkFunc := func(method string, route string, _ http.Handler, _ ...func(http.Handler) http.Handler) error {
		routeInfo = append(routeInfo, fmt.Sprintf("%s %s", method, route))
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		logger.Printf("Error collecting routes: %s\n", err)
	}

	uniformOutputRoutes := "Registered routes:\n"
	for _, info := range routeInfo {
		uniformOutputRoutes += info + "\n"
	}
	logger.Printf("%s", uniformOutputRoutes)
}

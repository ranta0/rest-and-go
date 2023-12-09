package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/ranta0/rest-and-go/internal/api"
	"github.com/ranta0/rest-and-go/internal/domain/user"
	"github.com/ranta0/rest-and-go/internal/middlewares"
)

func BindUser(router *chi.Mux, apiVersion string, controller *user.UserController, middleware middlewares.MiddlewareInterface) {
	userAPI := api.NewResource("users", controller, middleware)
	userAPI.BindRoutes(router, apiVersion);
}

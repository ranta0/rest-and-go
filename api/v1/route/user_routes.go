package route

import (
	"github.com/go-chi/chi/v5"

	"github.com/ranta0/rest-and-go/api"
	"github.com/ranta0/rest-and-go/domain/user"
	"github.com/ranta0/rest-and-go/middleware"
)

func BindUser(router *chi.Mux, apiVersion string, controller *user.UserController, middleware middleware.MiddlewareInterface) {
	userAPI := api.NewResource("/users", controller, middleware)
	userAPI.BindRoutes(router, apiVersion)
}

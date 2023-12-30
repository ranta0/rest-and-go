package route

import (
	"github.com/go-chi/chi/v5"

	"github.com/ranta0/rest-and-go/api"
	"github.com/ranta0/rest-and-go/domain/role"
	"github.com/ranta0/rest-and-go/middleware"
)

func BindRole(router *chi.Mux, apiVersion string, controller *role.RoleController, middleware middleware.MiddlewareInterface) {
	userAPI := api.NewResource("/roles", controller, middleware)
	userAPI.BindRoutes(router, apiVersion, nil)
}

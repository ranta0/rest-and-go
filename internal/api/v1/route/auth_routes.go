package route

import (
	"github.com/go-chi/chi/v5"

	"github.com/ranta0/rest-and-go/internal/domain/auth"
	"github.com/ranta0/rest-and-go/internal/middleware"
)

func BindAuth(router *chi.Mux, apiVersion string, controller *auth.JWTController, middleware middleware.MiddlewareInterface) {
	router.Post(apiVersion+"/register", controller.Register)
	router.Post(apiVersion+"/login", controller.Login)
	router.Post(apiVersion+"/refresh", controller.Refresh)

	router.Route(apiVersion, func(r chi.Router) {
		r.Use(middleware.Apply)

		r.Post("/logout", controller.Logout)
	})
}

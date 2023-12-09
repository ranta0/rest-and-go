package routes

import (
	"github.com/go-chi/chi/v5"

	"github.com/ranta0/rest-and-go/internal/domain/auth"
)

func BindAuth(router *chi.Mux, controller *auth.JWTController) {
	router.Post("/register", controller.Register)
	router.Post("/login", controller.Login)
	router.Post("/refresh", controller.Refresh)
}

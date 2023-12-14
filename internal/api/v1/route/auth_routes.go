package route

import (
	"github.com/go-chi/chi/v5"

	"github.com/ranta0/rest-and-go/internal/domain/auth"
)

func BindAuth(router *chi.Mux, apiVersion string, controller *auth.JWTController) {
	router.Post(apiVersion+"register", controller.Register)
	router.Post(apiVersion+"login", controller.Login)
	router.Post(apiVersion+"refresh", controller.Refresh)
}

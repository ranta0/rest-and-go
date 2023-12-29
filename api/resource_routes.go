package api

import (
	"github.com/go-chi/chi/v5"

	"github.com/ranta0/rest-and-go/domain"
	"github.com/ranta0/rest-and-go/middleware"
)

type Resource struct {
	name       string
	controller domain.ControllerInterface
	middleware middleware.MiddlewareInterface
	// ExtraRoutes []*chi.Router
}

func NewResource(name string, controller domain.ControllerInterface, middleware middleware.MiddlewareInterface) *Resource {
	return &Resource{
		name:       name,
		controller: controller,
		middleware: middleware,
		// ExtraRoutes: routes,
	}
}

// This can be used to simply the api building process
// - apiVersion must be in the form of "/api/v<number>/"
func (re *Resource) BindRoutes(router *chi.Mux, apiVersion string) {
	router.Route(apiVersion+re.name, func(r chi.Router) {
		r.Use(re.middleware.Apply)

		r.Get("/", re.controller.GetAll)
		r.Post("/", re.controller.Create)
		r.Get("/{id}", re.controller.GetByID)
		r.Patch("/{id}", re.controller.Update)
		r.Delete("/{id}", re.controller.Delete)
		// TODO: ExtraRoutes
	})
}

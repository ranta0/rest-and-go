package	domain

import (
	"net/http"
)

type ControllerInterface interface {
	// GET /api/v1/<resource> endpoint
	GetAll(w http.ResponseWriter, r *http.Request)

	// GET /api/v1/<resource>/{id} endpoint
	GetByID(w http.ResponseWriter, r *http.Request)

	// POST /api/v1/<resource> endpoint
	Create(w http.ResponseWriter, r *http.Request)

	// PATCH /api/v1/<resource>/{id} endpoint
	Update(w http.ResponseWriter, r *http.Request)

	// DELETE /api/v1/<resource>/{id} endpoint
	Delete(w http.ResponseWriter, r *http.Request)
}

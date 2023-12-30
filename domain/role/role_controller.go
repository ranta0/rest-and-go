package role

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ranta0/rest-and-go/form"
	"github.com/ranta0/rest-and-go/request"
	"github.com/ranta0/rest-and-go/response"
)

type RoleController struct {
	roleService *RoleService
}

func NewRoleController(roleService *RoleService) *RoleController {
	return &RoleController{
		roleService: roleService,
	}
}

// GET /api/v1/roles endpoint
func (c *RoleController) GetAll(w http.ResponseWriter, r *http.Request) {
	pagination := &form.Pagination{}

	err := request.Handler(r, pagination)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = request.Validator(pagination)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	count, err := c.roleService.CountAll()
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	paginator := form.NewPaginatorFromRequest(pagination, count)

	roles, err := c.roleService.GetAll(paginator.Limit(), paginator.Offset())
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	stub := response.NewOK(roles, paginator)
	response.OK(w, r, http.StatusOK, stub)
}

// GET /api/v1/roles/{id} endpoint
func (c *RoleController) GetByID(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "id")

	role, err := c.roleService.GetByID(roleID)
	if err != nil {
		response.Error(w, r, http.StatusNotFound, err.Error())
		return
	}

	stub := response.NewOK(role, nil)
	response.OK(w, r, http.StatusOK, stub)
}

// POST /api/v1/roles endpoint
func (c *RoleController) Create(w http.ResponseWriter, r *http.Request) {
	newRole := &Role{}

	err := request.Handler(r, newRole)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = request.Validator(newRole)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.roleService.Create(newRole); err != nil {
		response.Error(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	stub := response.NewOK(newRole, nil)
	response.OK(w, r, http.StatusCreated, stub)
}

// PATCH /api/v1/roles/{id} endpoint
func (c *RoleController) Update(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "id")
	updatedRole := &Role{}

	err := request.Handler(r, updatedRole)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = request.Validator(updatedRole)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.roleService.Update(roleID, updatedRole); err != nil {
		response.Error(w, r, http.StatusNotFound, err.Error())
		return
	}

	stub := response.NewOK(updatedRole, nil)
	response.OK(w, r, http.StatusOK, stub)
}

// DELETE /api/v1/roles/{id} endpoint
func (c *RoleController) Delete(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "id")

	if err := c.roleService.Delete(roleID); err != nil {
		response.Error(w, r, http.StatusNotFound, err.Error())
		return
	}

	stub := response.NewOK(nil, nil)
	response.OK(w, r, http.StatusOK, stub)
}

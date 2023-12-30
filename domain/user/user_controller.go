package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ranta0/rest-and-go/domain/role"
	"github.com/ranta0/rest-and-go/form"
	"github.com/ranta0/rest-and-go/request"
	"github.com/ranta0/rest-and-go/response"
)

type UserController struct {
	userService *UserService
	roleService *role.RoleService
}

func NewUserController(userService *UserService, roleService *role.RoleService) *UserController {
	return &UserController{
		userService: userService,
		roleService: roleService,
	}
}

// GET /api/v1/users endpoint
func (c *UserController) GetAll(w http.ResponseWriter, r *http.Request) {
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

	count, err := c.userService.CountAll()
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	paginator := form.NewPaginatorFromRequest(pagination, count)

	users, err := c.userService.GetAll(paginator.Limit(), paginator.Offset())
	if err != nil {
		response.Error(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	stub := response.NewOK(users, paginator)
	response.OK(w, r, http.StatusOK, stub)
}

// GET /api/v1/users/{id} endpoint
func (c *UserController) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	user, err := c.userService.GetByID(userID)
	if err != nil {
		response.Error(w, r, http.StatusNotFound, err.Error())
		return
	}

	stub := response.NewOK(user, nil)
	response.OK(w, r, http.StatusOK, stub)
}

// POST /api/v1/users endpoint
func (c *UserController) Create(w http.ResponseWriter, r *http.Request) {
	newUser := &User{}

	err := request.Handler(r, newUser)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = request.Validator(newUser)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.userService.Create(newUser); err != nil {
		response.Error(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	stub := response.NewOK(newUser, nil)
	response.OK(w, r, http.StatusCreated, stub)
}

// PATCH /api/v1/users/{id}/roles/{role_id} endpoint
func (c *UserController) SetRole(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	roleID := chi.URLParam(r, "role_id")

	user, err := c.userService.GetByID(userID)
	if err != nil {
		response.Error(w, r, http.StatusNotFound, err.Error())
		return
	}

	role, err := c.roleService.GetByID(roleID)
	if err != nil {
		response.Error(w, r, http.StatusNotFound, err.Error())
		return
	}

	if err := c.userService.userRepo.SetRole(user, role); err != nil {
		response.Error(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	user, err = c.userService.GetByID(userID)
	if err != nil {
		response.Error(w, r, http.StatusNotFound, err.Error())
		return
	}

	stub := response.NewOK(user, nil)
	response.OK(w, r, http.StatusOK, stub)
}

// DELETE /api/v1/users/{id}/roles/{role_id} endpoint
func (c *UserController) DeleteRole(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	roleID := chi.URLParam(r, "role_id")

	user, err := c.userService.GetByID(userID)
	if err != nil {
		response.Error(w, r, http.StatusNotFound, err.Error())
		return
	}

	role, err := c.roleService.GetByID(roleID)
	if err != nil {
		response.Error(w, r, http.StatusNotFound, err.Error())
		return
	}

	if err := c.userService.userRepo.DeleteRole(user, role); err != nil {
		response.Error(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	user, err = c.userService.GetByID(userID)
	if err != nil {
		response.Error(w, r, http.StatusNotFound, err.Error())
		return
	}

	stub := response.NewOK(user, nil)
	response.OK(w, r, http.StatusOK, stub)
}

// PATCH /api/v1/users/{id} endpoint
func (c *UserController) Update(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	updatedUser := &User{}

	err := request.Handler(r, updatedUser)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = request.Validator(updatedUser)
	if err != nil {
		response.Error(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.userService.Update(userID, updatedUser); err != nil {
		response.Error(w, r, http.StatusNotFound, err.Error())
		return
	}

	stub := response.NewOK(updatedUser, nil)
	response.OK(w, r, http.StatusOK, stub)
}

// DELETE /api/v1/users/{id} endpoint
func (c *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	if err := c.userService.Delete(userID); err != nil {
		response.Error(w, r, http.StatusNotFound, err.Error())
		return
	}

	stub := response.NewOK(nil, nil)
	response.OK(w, r, http.StatusOK, stub)
}

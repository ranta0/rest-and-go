package user

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ranta0/rest-and-go/internal/utils"
)

type UserController struct {
	userService *UserService
}

func NewUserController(userService *UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// GET /api/v1/users endpoint
func (c *UserController) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := c.userService.GetAll()
	if err != nil {
		utils.ErrorJsonResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessJsonResponse(w, r, http.StatusOK, users)
}

// GET /api/v1/users/{id} endpoint
func (c *UserController) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	user, err := c.userService.GetByID(userID)
	if err != nil {
		utils.ErrorJsonResponse(w, r, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessJsonResponse(w, r, http.StatusOK, user)
}

// POST /api/v1/users endpoint
func (c *UserController) Create(w http.ResponseWriter, r *http.Request) {
	var newUser User

	err := utils.HandlePayload(r, &newUser)
	if err != nil {
		utils.ErrorJsonResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = utils.ValidatePayload(&newUser)
	if err != nil {
		utils.ErrorJsonResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.userService.Create(&newUser); err != nil {
		utils.ErrorJsonResponse(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessJsonResponse(w, r, http.StatusCreated, newUser)
}

// PATCH /api/v1/users/{id} endpoint
func (c *UserController) Update(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	var updatedUser User

	err := utils.HandlePayload(r, &updatedUser)
	if err != nil {
		utils.ErrorJsonResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	err = utils.ValidatePayload(&updatedUser)
	if err != nil {
		utils.ErrorJsonResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.userService.Update(userID, &updatedUser); err != nil {
		utils.ErrorJsonResponse(w, r, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessJsonResponse(w, r, http.StatusOK, updatedUser)
}

// DELETE /api/v1/users/{id} endpoint
func (c *UserController) Delete(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")

	if err := c.userService.Delete(userID); err != nil {
		utils.ErrorJsonResponse(w, r, http.StatusNotFound, err.Error())
		return
	}

	utils.SuccessJsonResponse(w, r, http.StatusOK, nil)
}

package handler

import (
	"go-playground/internal/domain"
	"go-playground/internal/service"
	"go-playground/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// @title go-playground API
// @version 1.0
// @description User management API with PostgreSQL and Redis
// @host localhost:8080
// @BasePath /api

// @Create godoc
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body domain.CreateUserRequest true "User details"
// @Success 201 {object} domain.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	user, err := h.userService.Create(c.Request.Context(), &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

// @GetByID godoc
// @Summary Get user by ID
// @Description Get user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} domain.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid user ID",
		})
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetAll godoc
// @Summary      Get all users
// @Description  Retrieve all users from the system
// @Tags         users
// @Produce      json
// @Security BearerAuth
// @Security UserIdAuth
// @Success      200  {array}   domain.User
// @Failure      401  {object}  map[string]string
// @Router       /users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.userService.GetAll(c.Request.Context())
	if err != nil {
		util.HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

// Update godoc
// @Summary      Update user
// @Description  Update user information by their UUID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string                 true  "User ID"
// @Param        user  body      domain.UpdateUserRequest  true  "User information"
// @Success      200   {object}  domain.User
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid user ID",
		})
		return
	}

	var req domain.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	user, err := h.userService.Update(c.Request.Context(), id, &req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// Delete godoc
// @Summary      Delete user
// @Description  Delete a user by their UUID
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid user ID",
		})
		return
	}

	if err := h.userService.Delete(c.Request.Context(), id); err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// @Summary Get current user
// @Description Get current user's data
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Security UserIdAuth
// @Success 200 {object} domain.User
// @Failure 401 {object} map[string]string
// @Router /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		util.HandleError(c, domain.AuthenticationError{Message: "unauthorized"})
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// Implement other handler methods (Create, Update, Delete)

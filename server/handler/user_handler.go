package handler

import (
	"go-playground/pkg/logging"
	"go-playground/server/domain"
	"go-playground/server/service"
	"go-playground/server/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	userService *service.UserService
	logger      zerolog.Logger
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logging.GetLogger(),
	}
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
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming create user request")

	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to bind create user request")
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	user, err := h.userService.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error().
			Err(err).
			Interface("request", req).
			Msg("Failed to create user")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("user_id", user.ID).
		Str("email", user.Email).
		Str("name", user.Name).
		Str("status", string(user.Status)).
		Msg("User created successfully")

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
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming get user request")

	id := c.Param("id")
	if id == "" {
		h.logger.Error().
			Msg("Missing user ID")
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid user ID",
		})
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("user_id", id).
			Msg("Failed to get user")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("user_id", user.ID).
		Str("email", user.Email).
		Str("name", user.Name).
		Str("status", string(user.Status)).
		Msg("User retrieved successfully")

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
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming get all users request")

	users, err := h.userService.GetAll(c.Request.Context())
	if err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to get users")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Int("users_count", len(users)).
		Msg("Users retrieved successfully")

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
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming update user request")

	id := c.Param("id")
	if id == "" {
		h.logger.Error().
			Msg("Missing user ID")
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid user ID",
		})
		return
	}

	var req domain.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error().
			Err(err).
			Msg("Failed to bind update user request")
		util.HandleError(c, domain.ValidationError{Message: err.Error()})
		return
	}

	user, err := h.userService.Update(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.Error().
			Err(err).
			Str("user_id", id).
			Interface("request", req).
			Msg("Failed to update user")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("user_id", user.ID).
		Str("email", user.Email).
		Str("name", user.Name).
		Str("status", string(user.Status)).
		Msg("User updated successfully")

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
	h.logger.Info().
		Str("method", c.Request.Method).
		Str("url", c.Request.URL.RequestURI()).
		Str("user_agent", c.Request.UserAgent()).
		Dur("elapsed_ms", time.Since(time.Now())).
		Msg("incoming delete user request")

	id := c.Param("id")
	if id == "" {
		h.logger.Error().
			Msg("Missing user ID")
		util.HandleError(c, domain.ValidationError{
			Field:   "id",
			Message: "invalid user ID",
		})
		return
	}

	if err := h.userService.Delete(c.Request.Context(), id); err != nil {
		h.logger.Error().
			Err(err).
			Str("user_id", id).
			Msg("Failed to delete user")
		util.HandleError(c, err)
		return
	}

	h.logger.Info().
		Str("user_id", id).
		Msg("User deleted successfully")

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

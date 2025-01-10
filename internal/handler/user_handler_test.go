package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-cursor/internal/domain"
	"go-cursor/internal/handler"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of the user service
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(req *domain.CreateUserRequest) (*domain.User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) GetByID(id string) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) GetAll() ([]domain.User, error) {
	args := m.Called()
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserService) Update(id string, req *domain.UpdateUserRequest) (*domain.User, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupTest() (*gin.Engine, *MockUserService) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockUserService)
	userHandler := handler.NewUserHandler(mockService)

	router := gin.Default()
	api := router.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("/", userHandler.Create)
			users.GET("/", userHandler.GetAll)
			users.GET("/:id", userHandler.GetByID)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
		}
	}

	return router, mockService
}

func TestCreateUser(t *testing.T) {
	router, mockService := setupTest()

	t.Run("Success", func(t *testing.T) {
		user := &domain.User{
			ID:    "123e4567-e89b-12d3-a456-426614174000",
			Email: "test@example.com",
			Name:  "Test User",
			Phone: "1234567890",
		}

		mockService.On("Create", mock.AnythingOfType("*domain.CreateUserRequest")).Return(user, nil).Once()

		body := []byte(`{
			"email": "test@example.com",
			"password": "password123",
			"name": "Test User",
			"phone": "1234567890"
		}`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response domain.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, response.ID)
		assert.Equal(t, user.Email, response.Email)
	})

	t.Run("Invalid Request", func(t *testing.T) {
		body := []byte(`{
			"email": "invalid-email",
			"password": "123",
			"name": "",
			"phone": ""
		}`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetUser(t *testing.T) {
	router, mockService := setupTest()

	t.Run("Success", func(t *testing.T) {
		user := &domain.User{
			ID:    "123e4567-e89b-12d3-a456-426614174000",
			Email: "test@example.com",
			Name:  "Test User",
			Phone: "1234567890",
		}

		mockService.On("GetByID", user.ID).Return(user, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/users/"+user.ID, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, response.ID)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockService.On("GetByID", "non-existent-id").Return(nil, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/users/non-existent-id", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGetAllUsers(t *testing.T) {
	router, mockService := setupTest()

	users := []domain.User{
		{
			ID:    "123e4567-e89b-12d3-a456-426614174000",
			Email: "test1@example.com",
			Name:  "Test User 1",
			Phone: "1234567890",
		},
		{
			ID:    "223e4567-e89b-12d3-a456-426614174000",
			Email: "test2@example.com",
			Name:  "Test User 2",
			Phone: "0987654321",
		},
	}

	mockService.On("GetAll").Return(users, nil).Once()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/users", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []domain.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
}

func TestUpdateUser(t *testing.T) {
	router, mockService := setupTest()

	t.Run("Success", func(t *testing.T) {
		userID := "123e4567-e89b-12d3-a456-426614174000"
		updatedUser := &domain.User{
			ID:    userID,
			Email: "test@example.com",
			Name:  "Updated Name",
			Phone: "0987654321",
		}

		mockService.On("Update", userID, mock.AnythingOfType("*domain.UpdateUserRequest")).Return(updatedUser, nil).Once()

		body := []byte(`{
			"name": "Updated Name",
			"phone": "0987654321"
		}`)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/users/"+userID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response domain.User
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, updatedUser.Name, response.Name)
		assert.Equal(t, updatedUser.Phone, response.Phone)
	})
}

func TestDeleteUser(t *testing.T) {
	router, mockService := setupTest()

	t.Run("Success", func(t *testing.T) {
		userID := "123e4567-e89b-12d3-a456-426614174000"
		mockService.On("Delete", userID).Return(nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/users/"+userID, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Not Found", func(t *testing.T) {
		userID := "non-existent-id"
		mockService.On("Delete", userID).Return(errors.New("user not found")).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/api/users/"+userID, nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

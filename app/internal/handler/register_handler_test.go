package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/models"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/service"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/test_helper/funcs"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/test_helper/mocks/svc_mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegisterSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := map[string]string{
		"name":     "Test User",
		"email":    "testuser@example.com",
		"password": "securepassword",
	}
	jsonBody, _ := json.Marshal(body)

	reqBody := strings.NewReader(string(jsonBody))
	req := httptest.NewRequest("POST", "/", reqBody)
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	input := service.RegisterUserInput{
		Name:     "Test User",
		Email:    "testuser@example.com",
		Password: "securepassword",
	}

	user := models.User{
		UUID:         "some-uuid",
		Username:     "Test User",
		Email:        "testuser@example.com",
		PasswordHash: "securepassword",
	}

	registerUserMock := new(svc_mock.UserRegisterSvcStructMock)
	registerUserMock.On(
		"RegisterUser",
		input,
	).Return(user, nil)

	handler := NewRegisterHandler(registerUserMock)
	handler.Register(c)

	assert.Equal(t, http.StatusOK, w.Code)

	result := map[string]string{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, "some-uuid", result["uuid"])
	assert.Equal(t, "Test User", result["username"])
	assert.Equal(t, "testuser@example.com", result["email"])
}

func TestRegisterFailRegisterUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := map[string]string{
		"name":     "Test User",
		"email":    "testuser@example.com",
		"password": "securepassword",
	}
	jsonBody, _ := json.Marshal(body)

	reqBody := strings.NewReader(string(jsonBody))
	req := httptest.NewRequest("POST", "/", reqBody)
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	input := service.RegisterUserInput{
		Name:     "Test User",
		Email:    "testuser@example.com",
		Password: "securepassword",
	}

	registerUserMock := new(svc_mock.UserRegisterSvcStructMock)
	registerUserMock.On(
		"RegisterUser",
		input,
	).Return(models.User{}, assert.AnError)

	handler := NewRegisterHandler(registerUserMock)
	handler.Register(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestRegisterFailedValidation(t *testing.T) {
	expected := []*funcs.ValidationSetting{
		{
			Title:     "Name is required",
			Key:       "name",
			ErrorType: "required",
		},
		{
			Title:     "Email is required",
			Key:       "email",
			ErrorType: "required",
		},
		{
			Title:     "Email is invalid",
			Key:       "email",
			ErrorType: "email",
		},
		{
			Title:     "Password is required",
			Key:       "password",
			ErrorType: "required",
		},
		{
			Title:     "Password is too short",
			Key:       "password",
			ErrorType: "min",
		},
	}

	for _, exp := range expected {
		t.Run(exp.Title, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			name := funcs.CreateValidationTestDataName(exp)

			email := funcs.CreateValidationTestDataEmail(exp)

			password := funcs.CreateValidationTestDataPassword(exp)

			body := map[string]string{
				"name":     name,
				"email":    email,
				"password": password,
			}
			jsonBody, _ := json.Marshal(body)

			reqBody := strings.NewReader(string(jsonBody))
			req := httptest.NewRequest("POST", "/", reqBody)
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			registerUserMock := new(svc_mock.UserRegisterSvcStructMock)

			handler := NewRegisterHandler(registerUserMock)
			handler.Register(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

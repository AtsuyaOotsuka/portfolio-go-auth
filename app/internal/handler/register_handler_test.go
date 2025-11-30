package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/models"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/service"
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
	expected := map[string]any{
		"required_name": map[string]string{
			"key":   "name",
			"error": "required",
		},
		"required_email": map[string]string{
			"key":   "email",
			"error": "required",
		},
		"invalid_email": map[string]string{
			"key":   "email",
			"error": "email",
		},
		"required_password": map[string]string{
			"key":   "password",
			"error": "required",
		},
		"short_password": map[string]string{
			"key":   "password",
			"error": "min",
		},
	}

	for name, exp := range expected {
		t.Run(name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			name := func() string {
				if val, ok := exp.(map[string]string); ok && val["key"] == "name" {
					if val["error"] == "required" {
						return ""
					}
				}
				return "Test User"
			}()

			email := func() string {
				if val, ok := exp.(map[string]string); ok && val["key"] == "email" {
					if val["error"] == "required" {
						return ""
					}
					if val["error"] == "email" {
						return "invalid-email"
					}
				}
				return "testuser@example.com"
			}()

			password := func() string {
				if val, ok := exp.(map[string]string); ok && val["key"] == "password" {
					if val["error"] == "required" {
						return ""
					}
					if val["error"] == "min" {
						return "short"
					}
				}
				return "securepassword"
			}()

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

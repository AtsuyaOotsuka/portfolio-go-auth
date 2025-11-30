package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/service"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/test_helper/funcs"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/test_helper/mocks/svc_mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestLoginSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := map[string]string{
		"email":    "user@example.com",
		"password": "securepassword",
	}
	jsonBody, _ := json.Marshal(body)
	reqBody := strings.NewReader(string(jsonBody))
	req := httptest.NewRequest("POST", "/", reqBody)
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	input := service.LoginInput{
		Email:    "user@example.com",
		Password: "securepassword",
	}

	response := &service.AuthOutput{
		AccessToken:  "access_token_value",
		RefreshToken: "refresh_token_value",
	}

	authSvcMock := new(svc_mock.AuthSvcMock)
	authSvcMock.On("Login", input).Return(response, nil)

	handler := NewAuthHandler(authSvcMock)
	handler.Login(c)

	assert.Equal(t, http.StatusOK, w.Code)

	result := map[string]interface{}{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, "access_token_value", result["access_token"])
	assert.Equal(t, "refresh_token_value", result["refresh_token"])
	assert.Equal(t, "Bearer", result["token_type"])
	assert.Equal(t, float64(3600), result["expires_in"])
}

func TestLoginFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := map[string]string{
		"email":    "user@example.com",
		"password": "wrongpassword",
	}
	jsonBody, _ := json.Marshal(body)
	reqBody := strings.NewReader(string(jsonBody))
	req := httptest.NewRequest("POST", "/", reqBody)
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	input := service.LoginInput{
		Email:    "user@example.com",
		Password: "wrongpassword",
	}

	authSvcMock := new(svc_mock.AuthSvcMock)
	authSvcMock.On("Login", input).Return(&service.AuthOutput{}, fmt.Errorf("Invalid email or password"))

	handler := NewAuthHandler(authSvcMock)
	handler.Login(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	result := map[string]interface{}{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, "Invalid email or password", result["error"])
}

func TestLoginFailedValidation(t *testing.T) {
	expected := []*funcs.ValidationSetting{
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

			email := funcs.CreateValidationTestDataEmail(exp)
			password := funcs.CreateValidationTestDataPassword(exp)

			body := map[string]string{
				"email":    email,
				"password": password,
			}
			jsonBody, _ := json.Marshal(body)
			reqBody := strings.NewReader(string(jsonBody))
			req := httptest.NewRequest("POST", "/", reqBody)
			req.Header.Set("Content-Type", "application/json")
			c.Request = req

			authSvcMock := new(svc_mock.AuthSvcMock)
			handler := NewAuthHandler(authSvcMock)
			handler.Login(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestRefreshSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := map[string]string{
		"refresh_token": "valid_refresh_token",
	}
	jsonBody, _ := json.Marshal(body)
	reqBody := strings.NewReader(string(jsonBody))
	req := httptest.NewRequest("POST", "/", reqBody)
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	input := service.RefreshInput{
		RefreshToken: "valid_refresh_token",
		IpAddress:    c.ClientIP(),
	}

	response := &service.AuthOutput{
		AccessToken:  "new_access_token",
		RefreshToken: "new_refresh_token",
	}

	authSvcMock := new(svc_mock.AuthSvcMock)
	authSvcMock.On("Refresh", input).Return(response, nil)

	handler := NewAuthHandler(authSvcMock)
	handler.Refresh(c)

	assert.Equal(t, http.StatusOK, w.Code)

	result := map[string]interface{}{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, "new_access_token", result["access_token"])
	assert.Equal(t, "new_refresh_token", result["refresh_token"])
	assert.Equal(t, "Bearer", result["token_type"])
	assert.Equal(t, float64(3600), result["expires_in"])
}

func TestRefreshFail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := map[string]string{
		"refresh_token": "invalid_refresh_token",
	}
	jsonBody, _ := json.Marshal(body)
	reqBody := strings.NewReader(string(jsonBody))
	req := httptest.NewRequest("POST", "/", reqBody)
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	input := service.RefreshInput{
		RefreshToken: "invalid_refresh_token",
		IpAddress:    c.ClientIP(),
	}

	authSvcMock := new(svc_mock.AuthSvcMock)
	authSvcMock.On("Refresh", input).Return(&service.AuthOutput{}, fmt.Errorf("Invalid refresh token"))

	handler := NewAuthHandler(authSvcMock)
	handler.Refresh(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	result := map[string]interface{}{}
	err := json.Unmarshal(w.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, "Invalid refresh token", result["error"])
}

func TestRefreshFailedValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := map[string]string{
		"refresh_token": "",
	}
	jsonBody, _ := json.Marshal(body)
	reqBody := strings.NewReader(string(jsonBody))
	req := httptest.NewRequest("POST", "/", reqBody)
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	authSvcMock := new(svc_mock.AuthSvcMock)
	handler := NewAuthHandler(authSvcMock)
	handler.Refresh(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

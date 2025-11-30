package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/test_helper/mocks/svc_mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCsrfHandler(t *testing.T) {
	mockCsrfSvc := new(svc_mock.CsrfSvcMockStruct)
	mockCsrfSvc.On("Verify", "valid_token", mock.Anything, mock.AnythingOfType("int64")).Return(nil)

	r := gin.New()
	r.Use(NewCSRFMiddleware(mockCsrfSvc).Handler())
	r.POST("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "not set csrf token")
}

func TestCSRFMiddleware_InvalidToken(t *testing.T) {
	mockCsrfSvc := new(svc_mock.CsrfSvcMockStruct)
	mockCsrfSvc.On("Verify", "invalid-token", mock.Anything, mock.AnythingOfType("int64")).Return(fmt.Errorf("invalid"))

	r := gin.New()
	r.Use(NewCSRFMiddleware(mockCsrfSvc).Handler())
	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.Header.Set("X-CSRF-Token", "invalid-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "invalid csrf token")
}

func TestCSRFMiddlewareForGET(t *testing.T) {
	mockCsrfSvc := new(svc_mock.CsrfSvcMockStruct)
	mockCsrfSvc.On("Verify", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)

	r := gin.New()
	r.Use(NewCSRFMiddleware(mockCsrfSvc).Handler())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "GET success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "GET success")
}

func TestCSRFMiddlewareForCookie(t *testing.T) {
	mockCsrfSvc := new(svc_mock.CsrfSvcMockStruct)
	mockCsrfSvc.On("Verify", "cookie_token", mock.Anything, mock.AnythingOfType("int64")).Return(nil)

	r := gin.New()
	r.Use(NewCSRFMiddleware(mockCsrfSvc).Handler())
	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "POST success"})
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.AddCookie(&http.Cookie{Name: "csrf_token", Value: "cookie_token"})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "POST success")
}

func TestCSRFMiddlewareSuccess(t *testing.T) {
	mockCsrfSvc := new(svc_mock.CsrfSvcMockStruct)
	mockCsrfSvc.On("Verify", "valid_token", mock.Anything, mock.AnythingOfType("int64")).Return(nil)

	r := gin.New()
	r.Use(NewCSRFMiddleware(mockCsrfSvc).Handler())
	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "POST success"})
	})

	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	req.Header.Set("X-CSRF-Token", "valid_token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "POST success")
}

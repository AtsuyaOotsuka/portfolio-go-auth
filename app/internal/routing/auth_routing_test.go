package routing

import (
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/test_helper/funcs"
	"github.com/gin-gonic/gin"
)

type MockAuthHandler struct{}

func (m *MockAuthHandler) Login(c *gin.Context) {
	c.JSON(200, gin.H{"message": "logged in"})
}

func (m *MockAuthHandler) Refresh(c *gin.Context) {
	c.JSON(200, gin.H{"message": "token refreshed"})
}
func TestAuthRouting(t *testing.T) {
	expected := []funcs.ExpectedRoute{
		{
			Method: "POST",
			Path:   "/auth/login",
		},
		{
			Method: "POST",
			Path:   "/auth/refresh",
		},
	}

	g := gin.Default()
	r := NewRouting(g, nil)
	r.AuthRouting(&MockAuthHandler{})

	funcs.EachExepectedRoute(expected, g, t)
}

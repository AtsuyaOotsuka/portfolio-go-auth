package routing

import (
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/test_helper/funcs"
	"github.com/gin-gonic/gin"
)

type MockRgisterHandler struct{}

func (m *MockRgisterHandler) Register(c *gin.Context) {
	c.JSON(200, gin.H{"message": "registered"})
}

func TestRegisterRouting(t *testing.T) {
	expected := []funcs.ExpectedRoute{
		{Path: "/register", Method: "POST"},
	}

	g := gin.Default()
	r := NewRouting(g, nil)
	r.RegisterRouting(&MockRgisterHandler{})

	funcs.EachExepectedRoute(expected, g, t)
}

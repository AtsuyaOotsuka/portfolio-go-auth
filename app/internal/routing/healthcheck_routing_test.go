package routing

import (
	"net/http"
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/test_helper/funcs"
	"github.com/gin-gonic/gin"
)

type MockHealthCheckHandler struct{}

func (m *MockHealthCheckHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func TestHealthCheckRouting(t *testing.T) {
	expected := []funcs.ExpectedRoute{
		{Path: "/healthcheck", Method: "GET"},
	}

	g := gin.Default()
	r := NewRouting(g, nil)
	r.HealthCheckRoute(&MockHealthCheckHandler{})

	funcs.EachExepectedRoute(expected, g, t)
}

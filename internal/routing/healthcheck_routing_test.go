package routing

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockHealthCheckHandler struct{}

func (m *MockHealthCheckHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func TestHealthCheckRouting(t *testing.T) {
	expected := map[string]string{
		"/healthcheck": "GET",
	}

	g := gin.Default()
	r := NewRouting(g, nil)
	r.HealthCheckRoute(&MockHealthCheckHandler{})

	for path, method := range expected {
		found := false
		for _, route := range g.Routes() {
			if route.Path == path && route.Method == method {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected route %s [%s] to be registered", path, method)
	}
}

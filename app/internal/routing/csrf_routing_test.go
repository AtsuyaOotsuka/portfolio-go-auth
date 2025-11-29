package routing

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockCSRFHandler struct{}

func (m *MockCSRFHandler) CsrfGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"csrf_token": "mock_token"})
}

func TestCsrfRoute(t *testing.T) {
	expected := map[string]string{
		"/csrf/get": "GET",
	}

	g := gin.Default()
	r := NewRouting(g, nil)
	r.CsrfRoute(&MockCSRFHandler{})

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

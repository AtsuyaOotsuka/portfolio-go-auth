package routing

import (
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewRouting(t *testing.T) {
	g := &gin.Engine{}
	m := &middleware.Middleware{}
	r := NewRouting(g, m)

	assert.Equal(t, g, r.gin)
	assert.Equal(t, m, r.middleware)
}

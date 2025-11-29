package middleware

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewMiddleware(t *testing.T) {
	g := &gin.Engine{}
	m := NewMiddleware(g)

	assert.Equal(t, g, m.g)
}

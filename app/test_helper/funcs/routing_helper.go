package funcs

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type ExpectedRoute struct {
	Path   string
	Method string
}

func EachExepectedRoute(
	expected []ExpectedRoute,
	g *gin.Engine,
	t assert.TestingT,
) {
	for _, er := range expected {
		found := false
		for _, route := range g.Routes() {
			if route.Path == er.Path && route.Method == er.Method {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected route %s [%s] to be registered", er.Path, er.Method)
	}
}

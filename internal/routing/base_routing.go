package routing

import (
	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Routing struct {
	gin        *gin.Engine
	middleware *middleware.Middleware
}

func NewRouting(
	gin *gin.Engine,
	middleware *middleware.Middleware,
) *Routing {
	return &Routing{
		gin:        gin,
		middleware: middleware,
	}
}

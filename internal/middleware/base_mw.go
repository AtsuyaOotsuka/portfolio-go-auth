package middleware

import "github.com/gin-gonic/gin"

type Middleware struct {
	g *gin.Engine
}

func NewMiddleware(r *gin.Engine) *Middleware {
	return &Middleware{
		g: r,
	}
}

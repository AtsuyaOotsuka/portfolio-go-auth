package middleware

import (
	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/service"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabcsrf"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	g    *gin.Engine
	Csrf gin.HandlerFunc
}

func NewMiddleware(r *gin.Engine) *Middleware {

	csrf := NewCSRFMiddleware(
		service.NewCsrfSvcStruct(
			atylabcsrf.NewCsrfPkgStruct(),
		),
	)

	return &Middleware{
		g:    r,
		Csrf: csrf.Handler(),
	}
}

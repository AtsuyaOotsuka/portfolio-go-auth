package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/service"
	"github.com/gin-gonic/gin"
)

type CSRFMiddlewareInterface interface {
	Handler() gin.HandlerFunc
}

type CSRFMiddleware struct {
	csrf service.CsrfSvcInterface
}

func NewCSRFMiddleware(
	v service.CsrfSvcInterface,
) CSRFMiddlewareInterface {
	return &CSRFMiddleware{
		csrf: v,
	}
}

func (m *CSRFMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodHead {
			c.Next()
			return
		}
		token := c.GetHeader("X-CSRF-Token")
		if token == "" {
			token = c.PostForm("_token")
		}
		if token == "" {
			cookie, err := c.Cookie("csrf_token")
			if err == nil {
				token = cookie
			}
		}
		if token == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "not set csrf token"})
			return
		}
		if err := m.csrf.Verify(
			token,
			os.Getenv("CSRF_TOKEN"),
			time.Now().Unix(),
		); err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid csrf token"})
			return
		}
		c.Next()
	}
}

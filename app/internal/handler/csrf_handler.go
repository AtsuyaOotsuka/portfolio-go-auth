package handler

import (
	"os"
	"time"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/service"
	"github.com/gin-gonic/gin"
)

type CSRFHandlerInterface interface {
	CsrfGet(c *gin.Context)
}

type CSRFHandlerStruct struct {
	BaseHandler
	service service.CsrfSvcInterface
}

func NewCSRFHandler(
	service service.CsrfSvcInterface,
) *CSRFHandlerStruct {
	return &CSRFHandlerStruct{
		service: service,
	}
}

func (h *CSRFHandlerStruct) CsrfGet(c *gin.Context) {
	secret := os.Getenv("CSRF_TOKEN")
	token := h.service.CreateCSRFToken(time.Now().Unix(), secret)
	c.SetCookie("csrf_token", token, 3600, "/", "", false, true)
	c.JSON(200, gin.H{
		"csrf_token": token,
	})
}

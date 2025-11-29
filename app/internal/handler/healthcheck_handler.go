package handler

import "github.com/gin-gonic/gin"

type HealthCheckHandlerInterface interface {
	Check(c *gin.Context)
}

type HealthCheckHandler struct {
	BaseHandler
}

func NewHealthCheckHandler() *HealthCheckHandler {
	return &HealthCheckHandler{}
}

func (h *HealthCheckHandler) Check(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

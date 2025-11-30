package handler

import (
	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/service"
	"github.com/gin-gonic/gin"
)

type RegisterHandlerInterface interface {
	Register(c *gin.Context)
}

type RegisterHandlerStruct struct {
	BaseHandler
	service service.UserRegisterSvcInterface
}

func NewRegisterHandler(
	service service.UserRegisterSvcInterface,
) *RegisterHandlerStruct {
	return &RegisterHandlerStruct{
		service: service,
	}
}

type registerRequest struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=8"`
}

func (h *RegisterHandlerStruct) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	input := service.RegisterUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := h.service.RegisterUser(input)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"uuid":     user.UUID,
		"username": user.Username,
		"email":    user.Email,
	})

}

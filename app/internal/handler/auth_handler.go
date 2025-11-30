package handler

import (
	"net/http"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/repositories"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandlerInterface interface {
	Login(c *gin.Context)
	Refresh(c *gin.Context)
}

type AuthHandlerStruct struct {
	BaseHandler
	service service.AuthSvcInterface
	repo    repositories.UserRepoInterface
}

func NewAuthHandler(
	service service.AuthSvcInterface,
	repo repositories.UserRepoInterface,
) *AuthHandlerStruct {
	return &AuthHandlerStruct{
		service: service,
		repo:    repo,
	}
}

type loginRequest struct {
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=8"`
}

func (h *AuthHandlerStruct) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Login(service.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	resp := map[string]interface{}{
		"access_token":  response.AccessToken,
		"refresh_token": response.RefreshToken,
		"token_type":    "Bearer",
		"expires_in":    3600,
	}

	c.JSON(http.StatusOK, resp)
}

type refreshRequest struct {
	RefreshToken string `form:"refresh_token" json:"refresh_token" binding:"required"`
}

func (h *AuthHandlerStruct) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Refresh(service.RefreshInput{
		RefreshToken: req.RefreshToken,
		IpAddress:    c.ClientIP(),
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	resp := map[string]interface{}{
		"access_token":  response.AccessToken,
		"refresh_token": response.RefreshToken,
		"token_type":    "Bearer",
		"expires_in":    3600,
	}

	c.JSON(http.StatusOK, resp)
}

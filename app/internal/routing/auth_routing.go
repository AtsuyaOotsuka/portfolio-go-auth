package routing

import "github.com/AtsuyaOotsuka/portfolio-go-auth/internal/handler"

func (r *Routing) AuthRouting(
	authHandler handler.AuthHandlerInterface,
) {
	authGroup := r.gin.Group("/auth")
	authGroup.POST("/login", authHandler.Login)
	authGroup.POST("/refresh", authHandler.Refresh)
}

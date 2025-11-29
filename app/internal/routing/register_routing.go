package routing

import "github.com/AtsuyaOotsuka/portfolio-go-auth/internal/handler"

func (r *Routing) RegisterRouting(
	registerHandler handler.RegisterHandlerInterface,
) {
	r.gin.POST("/register", registerHandler.Register)
}

package routing

import "github.com/AtsuyaOotsuka/portfolio-go-auth/internal/handler"

func (r *Routing) HealthCheckRoute(
	handler handler.HealthCheckHandlerInterface,
) {
	r.gin.GET("/healthcheck", handler.Check)
}

package routing

import "github.com/AtsuyaOotsuka/portfolio-go-auth/internal/handler"

func (r *Routing) CsrfRoute(
	csrfHandler handler.CSRFHandlerInterface,
) {
	routerGroup := r.gin.Group("/csrf")
	routerGroup.GET("/get", csrfHandler.CsrfGet)
}

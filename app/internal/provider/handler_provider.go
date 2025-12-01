package provider

import "github.com/AtsuyaOotsuka/portfolio-go-auth/internal/handler"

func (p *Provider) BindRegisterHandler() *handler.RegisterHandlerStruct {
	return handler.NewRegisterHandler(
		p.bindRegisterSvc(),
	)
}

func (p *Provider) BindAuthHandler() *handler.AuthHandlerStruct {
	return handler.NewAuthHandler(
		p.bindAuthSvc(),
	)
}

func (p *Provider) BindCSRFHandler() *handler.CSRFHandlerStruct {
	return handler.NewCSRFHandler(
		p.bindCsrfSvc(),
	)
}

func (p *Provider) BindHealthCheckHandler() *handler.HealthCheckHandler {
	return handler.NewHealthCheckHandler()
}

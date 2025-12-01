package app

import (
	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/middleware"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/provider"
)

func (a *App) initProviders() {
	a.provider = provider.NewProvider(a.db)
}

func (a *App) initMiddlewares() {
	// ミドルウェアの初期化
	a.middleware = middleware.NewMiddleware(a.gin)
}

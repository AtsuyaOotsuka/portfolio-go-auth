package app

import "github.com/AtsuyaOotsuka/portfolio-go-auth/internal/routing"

func (a *App) entryBeforeGlobalMiddleware() {
	// 前処理系ミドルウェアをここに追加
	// a.gin.Use(a.middleware.Firewall)
	a.gin.Use(a.middleware.Csrf)
}

func (a *App) entryAfterGlobalMiddleware() {
	// 後処理系ミドルウェアをここに追加
}

func (a *App) entryRoutes() {
	// ルーティングの初期化
	routing := routing.NewRouting(a.gin, a.middleware)

	routing.HealthCheckRoute(
		a.provider.BindHealthCheckHandler(),
	)
	routing.CsrfRoute(
		a.provider.BindCSRFHandler(),
	)
	routing.RegisterRouting(
		a.provider.BindRegisterHandler(),
	)
	routing.AuthRouting(
		a.provider.BindAuthHandler(),
	)
}

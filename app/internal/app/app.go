package app

import (
	"database/sql"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/handler"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/middleware"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/repositories"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/routing"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/service"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/public_lib/atylabencrypt"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/public_lib/atylabjwt"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabcsrf"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	db         *gorm.DB
	middleware *middleware.Middleware
	gin        *gin.Engine
}

func NewApp(db *gorm.DB, sqlDB *sql.DB) (*App, func(), error) {

	app := &App{
		db: db,
	}

	cleanup := func() { sqlDB.Close() }
	return app, cleanup, nil
}

func (a *App) initMiddlewares() {
	// ミドルウェアの初期化
	a.middleware = middleware.NewMiddleware(a.gin)
}

func (a *App) entryBeforeGlobalMiddleware() {
	// 前処理系ミドルウェアをここに追加
	// a.gin.Use(a.middleware.Firewall)
	a.gin.Use(a.middleware.Csrf)
}

func (a *App) entryAfterGlobalMiddleware() {
	// 後処理系ミドルウェアをここに追加
}

func (a *App) initRoutes() {
	// ルーティングの初期化
	routing := routing.NewRouting(a.gin, a.middleware)
	routing.HealthCheckRoute(handler.NewHealthCheckHandler())
	routing.CsrfRoute(
		handler.NewCSRFHandler(
			service.NewCsrfSvcStruct(
				atylabcsrf.NewCsrfPkgStruct(),
			),
		),
	)
	routing.RegisterRouting(
		handler.NewRegisterHandler(
			service.NewUserRegisterSvc(
				a.db,
				atylabencrypt.NewEncryptPkg(),
			),
		),
	)
	routing.AuthRouting(
		handler.NewAuthHandler(
			service.NewAuthSvc(
				repositories.NewUserRepo(a.db),
				repositories.NewUserRefreshTokenRepo(a.db),
				atylabjwt.NewJwtSvc(),
			),
			repositories.NewUserRepo(a.db),
		),
	)
}

func (a *App) Init(g *gin.Engine) {
	a.gin = g
	a.initMiddlewares()
	a.entryBeforeGlobalMiddleware()
	a.initRoutes()
	a.entryAfterGlobalMiddleware()

}

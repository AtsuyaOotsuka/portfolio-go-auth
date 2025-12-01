package app

import (
	"database/sql"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/middleware"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/provider"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	db         *gorm.DB
	middleware *middleware.Middleware
	provider   *provider.Provider
	gin        *gin.Engine
}

func NewApp(db *gorm.DB, sqlDB *sql.DB) (*App, func(), error) {
	app := &App{
		db: db,
	}

	cleanup := func() { sqlDB.Close() }
	return app, cleanup, nil
}

func (a *App) Init(g *gin.Engine) {
	a.gin = g
	a.initProviders()
	a.initMiddlewares()
	a.entryBeforeGlobalMiddleware()
	a.entryRoutes()
	a.entryAfterGlobalMiddleware()
}

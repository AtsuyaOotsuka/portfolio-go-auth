package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/app"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabdatabase"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func SetupDB() (*gorm.DB, *sql.DB) {
	db_pkg := atylabdatabase.NewDBConnect(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_TZ"),
	)
	db, _ := db_pkg.ConnectDB()

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	return db, sqlDB
}

func SetupRouter(db *gorm.DB, sqlDB *sql.DB) (*gin.Engine, func()) {
	r := gin.New()
	r.Use(gin.Recovery())

	app, cleanup, _ := app.NewApp(db, sqlDB)

	app.Init(r)

	return r, cleanup
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("envファイルの読み込みに失敗しました。")
	}

	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = gin.DebugMode
	}
	gin.SetMode(mode)

	db, sqlDB := SetupDB()
	r, cleanup := SetupRouter(db, sqlDB)
	defer cleanup()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

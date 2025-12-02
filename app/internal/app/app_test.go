package app_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/app"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/test_helper/mocks/global_mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) (*gorm.DB, func()) {
	t.Helper()

	db, _, cleanup := global_mock.NewGormWithMock(t)

	// 必要なら  AutoMigrate(&models.User{})
	return db, cleanup
}

func TestNewAppAndInitRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// App生成とルート初期化
	db, cleanup := newTestDB(t)
	defer cleanup()
	sqlDB, _ := db.DB()
	a, cleanup, err := app.NewApp(db, sqlDB) // ← これで db.DB() もOK
	if err != nil {
		t.Fatal(err)
	}
	defer cleanup()
	a.Init(r)

	// テスト用リクエスト
	req := httptest.NewRequest("GET", "/healthcheck", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// 結果検証
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ok")
}

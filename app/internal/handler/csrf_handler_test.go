package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/test_helper/funcs"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/test_helper/mocks/svc_mock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCSRFHandler(t *testing.T) {

	funcs.WithEnv("CSRF_TOKEN", "test_secret", t, func() {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		csrfSvcMock := new(svc_mock.CsrfSvcMockStruct)
		csrfSvcMock.On(
			"CreateCSRFToken",
			mock.AnythingOfType("int64"),
			"test_secret",
		).Return("mocked_csrf_token")

		handler := NewCSRFHandler(
			csrfSvcMock,
		)
		handler.CsrfGet(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "mocked_csrf_token")
	})
}

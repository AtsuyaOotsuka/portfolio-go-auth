package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/repositories"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/test_helper/funcs"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabcsrf"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var testServer *http.Server
var baseURL string
var db *gorm.DB
var sqlDB *sql.DB
var dbRecords []funcs.DbRecords

func TestMain(m *testing.M) {
	var err error

	gin.SetMode(gin.TestMode)
	db, sqlDB = SetupDB()

	dbRecords, err = funcs.DbCleanup(sqlDB)
	if err != nil {
		panic(err)
	}

	assert.NoError(&testing.T{}, err)

	r, cleanup := SetupRouter(db, sqlDB)
	defer cleanup()

	testServer = &http.Server{
		Addr:    ":8880",
		Handler: r,
	}

	go func() {
		if err := testServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	time.Sleep(200 * time.Millisecond)
	baseURL = "http://localhost:8880"

	fmt.Println("Test server started at", baseURL)

	// 全テスト実行
	exitCode := m.Run()

	// サーバをシャットダウン
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = testServer.Shutdown(ctx)

	os.Exit(exitCode)
}

func createCsrf() string {
	csrf_token := os.Getenv("CSRF_TOKEN")
	nonce := funcs.GenerateCSRFCookieToken(
		csrf_token,
		time.Now().Add(1*time.Hour).Unix(),
	)
	return nonce
}

func request(method string, url string, body io.Reader, t *testing.T) (*http.Response, func() error) {
	csrf := createCsrf()

	client := &http.Client{}
	requestUrl := baseURL + url
	fmt.Println("Request URL:", requestUrl)
	req, err := http.NewRequest(method, requestUrl, body)
	if method != "GET" && err == nil {
		req.Header.Set("Content-Type", "application/json")
	}
	assert.NoError(t, err)
	if method != "GET" {
		req.Header.Set("X-CSRF-Token", csrf)
	}
	resp, err := client.Do(req)
	assert.NoError(t, err)

	return resp, resp.Body.Close
}

func TestHealth(t *testing.T) {
	resp, close := request("GET", "/healthcheck", nil, t)
	defer close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	bodyString := string(bodyBytes)
	var response map[string]string
	err = json.Unmarshal([]byte(bodyString), &response)
	assert.NoError(t, err)
	assert.Equal(t, "ok", response["status"])
}

func TestLogin(t *testing.T) {
	usersData := funcs.FilterRecordsByTableName(dbRecords, "users")
	id := usersData[0].Data[0]["id"].(int64)
	uuid := usersData[0].Data[0]["uuid"].(string)
	email := usersData[0].Data[0]["email"].(string)
	password := usersData[0].Data[0]["password"].(string)

	body := map[string]string{
		"email":    email,
		"password": password,
	}
	jsonBody, _ := json.Marshal(body)
	reqBody := strings.NewReader(string(jsonBody))

	resp, close := request("POST", "/auth/login", reqBody, t)
	defer close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var respData map[string]interface{}
	err = json.Unmarshal(bodyBytes, &respData)
	assert.NoError(t, err)

	assert.NotEmpty(t, respData)
	fmt.Println("access_token", respData["access_token"])

	jwt := respData["access_token"].(string)
	jwtInfo, err := funcs.JwtConvert(jwt)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmt.Sprintf("user%s", uuid), jwtInfo.Uuid)
	assert.Equal(t, email, jwtInfo.Email)

	insertResult := funcs.GetRecords(sqlDB, "user_refresh_tokens", map[string]interface{}{
		"user_id": id,
	})

	assert.Len(t, insertResult, 2)
}

func TestRefresh(t *testing.T) {
	usersData := funcs.FilterRecordsByTableName(dbRecords, "users")
	id := usersData[0].Data[0]["id"].(int64)
	uuid := usersData[0].Data[0]["uuid"].(string)
	email := usersData[0].Data[0]["email"].(string)

	body := map[string]string{
		"refresh_token": "refresh_token_sample" + fmt.Sprintf("%d", id),
	}
	jsonBody, _ := json.Marshal(body)
	reqBody := strings.NewReader(string(jsonBody))
	resp, close := request("POST", "/auth/refresh", reqBody, t)
	defer close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var respData map[string]interface{}
	err = json.Unmarshal(bodyBytes, &respData)
	assert.NoError(t, err)

	assert.NotEmpty(t, respData)

	jwt := respData["access_token"].(string)
	jwtInfo, err := funcs.JwtConvert(jwt)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmt.Sprintf("user%s", uuid), jwtInfo.Uuid)
	assert.Equal(t, email, jwtInfo.Email)

}

func TestRegister(t *testing.T) {
	body := map[string]string{
		"name":     "newuser",
		"email":    "newuser@example.com",
		"password": "newpassword123",
	}
	jsonBody, _ := json.Marshal(body)
	reqBody := strings.NewReader(string(jsonBody))

	resp, close := request("POST", "/register", reqBody, t)
	defer close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var respData map[string]interface{}
	err = json.Unmarshal(bodyBytes, &respData)
	assert.NoError(t, err)

	assert.NotEmpty(t, respData)
	assert.Equal(t, "newuser", respData["username"])
	assert.Equal(t, "newuser@example.com", respData["email"])

	exists := funcs.ExistsRecord(sqlDB, "users", map[string]interface{}{
		"email":    "newuser@example.com",
		"username": "newuser",
	})
	assert.True(t, exists)

	user := repositories.NewUserRepo(db)
	dbUser, err := user.GetByEmail("newuser@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, dbUser)

	passCheck := bcrypt.CompareHashAndPassword(
		[]byte(dbUser.PasswordHash),
		[]byte("newpassword123"),
	)
	assert.NoError(t, passCheck)
}

func TestCsrfGet(t *testing.T) {
	resp, close := request("GET", "/csrf/get", nil, t)
	defer close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	var respData map[string]string
	err = json.Unmarshal(bodyBytes, &respData)
	assert.NoError(t, err)

	assert.NotEmpty(t, respData)
	assert.NotEmpty(t, respData["csrf_token"])

	cookies := resp.Cookies()
	var csrfCookie *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "csrf_token" {
			csrfCookie = cookie
			break
		}
	}
	assert.NotNil(t, csrfCookie)
	assert.NotEmpty(t, csrfCookie.Value)

	// Cookieから出したため、URLデコードを実施
	decodeToken, err := url.QueryUnescape(csrfCookie.Value)
	assert.NoError(t, err)

	// クッキーの値とレスポンスの値が同じであることを確認
	assert.Equal(t, decodeToken, respData["csrf_token"])

	csrfStruct := atylabcsrf.NewCsrfPkgStruct()
	validErr := csrfStruct.ValidateCSRFCookieToken(
		decodeToken,
		os.Getenv("CSRF_TOKEN"),
		time.Now().Unix(),
	)
	assert.NoError(t, validErr)
}

//go:build integration

package middleware_test

import (
	"encoding/json"
	"github.com/t-kuni/go-web-api-template/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRecover_パニックが発生した場合ステータスコード500を返しエラーログが出力されること(t *testing.T) {
	//
	// Prepare
	//
	cont := beforeEach(t)
	defer afterEach(cont)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test-path", nil)
	req.Header.Set("Referer", "https://example.com/")
	req.Header.Set("Test-Header-Key1", "Test-Header-Value-1")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &model.User{
		ID:       10000,
		Username: "test",
		Name:     "テストユーザ",
		Email:    "test@test.test",
		Role:     "administrator",
	})
	mockHandler := func(c echo.Context) error { panic(errors.New("テストpanicです")) }

	//
	// Execute
	//
	m := middleware.NewRecover()
	err := m.Recover()(mockHandler)(c)

	//
	// Assert
	//
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	assert.Equal(t, 1, len(loggerHook.Entries))

	logStr, err := loggerHook.LastEntry().String()
	if err != nil {
		t.Fatal(err)
	}

	var log map[string]interface{}
	err = json.Unmarshal([]byte(logStr), &log)
	if err != nil {
		t.Fatal(err)
	}

	assert.Regexp(t, "テストpanicです", log["message"])
	assert.Regexp(t, "middleware/recover.go:[0-9]+", log["message"])
	assert.Equal(t, "error", log["level"])
	assert.Nil(t, log["params"])
	assert.NotEmpty(t, log["function"])
	assert.NotEmpty(t, log["file"])
	assert.NotEmpty(t, log["line"])
	assert.NotEmpty(t, log["host"])
	assert.Equal(t, "/test-path", log["uri"])
	assert.Equal(t, "192.0.2.1", log["ip"])
	assert.Equal(t, "GET", log["http_method"])
	assert.NotEmpty(t, log["server"])
	assert.Equal(t, "https://example.com/", log["referrer"])
	assert.Equal(t, "test", log["environment"])
	assert.Nil(t, log["input"])
	headers := log["header"].(map[string]interface{})
	assert.Len(t, headers, 2)
	assert.Equal(t, "Test-Header-Value-1", headers["Test-Header-Key1"].([]interface{})[0])
	assert.Equal(t, 10000, int(log["userId"].(float64))) // floatとしてパースされてしまうのでintに直す
	assert.Equal(t, true, log["panic"])
}

func TestRecover_すでにレスポンスを返却済みのタイミングでパニックが発生した場合ログだけ記録されること(t *testing.T) {
	//
	// Prepare
	//
	cont := beforeEach(t)
	defer afterEach(cont)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test-path", nil)
	req.Header.Set("Referer", "https://example.com/")
	req.Header.Set("Test-Header-Key1", "Test-Header-Value-1")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &model.User{
		ID:       10000,
		Username: "test",
		Name:     "テストユーザ",
		Email:    "test@test.test",
		Role:     "administrator",
	})
	mockHandler := func(c echo.Context) error {
		_ = c.NoContent(http.StatusBadRequest) // レスポンス返却済み
		panic(errors.New("テストpanicです"))
	}

	//
	// Execute
	//
	m := middleware.NewRecover()
	err := m.Recover()(mockHandler)(c)

	//
	// Assert
	//
	if err != nil {
		t.Error(err)
	}

	// handlerで返却したレスポンスが返却されていること
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	assert.Equal(t, 1, len(loggerHook.Entries))

	logStr, err := loggerHook.LastEntry().String()
	if err != nil {
		t.Fatal(err)
	}

	var log map[string]interface{}
	err = json.Unmarshal([]byte(logStr), &log)
	if err != nil {
		t.Fatal(err)
	}

	assert.Regexp(t, "テストpanicです", log["message"])
	assert.Equal(t, true, log["panic"])
}

func TestRecover_パニックが発生しない場合エラーログが出力されないこと(t *testing.T) {
	//
	// Prepare
	//
	cont := beforeEach(t)
	defer afterEach(cont)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/test-path", nil)
	req.Header.Set("Referer", "https://example.com/")
	req.Header.Set("Test-Header-Key1", "Test-Header-Value-1")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &model.User{
		ID:       10000,
		Username: "test",
		Name:     "テストユーザ",
		Email:    "test@test.test",
		Role:     "administrator",
	})
	mockHandler := func(c echo.Context) error { return nil }

	//
	// Execute
	//
	m := middleware.NewRecover()
	err := m.Recover()(mockHandler)(c)

	//
	// Assert
	//
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 0, len(loggerHook.Entries))
}

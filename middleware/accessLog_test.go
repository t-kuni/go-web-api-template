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

func TestAccessLog_RequestログとResponseログが出力されること(t *testing.T) {
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
	mockHandler := func(c echo.Context) error { return c.NoContent(http.StatusOK) }

	//
	// Execute
	//
	m := middleware.NewAccessLog()
	err := m.AccessLog(mockHandler)(c)

	//
	// Assert
	//
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 2, len(loggerHook.Entries))

	// リクエストログの検証
	reqLogStr, err := loggerHook.Entries[0].String()
	if err != nil {
		t.Fatal(err)
	}

	var reqLog map[string]interface{}
	err = json.Unmarshal([]byte(reqLogStr), &reqLog)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "[Request]", reqLog["message"])
	assert.Equal(t, "info", reqLog["level"])
	assert.Nil(t, reqLog["params"])
	assert.NotEmpty(t, reqLog["function"])
	assert.NotEmpty(t, reqLog["file"])
	assert.NotEmpty(t, reqLog["line"])
	assert.NotEmpty(t, reqLog["host"])
	assert.Equal(t, "/test-path", reqLog["uri"])
	assert.Equal(t, "192.0.2.1", reqLog["ip"])
	assert.Equal(t, "GET", reqLog["http_method"])
	assert.NotEmpty(t, reqLog["server"])
	assert.Equal(t, "https://example.com/", reqLog["referrer"])
	assert.Equal(t, "test", reqLog["environment"])
	assert.Nil(t, reqLog["input"])
	{
		headers := reqLog["header"].(map[string]interface{})
		assert.Len(t, headers, 2)
		assert.Equal(t, "Test-Header-Value-1", headers["Test-Header-Key1"].([]interface{})[0])
	}

	// レスポンスログの検証
	respLogStr, err := loggerHook.Entries[1].String()
	if err != nil {
		t.Fatal(err)
	}

	var respLog map[string]interface{}
	err = json.Unmarshal([]byte(respLogStr), &respLog)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "[Response]", respLog["message"])
	assert.Equal(t, "info", respLog["level"])
	assert.Nil(t, respLog["params"])
	assert.NotEmpty(t, respLog["function"])
	assert.NotEmpty(t, respLog["file"])
	assert.NotEmpty(t, respLog["line"])
	assert.NotEmpty(t, respLog["host"])
	assert.Equal(t, "/test-path", respLog["uri"])
	assert.Equal(t, "192.0.2.1", respLog["ip"])
	assert.Equal(t, "GET", respLog["http_method"])
	assert.NotEmpty(t, respLog["server"])
	assert.Equal(t, "https://example.com/", respLog["referrer"])
	assert.Equal(t, "test", respLog["environment"])
	assert.Nil(t, respLog["input"])
	{
		headers := respLog["header"].(map[string]interface{})
		assert.Len(t, headers, 2)
		assert.Equal(t, "Test-Header-Value-1", headers["Test-Header-Key1"].([]interface{})[0])
	}
	assert.Equal(t, 10000, int(respLog["userId"].(float64))) // floatとしてパースされてしまうのでintに直す
	assert.NotNil(t, respLog["latency"])
	assert.NotNil(t, respLog["latency_human"])
}

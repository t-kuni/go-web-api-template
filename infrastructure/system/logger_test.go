package system_test

import (
	"encoding/json"
	"fmt"
	"github.com/rotisserie/eris"
	"github.com/stretchr/testify/assert"
	"github.com/t-kuni/go-web-api-template/infrastructure/system"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	t.Run("Should output info log", func(t *testing.T) {
		os.Setenv("APP_ENV", "test")

		req := httptest.NewRequest(http.MethodGet, "/test-path", nil)
		req.Header.Set("Referer", "https://example.com/")
		req.Header.Set("Test-Header-Key1", "Test-Header-Value-1")
		req.Header.Set("X-Forwarded-For", "192.0.2.1")

		logger, loggerHook := system.NewTestLogger()
		logger.Info(req, "test message", map[string]interface{}{
			"testKey1": "testValue1",
		})

		assert.Equal(t, 1, len(loggerHook.Entries))

		logStr, err := loggerHook.LastEntry().String()
		assert.NoError(t, err)

		var log map[string]interface{}
		err = json.Unmarshal([]byte(logStr), &log)
		assert.NoError(t, err)

		assert.Equal(t, "test message", log["message"])
		assert.Equal(t, "info", log["level"])
		assert.Equal(t, "testValue1", log["params"].(map[string]interface{})["testKey1"])
		assert.NotEmpty(t, log["function"])
		assert.NotEmpty(t, log["file"])
		assert.NotEmpty(t, log["line"])
		assert.NotEmpty(t, log["host"])
		assert.Equal(t, "/test-path", log["uri"])
		assert.Equal(t, "192.0.2.1", log["ip"])
		assert.Equal(t, "GET", log["http_method"])
		assert.NotEmpty(t, log["server_ip"])
		assert.Equal(t, "https://example.com/", log["referrer"])
		assert.Equal(t, "test", log["environment"])
		assert.Nil(t, log["input"])
		headers := log["header"].(map[string]interface{})
		assert.Len(t, headers, 3)
		assert.Equal(t, "Test-Header-Value-1", headers["Test-Header-Key1"].([]interface{})[0])
	})

	t.Run("Should output error log with stack trace", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test-path", nil)
		req.Header.Set("Referer", "https://example.com/")
		req.Header.Set("Test-Header-Key1", "Test-Header-Value-1")
		req.Header.Set("X-Forwarded-For", "192.0.2.1")

		err := fmt.Errorf("test root error")
		wrappedErr := eris.Wrap(err, "wrapped error")
		logger, loggerHook := system.NewTestLogger()
		logger.Error(req, wrappedErr, map[string]interface{}{
			"testKey1": "testValue1",
		})

		assert.Equal(t, 1, len(loggerHook.Entries))

		logStr, err := loggerHook.LastEntry().String()
		assert.NoError(t, err)

		var log map[string]interface{}
		err = json.Unmarshal([]byte(logStr), &log)
		assert.NoError(t, err)

		assert.Regexp(t, "test root error", log["message"])
		assert.Regexp(t, "wrapped error", log["message"])
		assert.Regexp(t, "logger_test.go:[0-9]+", log["message"])
		assert.Equal(t, "error", log["level"])
		assert.Equal(t, "testValue1", log["params"].(map[string]interface{})["testKey1"])
		assert.NotEmpty(t, log["function"])
		assert.NotEmpty(t, log["file"])
		assert.NotEmpty(t, log["line"])
		assert.NotEmpty(t, log["host"])
		assert.Equal(t, "/test-path", log["uri"])
		assert.Equal(t, "192.0.2.1", log["ip"])
		assert.Equal(t, "GET", log["http_method"])
		assert.NotEmpty(t, log["server_ip"])
		assert.Equal(t, "https://example.com/", log["referrer"])
		assert.Equal(t, "test", log["environment"])
		assert.Nil(t, log["input"])
		headers := log["header"].(map[string]interface{})
		assert.Len(t, headers, 3)
		assert.Equal(t, "Test-Header-Value-1", headers["Test-Header-Key1"].([]interface{})[0])
		assert.Equal(t, false, log["panic"])
	})

	t.Run("Should output warn log with stack trace", func(t *testing.T) {
		reqBody := "{ testKey1: 'testValue1' }"
		req := httptest.NewRequest(http.MethodGet, "/test-path", strings.NewReader(reqBody))
		req.Header.Set("X-Forwarded-For", "192.0.2.1")

		err := fmt.Errorf("test root error")
		wrappedErr := eris.Wrap(err, "wrapped error")
		// テスト用のロガーを使用
		logger, loggerHook := system.NewTestLogger()
		logger.WarnWithError(req, wrappedErr, map[string]interface{}{
			"testKey1": "testValue1",
		})

		assert.Equal(t, 1, len(loggerHook.Entries))

		logStr, err := loggerHook.LastEntry().String()
		assert.NoError(t, err)

		var log map[string]interface{}
		err = json.Unmarshal([]byte(logStr), &log)
		assert.NoError(t, err)

		assert.Regexp(t, "test root error", log["message"])
		assert.Regexp(t, "wrapped error", log["message"])
		assert.Regexp(t, "logger_test.go:[0-9]+", log["message"])
		assert.Equal(t, "warning", log["level"])
	})
}

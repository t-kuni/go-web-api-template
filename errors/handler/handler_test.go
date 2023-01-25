package handler_test

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/t-kuni/go-web-api-template/errors/handler"
	"github.com/t-kuni/go-web-api-template/errors/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorHandler(t *testing.T) {
	t.Run("Should return status code 500", func(t *testing.T) {
		cont := beforeEach(t)
		defer afterEach(cont)

		err := fmt.Errorf("standard error")
		rec := callErrorHandler(t, err)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, 0, rec.Body.Len())
		assertLogLevel(t, "error")
	})

	t.Run("Should return status code 500 when wrapped error", func(t *testing.T) {
		cont := beforeEach(t)
		defer afterEach(cont)

		err := eris.Wrap(fmt.Errorf("standard error"), "")
		rec := callErrorHandler(t, err)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, 0, rec.Body.Len())
		assertLogLevel(t, "error")
	})

	t.Run("Should return status code 415 when bind error", func(t *testing.T) {
		cont := beforeEach(t)
		defer afterEach(cont)

		err := echo.NewBindingError("", []string{""}, "", fmt.Errorf(""))
		err = eris.Wrap(err, "")
		rec := callErrorHandler(t, err)

		assert.Equal(t, http.StatusUnsupportedMediaType, rec.Code)
		assert.Equal(t, 0, rec.Body.Len())
		assertLogLevel(t, "warning")
	})

	t.Run("Should return status code 400 when validation error", func(t *testing.T) {
		cont := beforeEach(t)
		defer afterEach(cont)

		err := eris.Wrap(validator.ValidationErrors{}, "")
		rec := callErrorHandler(t, err)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		expectBody := `""
`
		assert.Equal(t, expectBody, rec.Body.String())
		assertLogLevel(t, "warning")
	})

	t.Run("Should record a log if an error occurs after responded.", func(t *testing.T) {
		cont := beforeEach(t)
		defer afterEach(cont)

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Already responded
		if e := c.NoContent(http.StatusOK); e != nil {
			t.Fatal(e)
		}

		err := fmt.Errorf("standard error")
		dummyInjector := do.New()
		h, err := handler.NewErrorHandler(dummyInjector)
		assert.NoError(t, err)
		h.Handler(err, c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, 0, rec.Body.Len())
		assertLogLevel(t, "error")
	})

	t.Run("Should return status code and message held in the error when HTTPError", func(t *testing.T) {
		cont := beforeEach(t)
		defer afterEach(cont)

		err := eris.Wrap(echo.NewHTTPError(http.StatusNotFound, "dummy"), "")
		rec := callErrorHandler(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		expectBody := `{"message": "dummy"}`
		assert.JSONEq(t, expectBody, rec.Body.String())
		assertLogLevel(t, "warning")
	})

	t.Run("Should return status code and body held in the error when HTTPError", func(t *testing.T) {
		cont := beforeEach(t)
		defer afterEach(cont)

		err := eris.Wrap(echo.NewHTTPError(http.StatusNotFound, map[string]string{"aaa": "bbb"}), "")
		rec := callErrorHandler(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		expectBody := `{"aaa": "bbb"}`
		assert.JSONEq(t, expectBody, rec.Body.String())
		assertLogLevel(t, "warning")
	})

	t.Run("Should return status code and body held in the error when HTTPError is nested", func(t *testing.T) {
		cont := beforeEach(t)
		defer afterEach(cont)

		httpErr := echo.NewHTTPError(http.StatusForbidden, map[string]string{"ccc": "ddd"})
		httpErr.Internal = echo.NewHTTPError(http.StatusNotFound, map[string]string{"aaa": "bbb"})
		err := eris.Wrap(httpErr, "")
		rec := callErrorHandler(t, err)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		expectBody := `{"aaa": "bbb"}`
		assert.JSONEq(t, expectBody, rec.Body.String())
		assertLogLevel(t, "warning")
	})

	t.Run("Should return status code and body held in the error when HTTPError holds another error", func(t *testing.T) {
		cont := beforeEach(t)
		defer afterEach(cont)

		httpErr := echo.NewHTTPError(http.StatusForbidden, map[string]string{"ccc": "ddd"})
		httpErr.Internal = &json.UnmarshalTypeError{}
		err := eris.Wrap(httpErr, "")
		rec := callErrorHandler(t, err)

		assert.Equal(t, http.StatusForbidden, rec.Code)
		expectBody := `{"ccc": "ddd"}`
		assert.JSONEq(t, expectBody, rec.Body.String())
		assertLogLevel(t, "warning")
	})

	t.Run("Should return status code 422 when BasicBusinessError", func(t *testing.T) {
		cont := beforeEach(t)
		defer afterEach(cont)

		appErr := eris.Wrap(types.NewBasicBusinessError("test message", map[string]interface{}{
			"testKey1": "testValue1",
		}), "")
		rec := callErrorHandler(t, appErr)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

		// Assert response
		expectBody := `{"message": "test message"}`
		assert.JSONEq(t, expectBody, rec.Body.String())

		// Assert log
		assert.Len(t, loggerHook.Entries, 1)

		logStr, err := loggerHook.LastEntry().String()
		if err != nil {
			t.Fatal(err)
		}

		var log map[string]interface{}
		err = json.Unmarshal([]byte(logStr), &log)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "warning", log["level"])
		assert.Equal(t, "testValue1", log["params"].(map[string]interface{})["testKey1"])
	})
}

func callErrorHandler(t *testing.T, targetErr error) *httptest.ResponseRecorder {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	dummyInjector := do.New()
	h, err := handler.NewErrorHandler(dummyInjector)
	assert.NoError(t, err)

	h.Handler(targetErr, c)

	return rec
}

func assertLogLevel(t *testing.T, logLevel string) {
	assert.Len(t, loggerHook.Entries, 1)

	logStr, err := loggerHook.LastEntry().String()
	if err != nil {
		t.Fatal(err)
	}

	var log map[string]interface{}
	err = json.Unmarshal([]byte(logStr), &log)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, logLevel, log["level"])
}

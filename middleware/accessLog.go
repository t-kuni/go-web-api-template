package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/go-http-utils/headers"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"github.com/t-kuni/go-web-api-template/logger"
	"io/ioutil"
	"net/http"
	"time"
)

// AccessLog RequestログとResponseログを出力するミドルウェア
type AccessLog struct {
}

func NewAccessLog(i *do.Injector) (*AccessLog, error) {
	return &AccessLog{}, nil
}

type CustomResponseWriter struct {
	wrapped    http.ResponseWriter
	StatusCode int
}

func NewCustomResponseWriter(wrapped http.ResponseWriter) *CustomResponseWriter {
	return &CustomResponseWriter{wrapped: wrapped}
}

func (w *CustomResponseWriter) Header() http.Header {
	return w.wrapped.Header()
}

func (w *CustomResponseWriter) Write(content []byte) (int, error) {
	return w.wrapped.Write(content)
}

func (w *CustomResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.wrapped.WriteHeader(statusCode)
}

func (m AccessLog) AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqBody := getRequestBody(r)
		logger.RequestLogV2(r, reqBody)

		respWriter := NewCustomResponseWriter(w)
		latency, latencyHuman := measureLatency(func() {
			next.ServeHTTP(respWriter, r)
		})

		logger.ResponseLogV2(r, respWriter.StatusCode, latency, latencyHuman)
	})
}

func getRequestBody(req *http.Request) map[string]interface{} {
	var reqBody map[string]interface{}
	contentType := req.Header.Get(headers.ContentType)
	if contentType == echo.MIMEApplicationJSON {
		if req.Body != nil {
			reqBodyBytes, err := ioutil.ReadAll(req.Body)
			req.Body = ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes)) // Reset  see detail: https://stackoverflow.com/a/47295689
			if err == nil {
				json.Unmarshal(reqBodyBytes, &reqBody)
			}
		}
	}
	return reqBody
}

func measureLatency(proc func()) (latency time.Duration, latencyHuman string) {
	start := time.Now()
	proc()
	stop := time.Now()

	latency = stop.Sub(start)
	latencyHuman = latency.String()
	return
}

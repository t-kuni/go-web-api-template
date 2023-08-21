package logger

import (
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"
)

type (
	StackInfo struct {
		file     string
		line     int
		funcName string
	}
)

func SetupLogger() error {
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg: "message",
		},
	})
	logrus.SetOutput(os.Stdout)

	level, err := getLogLevel()
	if err != nil {
		return eris.Wrap(err, "")
	}
	logrus.SetLevel(level)

	return nil
}

func Info(c echo.Context, msg string, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	logrus.
		WithFields(makeCommonFields(stackInfo, params)).
		WithFields(makeHttpFields(c)).
		Info(msg)
}

func SimpleInfoF(format string, args ...interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	logrus.
		WithFields(makeCommonFields(stackInfo, nil)).
		Infof(format, args...)
}

func Warn(c echo.Context, msg string, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	logrus.
		WithFields(makeCommonFields(stackInfo, params)).
		WithFields(makeHttpFields(c)).
		Warn(msg)
}

func WarnWithError(c echo.Context, e error, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	logrus.
		WithFields(makeCommonFields(stackInfo, params)).
		WithFields(makeHttpFields(c)).
		WithField("panic", false).
		Warnf("%+v", e)
}

func Error(c echo.Context, e error, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	logrus.
		WithFields(makeCommonFields(stackInfo, params)).
		WithFields(makeHttpFields(c)).
		WithField("panic", false).
		Errorf("%+v", e)
}

func Debug(c echo.Context, msg string, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	logrus.
		WithFields(makeCommonFields(stackInfo, params)).
		WithFields(makeHttpFields(c)).
		Debug(msg)
}

func Fatal(c echo.Context, msg string, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	logrus.
		WithFields(makeCommonFields(stackInfo, params)).
		WithFields(makeHttpFields(c)).
		Fatal(msg)
}

func SimpleFatal(e error, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	logrus.
		WithFields(makeCommonFields(stackInfo, params)).
		Fatalf("%+v", e)
}

func Panic(c echo.Context, msg string, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	logrus.
		WithFields(makeCommonFields(stackInfo, params)).
		WithFields(makeHttpFields(c)).
		WithField("panic", true).
		Error(msg)
}

func PanicV2(req *http.Request, msg string, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	logrus.
		WithFields(makeCommonFields(stackInfo, params)).
		WithFields(makeHttpFieldsV2(req)).
		WithField("panic", true).
		Error(msg)
}

func RequestLog(c echo.Context) {
	stackInfo := makeStackInfo(runtime.Caller(1))

	url := c.Request().RequestURI
	method := c.Request().Method
	msg := fmt.Sprintf("[Request][%s]%s", url, method)

	logrus.
		WithFields(makeCommonFields(stackInfo, nil)).
		WithFields(makeHttpFields(c)).
		Info(msg)
}

func RequestLogV2(req *http.Request, reqBody map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))

	url := req.RequestURI
	method := req.Method
	msg := fmt.Sprintf("[Request][%s]%s", url, method)

	logrus.
		WithFields(makeCommonFields(stackInfo, nil)).
		WithFields(makeHttpFieldsV2(req)).
		WithField("input", reqBody).
		Info(msg)
}

func ResponseLog(c echo.Context, status int, latency time.Duration, latencyHuman string) {
	stackInfo := makeStackInfo(runtime.Caller(1))

	url := c.Request().RequestURI
	method := c.Request().Method
	msg := fmt.Sprintf("[Response][%s]%s", url, method)

	logrus.
		WithFields(makeCommonFields(stackInfo, nil)).
		WithFields(makeHttpFields(c)).
		WithField("latency", latency).
		WithField("latency_human", latencyHuman).
		WithField("http_status", status).
		Info(msg)
}

func ResponseLogV2(req *http.Request, status int, latency time.Duration, latencyHuman string) {
	stackInfo := makeStackInfo(runtime.Caller(1))

	url := req.RequestURI
	method := req.Method
	msg := fmt.Sprintf("[Response][%s]%s", url, method)

	logrus.
		WithFields(makeCommonFields(stackInfo, nil)).
		WithFields(makeHttpFieldsV2(req)).
		WithField("latency", latency).
		WithField("latency_human", latencyHuman).
		WithField("http_status", status).
		Info(msg)
}

func getLogLevel() (logrus.Level, error) {
	levelStr := os.Getenv("LOG_LEVEL")

	if levelStr == "" {
		return logrus.InfoLevel, nil
	}

	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		return 0, eris.Wrap(err, "")
	}

	return level, nil
}

func makeCommonFields(stackInfo *StackInfo, params map[string]interface{}) map[string]interface{} {
	var function *string
	var file *string
	var line *int
	if stackInfo != nil {
		function = &stackInfo.funcName
		file = &stackInfo.file
		line = &stackInfo.line
	}

	hostname, _ := os.Hostname()

	return map[string]interface{}{
		"params":   params,
		"function": function,
		"file":     file,
		"line":     line,
		"host":     hostname,
	}
}

func makeHttpFields(c echo.Context) map[string]interface{} {
	req := c.Request()

	return map[string]interface{}{
		"uri":         req.RequestURI,
		"ip":          c.RealIP(),
		"http_method": req.Method,
		"server_ip":   getLocalIP(),
		"referrer":    req.Referer(),
		"environment": os.Getenv("APP_ENV"),
		"header":      makeHeaderField(c),
	}
}

func makeHttpFieldsV2(req *http.Request) map[string]interface{} {
	return map[string]interface{}{
		"uri":         req.RequestURI,
		"ip":          req.Header.Get(headers.XForwardedFor),
		"http_method": req.Method,
		"server_ip":   getLocalIP(),
		"referrer":    req.Referer(),
		"environment": os.Getenv("APP_ENV"),
		"header":      makeHeaderFieldV2(req),
	}
}

func makeHeaderField(c echo.Context) map[string]interface{} {
	excludeHeaders := []string{
		"Authorization",
	}
	return filterHeaders(c.Request().Header, excludeHeaders)
}

func makeHeaderFieldV2(req *http.Request) map[string]interface{} {
	excludeHeaders := []string{
		"Authorization",
	}
	return filterHeaders(req.Header, excludeHeaders)
}

func makeStackInfo(pc uintptr, file string, line int, ok bool) *StackInfo {
	if !ok {
		return nil
	}

	funcName := runtime.FuncForPC(pc).Name()
	return &StackInfo{
		file:     file,
		line:     line,
		funcName: funcName,
	}
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func filterHeaders(headers http.Header, excludeHeaders []string) map[string]interface{} {
	if headers == nil {
		return nil
	}

	filteredHeaders := make(map[string]interface{})
	for k, v := range headers {
		if !includes(k, excludeHeaders) {
			filteredHeaders[k] = v
		}
	}
	return filteredHeaders
}

func includes(target string, arr []string) bool {
	for _, v := range arr {
		if target == v {
			return true
		}
	}
	return false
}

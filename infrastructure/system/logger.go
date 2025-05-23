package system

import (
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/rotisserie/eris"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	systemInterface "github.com/t-kuni/go-web-api-template/domain/infrastructure/system"
	"io"
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

	Logger struct {
		logger *logrus.Logger
	}
)

// NewLogger は新しいロガーインスタンスを生成します
func NewLogger() (systemInterface.ILogger, error) {
	logger := logrus.New()

	// ロガーの設定
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg: "message",
		},
	})
	logger.SetOutput(os.Stdout)

	// ログレベルの設定
	level, err := getLogLevel()
	if err != nil {
		return nil, eris.Wrap(err, "")
	}
	logger.SetLevel(level)

	return &Logger{
		logger: logger,
	}, nil
}

// NewTestLogger はテスト用の新しいロガーインスタンスを生成します
func NewTestLogger() (systemInterface.ILogger, *test.Hook) {
	testLogger, loggerHook := test.NewNullLogger()
	testLogger.SetLevel(logrus.TraceLevel)

	// JSONフォーマッタを設定
	testLogger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg: "message",
		},
	})

	// 出力先を設定（テスト時は出力を破棄）
	testLogger.SetOutput(io.Discard)

	return &Logger{
		logger: testLogger,
	}, loggerHook
}

// NewLoggerWithCustomLogger は指定されたlogrusロガーを使用して新しいロガーインスタンスを生成します
// テスト用途に使用します
func NewLoggerWithCustomLogger(customLogger *logrus.Logger) systemInterface.ILogger {
	return &Logger{
		logger: customLogger,
	}
}

func (l *Logger) Info(req *http.Request, msg string, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	l.logger.
		WithFields(makeCommonFields(stackInfo, params)).
		WithFields(makeHttpFieldsV2(req)).
		Info(msg)
}

func (l *Logger) SimpleInfoF(format string, args ...interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	l.logger.
		WithFields(makeCommonFields(stackInfo, nil)).
		Infof(format, args...)
}

func (l *Logger) Warn(req *http.Request, msg string, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	l.logger.
		WithFields(makeCommonFields(stackInfo, params)).
		WithFields(makeHttpFieldsV2(req)).
		Warn(msg)
}

func (l *Logger) WarnWithError(req *http.Request, e error, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	l.logger.
		WithFields(makeCommonFields(stackInfo, params)).
		WithFields(makeHttpFieldsV2(req)).
		WithField("panic", false).
		Warnf("%+v", e)
}

func (l *Logger) Error(req *http.Request, e error, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	l.logger.
		WithFields(makeCommonFields(stackInfo, params)).
		WithFields(makeHttpFieldsV2(req)).
		WithField("panic", false).
		Errorf("%+v", e)
}

func (l *Logger) Debug(req *http.Request, msg string, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	l.logger.
		WithFields(makeCommonFields(stackInfo, params)).
		WithFields(makeHttpFieldsV2(req)).
		Debug(msg)
}

func (l *Logger) Fatal(req *http.Request, msg string, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	l.logger.
		WithFields(makeCommonFields(stackInfo, params)).
		WithFields(makeHttpFieldsV2(req)).
		Fatal(msg)
}

func (l *Logger) SimpleFatal(e error, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	l.logger.
		WithFields(makeCommonFields(stackInfo, params)).
		Fatalf("%+v", e)
}

func (l *Logger) Panic(req *http.Request, msg string, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	l.logger.
		WithFields(makeCommonFields(stackInfo, params)).
		WithFields(makeHttpFieldsV2(req)).
		WithField("panic", true).
		Error(msg)
}

func (l *Logger) PanicV2(req *http.Request, msg string, params map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))
	l.logger.
		WithFields(makeCommonFields(stackInfo, params)).
		WithFields(makeHttpFieldsV2(req)).
		WithField("panic", true).
		Error(msg)
}

func (l *Logger) RequestLog(req *http.Request) {
	stackInfo := makeStackInfo(runtime.Caller(1))

	url := req.RequestURI
	method := req.Method
	msg := fmt.Sprintf("[Request][%s]%s", url, method)

	l.logger.
		WithFields(makeCommonFields(stackInfo, nil)).
		WithFields(makeHttpFieldsV2(req)).
		Info(msg)
}

func (l *Logger) RequestLogV2(req *http.Request, reqBody map[string]interface{}) {
	stackInfo := makeStackInfo(runtime.Caller(1))

	url := req.RequestURI
	method := req.Method
	msg := fmt.Sprintf("[Request][%s]%s", url, method)

	l.logger.
		WithFields(makeCommonFields(stackInfo, nil)).
		WithFields(makeHttpFieldsV2(req)).
		WithField("input", reqBody).
		Info(msg)
}

func (l *Logger) ResponseLog(req *http.Request, status int, latency time.Duration, latencyHuman string) {
	stackInfo := makeStackInfo(runtime.Caller(1))

	url := req.RequestURI
	method := req.Method
	msg := fmt.Sprintf("[Response][%s]%s", url, method)

	l.logger.
		WithFields(makeCommonFields(stackInfo, nil)).
		WithFields(makeHttpFieldsV2(req)).
		WithField("latency", latency).
		WithField("latency_human", latencyHuman).
		WithField("http_status", status).
		Info(msg)
}

func (l *Logger) ResponseLogV2(req *http.Request, status int, latency time.Duration, latencyHuman string) {
	stackInfo := makeStackInfo(runtime.Caller(1))

	url := req.RequestURI
	method := req.Method
	msg := fmt.Sprintf("[Response][%s]%s", url, method)

	l.logger.
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

func makeHttpFieldsV2(req *http.Request) map[string]interface{} {
	if req == nil {
		return nil
	}

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

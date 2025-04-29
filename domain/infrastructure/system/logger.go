//go:generate go tool mockgen -source=$GOFILE -destination=${GOFILE}_mock.go -package=$GOPACKAGE

package system

import (
	"net/http"
	"time"
)

// ILogger ロガーインターフェース
// ログ出力に関する機能を提供します
// アプリケーション内でのログ出力に使用することを想定しています
type ILogger interface {
	// Info 情報ログを出力します
	Info(req *http.Request, msg string, params map[string]interface{})

	// SimpleInfoF フォーマット付き情報ログを出力します
	SimpleInfoF(format string, args ...interface{})

	// Warn 警告ログを出力します
	Warn(req *http.Request, msg string, params map[string]interface{})

	// WarnWithError エラー情報付き警告ログを出力します
	WarnWithError(req *http.Request, e error, params map[string]interface{})

	// Error エラーログを出力します
	Error(req *http.Request, e error, params map[string]interface{})

	// Debug デバッグログを出力します
	Debug(req *http.Request, msg string, params map[string]interface{})

	// Fatal 致命的エラーログを出力します
	Fatal(req *http.Request, msg string, params map[string]interface{})

	// SimpleFatal シンプルな致命的エラーログを出力します
	SimpleFatal(e error, params map[string]interface{})

	// Panic パニックログを出力します
	Panic(req *http.Request, msg string, params map[string]interface{})

	// PanicV2 パニックログを出力します（V2）
	PanicV2(req *http.Request, msg string, params map[string]interface{})

	// RequestLog リクエストログを出力します
	RequestLog(req *http.Request)

	// RequestLogV2 リクエストログを出力します（V2）
	RequestLogV2(req *http.Request, reqBody map[string]interface{})

	// ResponseLog レスポンスログを出力します
	ResponseLog(req *http.Request, status int, latency time.Duration, latencyHuman string)

	// ResponseLogV2 レスポンスログを出力します（V2）
	ResponseLogV2(req *http.Request, status int, latency time.Duration, latencyHuman string)
}

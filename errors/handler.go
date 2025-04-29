package errors

import (
	"github.com/t-kuni/go-web-api-template/infrastructure/system"
	"net/http"
)

// DIコンテナからLoggerを取得するためのグローバル変数
var logger *system.Logger

// SetLogger はDIコンテナからLoggerを設定するための関数です
func SetLogger(l *system.Logger) {
	logger = l
}

// CustomServeError はエラーをログに出力し、500エラーを返すカスタムエラーハンドラです
func CustomServeError(w http.ResponseWriter, r *http.Request, err error) {
	// エラーをログに出力
	if logger != nil {
		logger.PanicV2(r, err.Error(), nil)
	}

	// 常に500エラーを返す
	w.WriteHeader(http.StatusInternalServerError)
}

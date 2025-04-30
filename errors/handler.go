package errors

import (
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/system"
	"net/http"
)

// NewCustomServeError はカスタムエラーハンドラを生成する関数です
func NewCustomServeError(logger system.ILogger) func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		// エラーをログに出力
		if logger != nil {
			logger.PanicV2(r, err.Error(), nil)
		}

		// 常に500エラーを返す
		w.WriteHeader(http.StatusInternalServerError)
	}
}

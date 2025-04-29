package errors

import (
	"github.com/t-kuni/go-web-api-template/logger"
	"net/http"
)

// CustomServeError はエラーをログに出力し、500エラーを返すカスタムエラーハンドラです
func CustomServeError(w http.ResponseWriter, r *http.Request, err error) {
	// エラーをログに出力
	logger.PanicV2(r, err.Error(), nil)

	// 常に500エラーを返す
	w.WriteHeader(http.StatusInternalServerError)
}

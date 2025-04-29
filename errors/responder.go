package errors

import (
	"github.com/go-openapi/runtime/middleware"
)

// NewErrorResponder はエラーをResponderに変換するヘルパー関数です
func NewErrorResponder(err error) middleware.Responder {
	// エラーを直接返す
	// グローバルエラーハンドラがこのエラーを処理します
	return middleware.Error(500, err)
}

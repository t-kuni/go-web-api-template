package errors

import (
	"github.com/go-openapi/errors"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/system"
	"net/http"
)

// NewCustomServeError はカスタムエラーハンドラを生成する関数です
func NewCustomServeError(logger system.ILogger) func(w http.ResponseWriter, r *http.Request, err error) {
	return func(rw http.ResponseWriter, r *http.Request, err error) {
		switch err.(type) {
		case *errors.CompositeError,
			*errors.MethodNotAllowedError,
			errors.Error:
			// go-openapi管轄のエラーは標準エラーハンドラであるerrors.ServeErrorに移譲する
			if logger != nil {
				logger.WarnWithError(r, err, nil)
			}
			errors.ServeError(rw, r, err)
		case nil:
			if logger != nil {
				logger.Panic(r, "Unknown Error", nil)
			}
			rw.WriteHeader(http.StatusInternalServerError)
		default:
			if logger != nil {
				logger.PanicV2(r, err.Error(), nil)
			}
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}
}

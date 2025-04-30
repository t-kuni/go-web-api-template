package middleware

import (
	"fmt"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/system"
	"net/http"
	"runtime"
)

const (
	StackSize         = 4 << 10 // 4 KB
	DisableStackAll   = false
	DisablePrintStack = false
)

// Recover PanicをRecoverするミドルウェアです
type Recover struct {
	logger system.ILogger
}

func NewRecover(logger system.ILogger) (*Recover, error) {
	return &Recover{
		logger: logger,
	}, nil
}

// Recover returns a middleware which recovers from panics anywhere in the chain
// and handles the control to the centralized HTTPErrorHandler.
func (m Recover) Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				m.WritePanicLog(r, err)

				//shouldResponse := !c.Response().Committed
				//if shouldResponse {
				//	err := c.NoContent(http.StatusInternalServerError)
				//	if err != nil {
				//		logger.Error(c, err, nil)
				//	}
				//}
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (m Recover) WritePanicLog(r *http.Request, panicErr interface{}) {
	err, ok := panicErr.(error)
	if !ok {
		err = fmt.Errorf("%v", panicErr)
	}
	stack := make([]byte, StackSize)
	length := runtime.Stack(stack, !DisableStackAll)
	if !DisablePrintStack {
		msg := fmt.Sprintf("%v %s\n", err, stack[:length])
		m.logger.PanicV2(r, msg, nil)
	}
}

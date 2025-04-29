package system_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/t-kuni/go-web-api-template/infrastructure/system"
)

var loggerHook *test.Hook

func TestMain(m *testing.M) {
	godotenv.Load(filepath.Join("..", "..", ".env.feature"))

	// setupLoggerはprivateになったので直接呼び出せない
	// 代わりにNewLoggerを使用してロガーを初期化する
	_ = system.NewLogger()
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetOutput(io.Discard)

	loggerHook = test.NewGlobal()

	code := m.Run()
	os.Exit(code)
}

type TestCaseContainer struct {
	t *testing.T
}

func beforeEach(t *testing.T) *TestCaseContainer {
	loggerHook.Reset()

	return &TestCaseContainer{
		t: t,
	}
}

func afterEach(cont *TestCaseContainer) {
}

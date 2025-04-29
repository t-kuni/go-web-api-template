package system_test

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/t-kuni/go-web-api-template/infrastructure/system"
	"io"
	"os"
	"path/filepath"
	"testing"
)

var loggerHook *test.Hook

func TestMain(m *testing.M) {
	godotenv.Load(filepath.Join("..", "..", ".env.feature"))

	err := system.SetupLogger()
	if err != nil {
		panic(err)
	}
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

package testUtil

import (
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/t-kuni/go-web-api-template/di"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/system"
	"github.com/t-kuni/go-web-api-template/ent"
	dbImpl "github.com/t-kuni/go-web-api-template/infrastructure/db"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"
)

func TestMain(m *testing.M) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("Could not get current file path")
	}
	directory := filepath.Dir(file)
	godotenv.Load(filepath.Join(directory, "..", ".env.testing"))

	dbImpl.RegisterTxdbDriver()

	code := m.Run()
	os.Exit(code)
}

type TestCaseContainer struct {
	t           *testing.T
	MockCtrl    *gomock.Controller
	LoggerHook  *test.Hook
	FxOptions   []fx.Option
	FxContainer *fx.App
}

func Prepare(t *testing.T) *TestCaseContainer {
	cont := TestCaseContainer{}

	cont.t = t

	cont.MockCtrl = gomock.NewController(t)

	cont.FxOptions = []fx.Option{
		//fx.WithLogger(func() fxevent.Logger {
		//	// テスト時のfxに関するログは異なるロガーインスタンスから出力させる
		//	// （テスト対象の処理から出力されたログと分離するため）
		//	l, err := logger.NewLogger()
		//	if err != nil {
		//		t.Fatal(err)
		//	}
		//	return l
		//}),
		//fx.Invoke(func(log *logger.Logger) {
		//	cont.LoggerHook = log.SetupForTest()
		//}),
		fx.Decorate(dbImpl.NewTestConnector),
	}

	return &cont
}

// Finish テストケースの後処理
func (c *TestCaseContainer) Finish() {
	c.MockCtrl.Finish()
	err := c.FxContainer.Stop(c.t.Context())
	assert.NoError(c.t, err)
}

func (c *TestCaseContainer) Invoke(closure any) {
	c.FxOptions = append(c.FxOptions, fx.Invoke(closure))
}

func (c *TestCaseContainer) Decorate(closure any) {
	c.FxOptions = append(c.FxOptions, fx.Decorate(closure))
}

func (c *TestCaseContainer) Exec(closure any) {
	opts := append(c.FxOptions, fx.Invoke(closure))
	c.FxContainer = di.NewApp(opts...)
	err := c.FxContainer.Start(c.t.Context())
	assert.NoError(c.t, err)
}

// PrepareTestData 外部キー制約のチェックを無効化した状態で第二引数の処理を実行します
func (c *TestCaseContainer) PrepareTestData(closure func(entClient *ent.Client)) {
	c.Invoke(func(conn db.IConnector) {
		MustExec(conn.GetDB(), "SET FOREIGN_KEY_CHECKS = 0")
		closure(conn.GetEnt())
		MustExec(conn.GetDB(), "SET FOREIGN_KEY_CHECKS = 1")
	})
}

func (c *TestCaseContainer) SetTime(timeStr string) {
	mock := system.NewMockITimer(c.MockCtrl)
	mock.EXPECT().Now().Return(MustNewDateTime(timeStr)).AnyTimes()
	Override[system.ITimer](c, mock)
}

func Override[T any](cont *TestCaseContainer, impl any) {
	decorator := fx.Decorate(func() T { return impl.(T) })
	cont.FxOptions = append(cont.FxOptions, decorator)
}

// MustExec SQLを実行し、エラーが発生した場合はpanicを発生させます
func MustExec(db *sql.DB, sql string) sql.Result {
	result, err := db.Exec(sql)
	if err != nil {
		panic(err)
	}
	return result
}

// MustNewDateTime 日時を生成します
func MustNewDateTime(dateTime string) time.Time {
	now, err := time.Parse(time.RFC3339, dateTime)
	if err != nil {
		panic(err)
	}
	return now
}

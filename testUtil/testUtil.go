package testUtil

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/t-kuni/go-web-api-template/di"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/system"
	"github.com/t-kuni/go-web-api-template/ent"
	"github.com/t-kuni/go-web-api-template/infrastructure/db"
	"go.uber.org/fx"
	"go.uber.org/mock/gomock"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("Could not get current file path")
	}
	directory := filepath.Dir(file)
	godotenv.Load(filepath.Join(directory, "..", ".env.feature"))

	db.RegisterTxdbDriver()

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
		fx.Decorate(db.NewTestConnector),
	}

	return &cont
}

// Finish テストケースの後処理
func (c *TestCaseContainer) Finish() {
	c.MockCtrl.Finish()
	err := c.FxContainer.Stop(c.t.Context())
	if err != nil {
		c.t.Fatal(err)
	}
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
	if err != nil {
		c.t.Fatal(err)
	}
}

// PrepareTestData 外部キー制約のチェックを無効化した状態で第二引数の処理を実行します
func (c *TestCaseContainer) PrepareTestData(closure func(entClient *ent.Client)) {
	c.Invoke(func(conn db.Connector) {
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

// MustInsert データを挿入し、エラーが発生した場合はpanicを発生させます
func MustInsert(db *sql.DB, table string, records []map[string]interface{}) {
	if len(records) == 0 {
		return
	}

	// カラム名とプレースホルダーを生成
	columns := make([]string, 0, len(records[0]))
	placeholders := make([]string, 0, len(records[0]))
	for column := range records[0] {
		columns = append(columns, column)
		placeholders = append(placeholders, "?")
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)

	// 各レコードを挿入
	for _, record := range records {
		values := make([]interface{}, 0, len(record))
		for _, column := range columns {
			values = append(values, record[column])
		}

		_, err := db.Exec(query, values...)
		if err != nil {
			panic(err)
		}
	}
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

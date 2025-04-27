//go:build feature

package handler_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/t-kuni/go-web-api-template/application/handler"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/api"
	"github.com/t-kuni/go-web-api-template/domain/service"
	"github.com/t-kuni/go-web-api-template/ent"
	"github.com/t-kuni/go-web-api-template/testUtil"
	"go.uber.org/mock/gomock"
)

// infrastructureをモック化するパターン
func TestHello(t *testing.T) {
	//
	// Prepare
	//
	cont := testUtil.Prepare(t)
	defer cont.Finish()

	cont.SetTime("2020-04-10T00:00:00+09:00")

	{
		mock := api.NewMockIBinanceApi(cont.MockCtrl)
		mock.
			EXPECT().
			GetExchangeInfo(gomock.Eq("BNB")).
			Return(&api.GetExchangeInfoResult{
				Symbols: []api.GetExchangeInfoResultSymbol{
					{
						Status: "DUMMY",
					},
				},
			}, nil)
		testUtil.Override[api.IBinanceApi](cont, mock)
	}

	//
	// Execute
	//
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	cont.Exec(func(testee *handler.HelloHandler) {
		err := testee.Hello(c)

		//
		// Assert
		//
		assert.NoError(t, err)
		assert.Contains(t, rec.Body.String(), "DUMMY")
	})
}

// serviceをモック化するパターン
func TestHello2(t *testing.T) {
	//
	// Prepare
	//
	cont := testUtil.Prepare(t)
	defer cont.Finish()

	cont.SetTime("2020-04-10T00:00:00+09:00")

	{
		mock := service.NewMockIExampleService(cont.MockCtrl)
		createdAt, err := time.Parse("2006-01-02 15:04:05 MST", "2014-12-31 12:31:24 JST")
		if err != nil {
			return
		}
		mock.
			EXPECT().
			Exec(gomock.Any(), gomock.Eq("BNB")).
			Return("DUMMY", []*ent.Company{
				{
					ID:        1,
					Name:      "TEST",
					CreatedAt: createdAt,
					Edges:     ent.CompanyEdges{},
				},
			}, nil)
		testUtil.Override[service.IExampleService](cont, mock)
	}

	//
	// Execute
	//
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	cont.Exec(func(testee *handler.HelloHandler) {
		err := testee.Hello(c)

		//
		// Assert
		//
		assert.NoError(t, err)
		assert.Contains(t, rec.Body.String(), "DUMMY")
	})
}

// DBにテストデータを挿入するパターン
func TestHello3(t *testing.T) {
	//
	// Prepare
	//
	cont := testUtil.Prepare(t)
	defer cont.Finish()

	cont.SetTime("2020-04-10T00:00:00+09:00")

	cont.PrepareTestData(func(db *ent.Client) {
		db.User.Create().SetID(1).SetAge(10).SetName("テストユーザ").SaveX(t.Context())
		db.Company.Create().SetID(1).SetName("テスト企業").SaveX(t.Context())
	})

	{
		mock := api.NewMockIBinanceApi(cont.MockCtrl)
		mock.
			EXPECT().
			GetExchangeInfo(gomock.Eq("BNB")).
			Return(&api.GetExchangeInfoResult{
				Symbols: []api.GetExchangeInfoResultSymbol{
					{
						Status: "DUMMY",
					},
				},
			}, nil)
		testUtil.Override[api.IBinanceApi](cont, mock)
	}

	//
	// Execute
	//
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	cont.Exec(func(testee *handler.HelloHandler) {
		err := testee.Hello(c)

		//
		// Assert
		//
		assert.NoError(t, err)

		var res handler.HelloResponse
		buf, err := io.ReadAll(rec.Body)
		assert.NoError(t, err)
		err = json.Unmarshal(buf, &res)

		assert.NoError(t, err)
		assert.Equal(t, "DUMMY", res.Status)
		assert.Len(t, res.Companies, 1)
		assert.Equal(t, "テスト企業", res.Companies[0].Name)
		assert.Len(t, res.Companies[0].Users, 0)
	})
}

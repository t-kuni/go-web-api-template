//go:build feature

package handler_test

import (
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/t-kuni/go-web-api-skeleton/domain/infrastructure/api"
	"github.com/t-kuni/go-web-api-skeleton/domain/service"
	"github.com/t-kuni/go-web-api-skeleton/ent"
	"github.com/t-kuni/go-web-api-skeleton/interface/handler"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// infrastructureをモック化するパターン
func TestHello(t *testing.T) {
	//
	// Prepare
	//
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	binanceApiMock := api.NewMockBinanceApiInterface(ctrl)

	binanceApiMock.
		EXPECT().
		GetExchangeInfo(gomock.Eq("BNB")).
		Return(&api.GetExchangeInfoResult{
			Symbols: []api.GetExchangeInfoResultSymbol{
				{
					Status: "DUMMY",
				},
			},
		}, nil)

	//
	// Execute
	//
	h := handler.ProvideHello(service.ProvideExampleService(binanceApiMock, dbConnector))
	err := h.Hello(c)

	//
	// Assert
	//
	assert.NoError(t, err)
	assert.Contains(t, rec.Body.String(), "DUMMY")
}

// serviceをモック化するパターン
func TestHello2(t *testing.T) {
	//
	// Prepare
	//
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	exampleServiceMock := service.NewMockExampleServiceInterface(ctrl)

	createdAt, err := time.Parse("2006-01-02 15:04:05 MST", "2014-12-31 12:31:24 JST")
	if err != nil {
		return
	}
	exampleServiceMock.
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

	//
	// Execute
	//
	h := handler.ProvideHello(exampleServiceMock)
	err = h.Hello(c)

	//
	// Assert
	//
	assert.NoError(t, err)
	assert.Contains(t, rec.Body.String(), "DUMMY")
}

// DBにテストデータを挿入するパターン
func TestHello3(t *testing.T) {
	//
	// Prepare
	//
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	binanceApiMock := api.NewMockBinanceApiInterface(ctrl)
	binanceApiMock.
		EXPECT().
		GetExchangeInfo(gomock.Eq("BNB")).
		Return(&api.GetExchangeInfoResult{
			Symbols: []api.GetExchangeInfoResultSymbol{
				{
					Status: "DUMMY",
				},
			},
		}, nil)

	db := dbConnector.GetDB()

	_, err := db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	assert.NoError(t, err)
	_, err = db.Exec("START TRANSACTION")
	assert.NoError(t, err)
	t.Cleanup(func() { db.Exec("ROLLBACK") })

	_, err = db.Exec("INSERT INTO users(id, age, name, created_at) VALUES (1, 10, 'テストユーザ', '2010-12-31 23:59:59')")
	assert.NoError(t, err)
	_, err = db.Exec("INSERT INTO companies(id, name, created_at) VALUES (1, 'テスト企業', '2010-12-31 23:59:59')")
	assert.NoError(t, err)

	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	assert.NoError(t, err)

	//
	// Execute
	//
	h := handler.ProvideHello(service.ProvideExampleService(binanceApiMock, dbConnector))
	err = h.Hello(c)

	//
	// Assert
	//
	assert.NoError(t, err)

	var res handler.HelloResponse
	buf, err := ioutil.ReadAll(rec.Body)
	assert.NoError(t, err)
	err = json.Unmarshal(buf, &res)

	assert.NoError(t, err)
	assert.Equal(t, "DUMMY", res.Status)
	assert.Len(t, res.Companies, 1)
	assert.Equal(t, "テスト企業", res.Companies[0].Name)
	assert.Len(t, res.Companies[0].Users, 0)
}

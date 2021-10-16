//go:build feature

package handler_test

import (
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/t-kuni/go-web-api-skeleton/domain/infrastructure/api"
	"github.com/t-kuni/go-web-api-skeleton/interface/handler"
	"github.com/t-kuni/go-web-api-skeleton/wire"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHello(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	app := wire.InitializeApp()

	binanceApiMock := api.NewMockBinanceApi(ctrl)

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

	app.ExampleService.BinanceApi = binanceApiMock

	err := handler.Hello(app)(c)

	assert.NoError(t, err)
	assert.Contains(t, rec.Body.String(), "DUMMY")
}

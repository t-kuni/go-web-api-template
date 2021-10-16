//go:build feature

package handler_test

import (
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/t-kuni/go-web-api-skeleton/domain/infrastructure/api"
	"github.com/t-kuni/go-web-api-skeleton/domain/service"
	"github.com/t-kuni/go-web-api-skeleton/interface/handler"
	"net/http"
	"net/http/httptest"
	"testing"
)

// infrastructureをモック化するパターン
func TestHello(t *testing.T) {
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

	h := handler.ProvideHello(service.ProvideExampleService(binanceApiMock))
	err := h.Hello(c)

	assert.NoError(t, err)
	assert.Contains(t, rec.Body.String(), "DUMMY")
}

// serviceをモック化するパターン
func TestHello2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	exampleServiceMock := service.NewMockExampleServiceInterface(ctrl)

	exampleServiceMock.
		EXPECT().
		Exec(gomock.Eq("BNB")).
		Return("DUMMY", nil)

	h := handler.ProvideHello(exampleServiceMock)
	err := h.Hello(c)

	assert.NoError(t, err)
	assert.Contains(t, rec.Body.String(), "DUMMY")
}

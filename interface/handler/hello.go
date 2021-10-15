package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/t-kuni/go-cli-app-skeleton/wire"
	"net/http"
)

func Hello(app wire.App) func(c echo.Context) error {
	return func(c echo.Context) error {
		status, err := app.ExampleService.Exec("BNB")
		if err != nil {
			return err
		}
		return c.String(http.StatusOK, "Hello, World! Status:"+status)
	}
}

package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/t-kuni/go-web-api-skeleton/domain/service"
	"net/http"
)

type HelloHandler struct {
	ExampleService service.ExampleServiceInterface
}

func ProvideHello(exampleService service.ExampleServiceInterface) HelloHandler {
	return HelloHandler{exampleService}
}

func (h HelloHandler) Hello(c echo.Context) error {
	status, err := h.ExampleService.Exec("BNB")
	if err != nil {
		return err
	}
	return c.String(http.StatusOK, "Hello, World! Status:"+status)
}

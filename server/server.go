package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"github.com/t-kuni/go-web-api-template/errors/handler"
	"os"
)

type Server struct {
	Echo *echo.Echo
}

func NewServer(i *do.Injector) (*Server, error) {
	e := echo.New()

	e.HTTPErrorHandler = do.MustInvoke[*handler.ErrorHandler](i).Handler
	e.Validator = do.MustInvoke[echo.Validator](i)

	//e.Use(xxx)

	return &Server{e}, nil
}

func (s Server) Start() error {
	port := os.Getenv("SERVER_PORT")
	return s.Echo.Start(fmt.Sprintf(":%s", port))
}

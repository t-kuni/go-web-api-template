package router

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"github.com/t-kuni/go-web-api-template/application/handler"
)

type Router struct {
	HelloHandler    *handler.HelloHandler
	PostUserHandler *handler.PostUserHandler
}

func NewRouter(i *do.Injector) (*Router, error) {
	return &Router{
		HelloHandler:    do.MustInvoke[*handler.HelloHandler](i),
		PostUserHandler: do.MustInvoke[*handler.PostUserHandler](i),
	}, nil
}

func (r Router) Attach(e *echo.Echo) {
	g := e.Group("")

	g.GET("/", r.HelloHandler.Hello)
	g.POST("/users", r.PostUserHandler.PostUser)
}

func group(g *echo.Group, prefix string, m []echo.MiddlewareFunc, cb func(*echo.Group)) {
	childGroup := g.Group(prefix, m...)
	cb(childGroup)
}

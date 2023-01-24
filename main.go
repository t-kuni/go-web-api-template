package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4/middleware"
	"github.com/samber/do"
	"github.com/t-kuni/go-web-api-template/di"
	"github.com/t-kuni/go-web-api-template/interface/handler"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	godotenv.Load()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	container := di.NewContainer()
	defer container.Shutdown()

	// Routes
	e.GET("/", do.MustInvoke[*handler.HelloHandler](container).Hello)
	e.POST("/users", do.MustInvoke[*handler.PostUserHandler](container).PostUser)

	// Start server
	port := os.Getenv("SERVER_PORT")
	e.Logger.Fatal(e.Start(":" + port))
}

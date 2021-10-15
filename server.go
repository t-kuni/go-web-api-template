package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4/middleware"
	"github.com/t-kuni/go-cli-app-skeleton/interface/handler"
	"github.com/t-kuni/go-cli-app-skeleton/wire"
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

	app := wire.InitializeApp()

	// Routes
	e.GET("/", handler.Hello(app))

	// Start server
	port := os.Getenv("SERVER_PORT")
	e.Logger.Fatal(e.Start(":" + port))
}

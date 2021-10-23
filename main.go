package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4/middleware"
	"github.com/t-kuni/go-web-api-skeleton/wire"
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

	app, cleanup, err := wire.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// Routes
	e.GET("/", app.HelloHandler.Hello)
	e.POST("/users", app.PostUserHandler.PostUser)

	// Start server
	port := os.Getenv("SERVER_PORT")
	e.Logger.Fatal(e.Start(":" + port))
}

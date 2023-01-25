package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/samber/do"
	"github.com/t-kuni/go-web-api-template/di"
	routerPackage "github.com/t-kuni/go-web-api-template/router"
	serverPackage "github.com/t-kuni/go-web-api-template/server"
	"os"
)

func main() {
	godotenv.Load()

	app := di.NewApp()
	defer app.Shutdown()

	server := do.MustInvoke[*serverPackage.Server](app)
	router := do.MustInvoke[*routerPackage.Router](app)

	router.Attach(server.Echo)
	err := server.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

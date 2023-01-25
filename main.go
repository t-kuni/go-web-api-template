package main

import (
	"github.com/joho/godotenv"
	"github.com/samber/do"
	"github.com/t-kuni/go-web-api-template/di"
	"github.com/t-kuni/go-web-api-template/logger"
	routerPackage "github.com/t-kuni/go-web-api-template/router"
	serverPackage "github.com/t-kuni/go-web-api-template/server"
	"log"
	"os"
)

func main() {
	godotenv.Load()

	if err := logger.SetupLogger(); err != nil {
		log.Fatalf("Logger initialization failed: %+v", err)
		os.Exit(1)
	}

	app := di.NewApp()
	defer app.Shutdown()

	server := do.MustInvoke[*serverPackage.Server](app)
	router := do.MustInvoke[*routerPackage.Router](app)

	router.Attach(server.Echo)
	err := server.Start()
	if err != nil {
		logger.SimpleFatal(err, nil)
		os.Exit(1)
	}
}

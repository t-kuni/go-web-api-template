package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/samber/do"
	"github.com/t-kuni/go-web-api-template/di"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	"os"
	"path/filepath"
)

func main() {
	godotenv.Load(filepath.Join(".env"))

	var (
		database = flag.String("database", os.Getenv("DB_DATABASE"), "database name")
	)
	flag.Parse()
	err := os.Setenv("DB_DATABASE", *database)
	if err != nil {
		panic(err)
	}

	fmt.Println("Target Database: " + os.Getenv("DB_DATABASE"))

	app := di.NewApp()
	defer app.Shutdown()

	dbConnector := do.MustInvoke[db.Connector](app)

	if err := dbConnector.Migrate(context.Background()); err != nil {
		panic(err)
	}

	fmt.Println("Migrate successfully!")
}

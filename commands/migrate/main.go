package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/t-kuni/go-web-api-template/wire"
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

	app, cleanup, err := wire.InitializeApp()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.DBConnector.Migrate(context.Background()); err != nil {
		panic(err)
	}

	fmt.Println("Migrate successfully!")
}

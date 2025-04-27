package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/t-kuni/go-web-api-template/di"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	"go.uber.org/fx"
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

	ctx := context.Background()
	app := di.NewApp(fx.Invoke(func(conn db.IConnector) {
		if err := conn.Migrate(context.Background()); err != nil {
			panic(err)
		}
	}))
	defer app.Stop(ctx)

	err = app.Start(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Migrate successfully!")
}

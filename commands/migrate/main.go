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
	"runtime"
)

func main() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("Could not get current file path")
	}
	directory := filepath.Dir(file)
	godotenv.Load(filepath.Join(directory, "..", "..", ".env"))

	var (
		reset = flag.Bool("reset", false, "reset database before migration")
	)
	flag.Parse()

	database := os.Getenv("DB_DATABASE")
	fmt.Println("Target Database: " + database)

	ctx := context.Background()
	app := di.NewApp(fx.Invoke(func(conn db.IConnector) {
		if *reset {
			fmt.Println("Resetting database...")

			// データベースを削除して再作成するためにはinformation_schemaに接続する必要がある
			db := conn.GetDB()

			// データベースを削除
			_, err := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS `%s`", database))
			if err != nil {
				panic(fmt.Errorf("failed to drop database: %w", err))
			}

			// データベースを作成
			_, err = db.Exec(fmt.Sprintf("CREATE DATABASE `%s`", database))
			if err != nil {
				panic(fmt.Errorf("failed to create database: %w", err))
			}

			_, err = db.Exec(fmt.Sprintf("USE `%s`", database))
			if err != nil {
				panic(fmt.Errorf("failed to use database: %w", err))
			}

			fmt.Println("Database reset successfully!")
		}

		if err := conn.Migrate(ctx); err != nil {
			panic(err)
		}
	}))
	defer app.Stop(ctx)

	err := app.Start(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("Migrate successfully!")
}

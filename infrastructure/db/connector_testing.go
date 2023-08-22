//go:build feature

package db

import (
	"database/sql"
	sql2 "entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/DATA-DOG/go-txdb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/samber/do"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	"github.com/t-kuni/go-web-api-template/ent"
	"os"
)

func NewTestConnector(i *do.Injector) (db.Connector, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, database)

	txdb.Register("txdb", "mysql", dsn)

	db, err := sql.Open("txdb", "identifier")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	drv := sql2.OpenDB("mysql", db)
	client := ent.NewClient(ent.Driver(drv))
	return &Connector{DB: db, Client: client}, nil
}

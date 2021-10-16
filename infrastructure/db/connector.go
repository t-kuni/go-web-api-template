package db

import (
	"context"
	"database/sql"
	sql2 "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/t-kuni/go-web-api-skeleton/ent"
	"os"
)

type Connector struct {
	DB     *sql.DB
	Client *ent.Client
}

func ProvideConnector() (*Connector, func(), error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		if err := db.Close(); err != nil {
			// TODO Logging
		}
	}

	drv := sql2.OpenDB("mysql", db)
	return &Connector{DB: db, Client: ent.NewClient(ent.Driver(drv))}, cleanup, nil
}

func (c Connector) GetDB() *sql.DB {
	return c.DB
}

func (c Connector) GetEnt() *ent.Client {
	return c.Client
}

func (c Connector) Migrate(ctx context.Context, opts ...schema.MigrateOption) error {
	return c.Client.Schema.Create(ctx, opts...)
}

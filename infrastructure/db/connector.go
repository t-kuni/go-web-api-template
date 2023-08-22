package db

import (
	"context"
	"database/sql"
	sql2 "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/samber/do"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	"github.com/t-kuni/go-web-api-template/ent"
	"os"
)

type Connector struct {
	DB     *sql.DB
	Client *ent.Client
}

func NewConnector(i *do.Injector) (db.Connector, error) {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, database)
	db, err := sql.Open("mysql", dsn)
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

func (c Connector) GetDB() *sql.DB {
	return c.DB
}

func (c Connector) GetEnt() *ent.Client {
	return c.Client
}

func (c Connector) Migrate(ctx context.Context, opts ...schema.MigrateOption) error {
	return c.Client.Schema.Create(ctx, opts...)
}

func (c Connector) Transaction(ctx context.Context, fn func(tx *ent.Client) error) error {
	tx, err := c.Client.Tx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if v := recover(); v != nil {
			tx.Rollback()
			panic(v)
		}
	}()
	txClient := tx.Client()
	if err := fn(txClient); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}
	return nil
}

func (c Connector) Shutdown() error {
	return c.DB.Close()
}

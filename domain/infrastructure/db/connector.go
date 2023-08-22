package db

import (
	"context"
	"database/sql"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/t-kuni/go-web-api-template/ent"
)

type Connector interface {
	GetDB() *sql.DB
	GetEnt() *ent.Client
	Transaction(ctx context.Context, fn func(tx *ent.Client) error) error
	Migrate(ctx context.Context, opts ...schema.MigrateOption) error
}

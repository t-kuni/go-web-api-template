package db

import (
	"context"
	"database/sql"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/t-kuni/go-web-api-skeleton/ent"
)

type ConnectorInterface interface {
	GetDB() *sql.DB
	GetEnt() *ent.Client
	Migrate(ctx context.Context, opts ...schema.MigrateOption) error
}

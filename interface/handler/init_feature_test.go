//go:build feature

package handler_test

import (
	"github.com/joho/godotenv"
	"github.com/t-kuni/go-web-api-template/infrastructure/db"
	"os"
	"path/filepath"
	"testing"
)

var dbConnector *db.Connector

func TestMain(m *testing.M) {
	godotenv.Load(filepath.Join("..", "..", ".env.feature"))

	code := m.Run()
	os.Exit(code)
}

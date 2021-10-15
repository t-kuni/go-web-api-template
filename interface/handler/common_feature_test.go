package handler

import (
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	godotenv.Load(filepath.Join("..", "..", ".env.feature"))
	code := m.Run()
	os.Exit(code)
}

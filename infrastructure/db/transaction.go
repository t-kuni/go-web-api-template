//go:generate mockgen -source=$GOFILE -destination=${GOFILE}_mock.go -package=$GOPACKAGE

package db

import (
	"database/sql"
)

type TransactionInterface interface {
	Begin() error
	Commit() error
	Rollback() error
}

type Transaction struct {
	DB *sql.DB
}

func NewTransaction(db *sql.DB) *Transaction {
	return &Transaction{DB: db}
}

func (c Transaction) Begin() error {
	_, err := c.DB.Exec("START TRANSACTION")
	return err
}

func (c Transaction) Commit() error {
	_, err := c.DB.Exec("COMMIT")
	return err
}

func (c Transaction) Rollback() error {
	_, err := c.DB.Exec("ROLLBACK")
	return err
}

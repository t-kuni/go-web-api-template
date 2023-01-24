//go:generate mockgen -source=$GOFILE -destination=${GOFILE}_mock.go -package=$GOPACKAGE

package db

import (
	"context"
	"github.com/t-kuni/go-web-api-template/ent"
)

type TransactionInterface interface {
	Begin(ctx context.Context) (*ent.Tx, error)
	Commit(tx *ent.Tx) error
	Rollback(tx *ent.Tx) error
}

type Transaction struct {
	Client *ent.Client
}

func NewTransaction(client *ent.Client) *Transaction {
	return &Transaction{client}
}

func (c Transaction) Begin(ctx context.Context) (*ent.Tx, error) {
	return c.Client.Tx(ctx)
}

func (c Transaction) Commit(tx *ent.Tx) error {
	return tx.Commit()
}

func (c Transaction) Rollback(tx *ent.Tx) error {
	return tx.Rollback()
}

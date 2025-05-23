//go:generate go tool mockgen -source=$GOFILE -destination=${GOFILE}_mock.go -package=$GOPACKAGE
package service

import (
	"context"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/api"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	"github.com/t-kuni/go-web-api-template/ent"
)

type ExampleService struct {
	BinanceApi  api.IBinanceApi
	DBConnector db.IConnector
}

type IExampleService interface {
	Exec(ctx context.Context, baseAsset string) (string, []*ent.Company, error)
}

func NewExampleService(conn db.IConnector, binanceApi api.IBinanceApi) (IExampleService, error) {
	return &ExampleService{
		binanceApi,
		conn,
	}, nil
}

func (s ExampleService) Exec(ctx context.Context, baseAsset string) (string, []*ent.Company, error) {
	info, err := s.BinanceApi.GetExchangeInfo(baseAsset)
	if err != nil {
		return "", nil, err
	}

	companies, err := s.DBConnector.GetEnt().Company.Query().
		WithUsers().
		All(ctx)
	if err != nil {
		return "", nil, err
	}

	return info.Symbols[0].Status, companies, nil
}

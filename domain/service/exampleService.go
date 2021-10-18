//go:generate mockgen -source=$GOFILE -destination=${GOFILE}_mock.go -package=$GOPACKAGE
package service

import (
	"context"
	"github.com/t-kuni/go-web-api-skeleton/domain/infrastructure/api"
	"github.com/t-kuni/go-web-api-skeleton/domain/infrastructure/db"
	"github.com/t-kuni/go-web-api-skeleton/ent"
)

type ExampleService struct {
	BinanceApi  api.BinanceApiInterface
	DBConnector db.ConnectorInterface
}

type ExampleServiceInterface interface {
	Exec(ctx context.Context, baseAsset string) (string, []*ent.Company, error)
}

func ProvideExampleService(
	binanceApi api.BinanceApiInterface,
	dbConnector db.ConnectorInterface,
) *ExampleService {
	return &ExampleService{binanceApi, dbConnector}
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

//go:generate mockgen -source=$GOFILE -destination=${GOFILE}_mock.go -package=$GOPACKAGE
package service

import (
	"context"
	"github.com/samber/do"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/api"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	"github.com/t-kuni/go-web-api-template/ent"
)

type ExampleService struct {
	BinanceApi  api.BinanceApiInterface
	DBConnector db.Connector
}

type ExampleServiceInterface interface {
	Exec(ctx context.Context, baseAsset string) (string, []*ent.Company, error)
}

func NewExampleService(i *do.Injector) (ExampleServiceInterface, error) {
	return &ExampleService{
		do.MustInvoke[api.BinanceApiInterface](i),
		do.MustInvoke[db.Connector](i),
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

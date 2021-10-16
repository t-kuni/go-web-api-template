//go:generate mockgen -source=$GOFILE -destination=${GOFILE}_mock.go -package=$GOPACKAGE
package service

import (
	"github.com/t-kuni/go-web-api-skeleton/domain/infrastructure/api"
)

type ExampleService struct {
	BinanceApi api.BinanceApiInterface
}

type ExampleServiceInterface interface {
	Exec(baseAsset string) (string, error)
}

func ProvideExampleService(binanceApi api.BinanceApiInterface) ExampleService {
	return ExampleService{binanceApi}
}

func (s ExampleService) Exec(baseAsset string) (string, error) {
	info, err := s.BinanceApi.GetExchangeInfo(baseAsset)
	if err != nil {
		return "", err
	}
	return info.Symbols[0].Status, nil
}

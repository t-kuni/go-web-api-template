//go:generate mockgen -source=$GOFILE -destination=${GOFILE}_mock.go -package=$GOPACKAGE
package api

type GetExchangeInfoResult struct {
	Symbols []GetExchangeInfoResultSymbol
}

type GetExchangeInfoResultSymbol struct {
	Symbol     string
	Status     string
	BaseAsset  string
	QuoteAsset string
}

type IBinanceApi interface {
	GetExchangeInfo(baseAsset string) (*GetExchangeInfoResult, error)
}

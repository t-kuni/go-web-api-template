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

type BinanceApi interface {
	GetExchangeInfo(baseAsset string) (*GetExchangeInfoResult, error)
}

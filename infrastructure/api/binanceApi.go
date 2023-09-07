package api

import (
	"encoding/json"
	"fmt"
	"github.com/samber/do"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/api"
	"io/ioutil"
	"net/http"
)

type BinanceApi struct {
}

func NewBinanceApi(i *do.Injector) (api.IBinanceApi, error) {
	return &BinanceApi{}, nil
}

func (b *BinanceApi) GetExchangeInfo(baseAsset string) (*api.GetExchangeInfoResult, error) {
	url := fmt.Sprintf("https://api.binance.com/api/v3/exchangeInfo?symbol=%sBTC", baseAsset)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result api.GetExchangeInfoResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

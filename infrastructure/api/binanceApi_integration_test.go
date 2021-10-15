//go:build integration

package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetExchangeInfo(t *testing.T) {
	var api BinanceApi
	actual, err := api.GetExchangeInfo("BNB")
	assert.NoError(t, err)

	assert.Equal(t, "BNB", actual.Symbols[0].BaseAsset)
}

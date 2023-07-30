package exchange

import (
	"testing"
)

func TestCoinbaseGetCurrencyRate(t *testing.T) {
	loadEnvironment(t)

	coinbase := CoinbaseFactory{}.CreateCoinbaseFactory()
	GetCurrencyRateTest(coinbase, t)
}

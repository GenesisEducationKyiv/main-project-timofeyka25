package exchange

import "testing"

func TestCoingeckoGetCurrencyRate(t *testing.T) {
	loadEnvironment(t)

	coingecko := CoingeckoFactory{}.CreateCoingeckoFactory()
	GetCurrencyRateTest(coingecko, t)
}

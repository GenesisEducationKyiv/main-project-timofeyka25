package exchange

import "testing"

func TestBinanceGetCurrencyRate(t *testing.T) {
	loadEnvironment(t)

	binance := BinanceFactory{}.CreateBinanceFactory()
	GetCurrencyRateTest(binance, t)
}

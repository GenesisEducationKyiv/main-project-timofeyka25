package exchange

import "testing"

func TestBTCTradeUAGetCurrencyRate(t *testing.T) {
	loadEnvironment(t)

	btcTradeUA := BTCTradeUAFactory{}.CreateBTCTradeUAFactory()
	GetCurrencyRateTest(btcTradeUA, t)
}

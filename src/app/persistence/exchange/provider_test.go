package exchange

import (
	"genesis-test/src/app/domain"
	"genesis-test/src/app/service"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func loadEnvironment(t *testing.T) {
	if err := godotenv.Load("../../../../.env"); err != nil {
		t.Fatal("Failed to load .env file")
	}
}

func GetCurrencyRateTest(chain service.ExchangeChain, t *testing.T) {
	pair := &domain.CurrencyPair{
		BaseCurrency:  "BTC",
		QuoteCurrency: "UAH",
	}

	rate, err := chain.GetCurrencyRate(pair)

	require.NoError(t, err)
	require.Equal(t, pair.GetBaseCurrency(), rate.GetBaseCurrency())
	require.Equal(t, pair.GetQuoteCurrency(), rate.GetQuoteCurrency())
	require.NotEmpty(t, rate.Price)
}

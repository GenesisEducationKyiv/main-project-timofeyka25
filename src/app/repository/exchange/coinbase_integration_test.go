package exchange

import (
	"genesis-test/src/app/domain"
	"genesis-test/src/config"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestExchangeRepository_GetCurrencyRate(t *testing.T) {
	if err := godotenv.Load("../../../../.env"); err != nil {
		t.Fatal("Failed to load .env file")
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coinbaseRepo := NewExchangeCoinbaseRepository(config.Get().CryptoAPIFormatURL)
	pair := &domain.CurrencyPair{
		BaseCurrency:  "BTC",
		QuoteCurrency: "UAH",
	}

	rate, err := coinbaseRepo.GetCurrencyRate(pair)
	require.NoError(t, err)
	require.Equal(t, pair.BaseCurrency, rate.BaseCurrency)
	require.Equal(t, pair.QuoteCurrency, rate.QuoteCurrency)
	require.NotEmpty(t, rate.Price)
}

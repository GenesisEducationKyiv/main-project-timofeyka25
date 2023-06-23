package repository

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestExchangeRepository_GetCurrencyRate(t *testing.T) {
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Fatal("Failed to load .env file")
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	exchangeRepo := NewExchangeRepository()
	baseCurrency := "BTC"
	quoteCurrency := "UAH"

	rate, err := exchangeRepo.GetCurrencyRate(baseCurrency, quoteCurrency)
	require.NoError(t, err)
	require.Equal(t, rate.BaseCurrency, baseCurrency)
	require.Equal(t, rate.QuoteCurrency, quoteCurrency)
	require.NotEmpty(t, rate.Price)
}

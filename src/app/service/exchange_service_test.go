package service

import (
	"genesis-test/src/app/domain"
	"genesis-test/src/app/domain/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestExchangeService_GetCurrencyRate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExchangeRepo := mocks.NewMockExchangeRepository(ctrl)

	BTCUAHPair := &domain.CurrencyPair{
		BaseCurrency:  "BTC",
		QuoteCurrency: "UAH",
	}

	excService := NewExchangeService(BTCUAHPair, mockExchangeRepo)

	mockResponse := &domain.CurrencyRate{
		Price:        123456,
		CurrencyPair: *BTCUAHPair,
	}

	mockExchangeRepo.EXPECT().GetCurrencyRate(BTCUAHPair).Return(mockResponse, nil)
	rate, err := excService.GetCurrencyRate()
	require.NoError(t, err)

	require.NoError(t, err)
	require.Equalf(t, mockResponse.Price, rate, "rates are not equal")
}

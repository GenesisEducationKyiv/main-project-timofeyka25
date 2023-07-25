package exchange

import (
	"genesis-test/src/app/application/mocks"
	"genesis-test/src/app/domain/model"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestExchangeService_GetCurrencyRate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProvider := mocks.NewMockExchangeProvider(ctrl)

	BTCUAHPair := &model.CurrencyPair{
		BaseCurrency:  "BTC",
		QuoteCurrency: "UAH",
	}

	excService := NewExchangeService(mockProvider)

	mockResponse := &model.CurrencyRate{
		Price:        123456,
		CurrencyPair: *BTCUAHPair,
	}

	mockProvider.EXPECT().GetCurrencyRate(BTCUAHPair).Return(mockResponse, nil)
	rate, err := excService.GetCurrencyRate(BTCUAHPair)
	require.NoError(t, err)

	require.NoError(t, err)
	require.Equalf(t, mockResponse.Price, rate, "rates are not equal")
}

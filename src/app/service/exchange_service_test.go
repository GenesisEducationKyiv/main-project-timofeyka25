package service

import (
	"genesis-test/src/app/domain"
	mocks "genesis-test/src/app/domain/mocks"
	"genesis-test/src/app/repository"
	"genesis-test/src/config"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestExchangeService_GetCurrencyRate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockExchangeRepo := mocks.NewMockExchangeRepository(ctrl)
	repos := &repository.Repositories{
		Exchange: mockExchangeRepo,
	}

	excService := NewExchangeService(repos)

	cfg := &config.Config{
		BaseCurrency:  "BTC",
		QuoteCurrency: "UAH",
	}

	mockResponse := &domain.CurrencyRate{
		Price:         "888888",
		BaseCurrency:  cfg.BaseCurrency,
		QuoteCurrency: cfg.QuoteCurrency,
	}

	mockExchangeRepo.EXPECT().GetCurrencyRate(cfg.BaseCurrency, cfg.QuoteCurrency).Return(mockResponse, nil)
	rate, err := excService.GetCurrencyRate(cfg)
	require.NoError(t, err)

	rateInt, err := strconv.Atoi(mockResponse.Price)
	require.NoError(t, err)
	require.Equal(t, rateInt, rate)
}

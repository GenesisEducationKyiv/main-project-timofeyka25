package service

import (
	"genesis-test/src/app/domain"
	"genesis-test/src/app/handler"

	"github.com/pkg/errors"
)

type exchangeService struct {
	pair         *domain.CurrencyPair
	exchangeRepo ExchangeRepository
}

func NewExchangeService(
	pair *domain.CurrencyPair,
	exchangeRepo ExchangeRepository,
) handler.ExchangeService {
	return &exchangeService{
		pair:         pair,
		exchangeRepo: exchangeRepo,
	}
}

func (c *exchangeService) GetCurrencyRate() (float64, error) {
	rate, err := c.exchangeRepo.GetCurrencyRate(c.pair)
	if err != nil {
		return 0, errors.Wrap(err, "get rate")
	}

	return rate.Price, nil
}

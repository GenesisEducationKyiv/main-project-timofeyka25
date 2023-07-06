package exchange

import (
	"genesis-test/src/app/application"
	"genesis-test/src/app/domain"
	"genesis-test/src/app/domain/model"

	"github.com/pkg/errors"
)

type exchangeService struct {
	pair             *model.CurrencyPair
	exchangeProvider application.ExchangeChain
}

func NewExchangeService(
	pair *model.CurrencyPair,
	exchangeProvider application.ExchangeChain,
) domain.ExchangeService {
	return &exchangeService{
		pair:             pair,
		exchangeProvider: exchangeProvider,
	}
}

func (c *exchangeService) GetCurrencyRate() (float64, error) {
	rate, err := c.exchangeProvider.GetCurrencyRate(c.pair)
	if err != nil {
		return 0, errors.Wrap(err, "get rate")
	}

	return rate.Price, nil
}

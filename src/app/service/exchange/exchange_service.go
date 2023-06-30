package exchange

import (
	"genesis-test/src/app/domain"
	"genesis-test/src/app/handler"
	"genesis-test/src/app/service"

	"github.com/pkg/errors"
)

type exchangeService struct {
	pair             *domain.CurrencyPair
	exchangeProvider service.ExchangeChain
}

func NewExchangeService(
	pair *domain.CurrencyPair,
	exchangeProvider service.ExchangeChain,
) handler.ExchangeService {
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

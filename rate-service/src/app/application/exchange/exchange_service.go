package exchange

import (
	"genesis-test/src/app/application"
	"genesis-test/src/app/domain"
	"genesis-test/src/app/domain/model"

	"github.com/pkg/errors"
)

type exchangeService struct {
	exchangeProvider application.ExchangeProvider
}

func NewExchangeService(
	exchangeProvider application.ExchangeProvider,
) domain.ExchangeService {
	return &exchangeService{
		exchangeProvider: exchangeProvider,
	}
}

func (c *exchangeService) GetCurrencyRate(pair *model.CurrencyPair) (float64, error) {
	rate, err := c.exchangeProvider.GetCurrencyRate(pair)
	if err != nil {
		return 0, errors.Wrap(err, "get rate")
	}

	return rate.Price, nil
}

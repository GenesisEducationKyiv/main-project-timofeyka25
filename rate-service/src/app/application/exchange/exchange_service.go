package exchange

import (
	"genesis-test/src/app/domain/model"

	"github.com/pkg/errors"
)

type ExchangeProvider interface {
	GetCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error)
}

type ExchangeService struct {
	exchangeProvider ExchangeProvider
}

func NewExchangeService(
	exchangeProvider ExchangeProvider,
) *ExchangeService {
	return &ExchangeService{
		exchangeProvider: exchangeProvider,
	}
}

func (c *ExchangeService) GetCurrencyRate(pair *model.CurrencyPair) (float64, error) {
	rate, err := c.exchangeProvider.GetCurrencyRate(pair)
	if err != nil {
		return 0, errors.Wrap(err, "get rate")
	}

	return rate.Price, nil
}

package exchange

import (
	"genesis-test/src/app/domain/model"
	"genesis-test/src/app/utils"
	"genesis-test/src/config"
)

type CoingeckoFactory struct{}

type CoingeckoProvider struct {
	coingeckoURL string
}

func (f CoingeckoFactory) CreateCoingeckoFactory() *CoingeckoProvider {
	return &CoingeckoProvider{
		coingeckoURL: config.Get().CoingeckoURL,
	}
}

func (c *CoingeckoProvider) GetCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error) {
	rate, err := c.getCurrencyRate(pair)
	if err != nil {
		return nil, err
	}

	return rate, nil
}

func (c *CoingeckoProvider) getCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error) {
	resp, err := c.doRequest(pair)
	if err != nil {
		return nil, err
	}
	return resp.toDefaultRate()
}

func (c *CoingeckoProvider) doRequest(pair *model.CurrencyPair) (*coingeckoResponse, error) {
	rate := new(coingeckoResponse)
	err := utils.GetAndParse(c.coingeckoURL, &rate)
	if err != nil {
		return nil, err
	}
	rate.BaseCurrency = pair.BaseCurrency.ToString()
	rate.QuoteCurrency = pair.QuoteCurrency.ToString()

	return rate, nil
}

type coingeckoResponse struct {
	Rates struct {
		QuoteCurrency struct {
			Price float64 `json:"value"`
		} `json:"uah"`
	} `json:"rates"`
	BaseCurrency  string
	QuoteCurrency string
}

func (c *coingeckoResponse) toDefaultRate() (*model.CurrencyRate, error) {
	return &model.CurrencyRate{
		Price: c.Rates.QuoteCurrency.Price,
		CurrencyPair: model.CurrencyPair{
			BaseCurrency:  model.CurrencyFactory{}.CreateCurrency(c.BaseCurrency),
			QuoteCurrency: model.CurrencyFactory{}.CreateCurrency(c.QuoteCurrency),
		},
	}, nil
}

package exchange

import (
	"genesis-test/src/app/application"
	"genesis-test/src/app/domain/model"
	"genesis-test/src/app/utils"
	"genesis-test/src/config"
)

type CoingeckoFactory struct{}

func (f CoingeckoFactory) CreateCoingeckoFactory() application.ExchangeChain {
	return &coingeckoProvider{
		coingeckoURL: config.Get().CoingeckoURL,
	}
}

type coingeckoProvider struct {
	coingeckoURL string
	next         application.ExchangeChain
}

func (c *coingeckoProvider) GetCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error) {
	rate, err := c.getCurrencyRate(pair)
	if err != nil && c.next != nil {
		return c.next.GetCurrencyRate(pair)
	}

	return rate, nil
}

func (c *coingeckoProvider) SetNext(chain application.ExchangeChain) {
	c.next = chain
}

func (c *coingeckoProvider) getCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error) {
	resp, err := c.doRequest(pair)
	if err != nil {
		return nil, err
	}
	return resp.toDefaultRate()
}

func (c *coingeckoProvider) doRequest(pair *model.CurrencyPair) (*coingeckoResponse, error) {
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

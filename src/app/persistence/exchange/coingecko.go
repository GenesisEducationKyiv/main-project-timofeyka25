package exchange

import (
	"genesis-test/src/app/domain"
	"genesis-test/src/app/service"
	"genesis-test/src/app/utils"
	"genesis-test/src/config"
)

type CoingeckoFactory struct{}

func (f CoingeckoFactory) CreateCoingeckoFactory() service.ExchangeChain {
	return &coingeckoProvider{
		coingeckoURL: config.Get().CoingeckoURL,
	}
}

type coingeckoProvider struct {
	coingeckoURL string
	next         service.ExchangeChain
}

func (c *coingeckoProvider) GetCurrencyRate(pair *domain.CurrencyPair) (*domain.CurrencyRate, error) {
	rate, err := c.getCurrencyRate(pair)
	if err != nil && c.next != nil {
		return c.next.GetCurrencyRate(pair)
	}

	return rate, nil
}

func (c *coingeckoProvider) SetNext(chain service.ExchangeChain) {
	c.next = chain
}

func (c *coingeckoProvider) getCurrencyRate(pair *domain.CurrencyPair) (*domain.CurrencyRate, error) {
	resp, err := c.doRequest(pair)
	if err != nil {
		return nil, err
	}
	return resp.toDefaultRate()
}

func (c *coingeckoProvider) doRequest(pair *domain.CurrencyPair) (*coingeckoResponse, error) {
	rate := new(coingeckoResponse)
	err := utils.GetAndParse(c.coingeckoURL, &rate)
	if err != nil {
		return nil, err
	}
	rate.BaseCurrency = pair.BaseCurrency
	rate.QuoteCurrency = pair.QuoteCurrency

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

func (c *coingeckoResponse) toDefaultRate() (*domain.CurrencyRate, error) {
	return &domain.CurrencyRate{
		Price: c.Rates.QuoteCurrency.Price,
		CurrencyPair: domain.CurrencyPair{
			BaseCurrency:  c.BaseCurrency,
			QuoteCurrency: c.QuoteCurrency,
		},
	}, nil
}

package exchange

import (
	"fmt"
	"genesis-test/src/app/domain"
	"genesis-test/src/app/service"
	"genesis-test/src/app/utils"
	"genesis-test/src/config"
	"strconv"
)

type CoinbaseFactory struct{}

func (f CoinbaseFactory) CreateCoinbaseFactory() service.ExchangeChain {
	return &coinbaseProvider{
		coinbaseURL: config.Get().CoinbaseURL,
	}
}

type coinbaseProvider struct {
	coinbaseURL string
	next        service.ExchangeChain
}

type coinbaseResponse struct {
	Data struct {
		Amount        string `json:"amount"`
		BaseCurrency  string `json:"base"`
		QuoteCurrency string `json:"currency"`
	} `json:"data"`
}

func (c *coinbaseProvider) GetCurrencyRate(pair *domain.CurrencyPair) (*domain.CurrencyRate, error) {
	rate, err := c.getCurrencyRate(pair)
	if err != nil && c.next != nil {
		return c.next.GetCurrencyRate(pair)
	}

	return rate, nil
}

func (c *coinbaseProvider) SetNext(chain service.ExchangeChain) {
	c.next = chain
}

func (c *coinbaseProvider) getCurrencyRate(pair *domain.CurrencyPair) (*domain.CurrencyRate, error) {
	resp, err := c.doRequest(pair)
	if err != nil {
		return nil, err
	}
	return resp.toDefaultRate()
}

func (c *coinbaseProvider) doRequest(pair *domain.CurrencyPair) (*coinbaseResponse, error) {
	url := fmt.Sprintf(
		c.coinbaseURL,
		pair.GetBaseCurrency(),
		pair.GetQuoteCurrency(),
	)
	rate := new(coinbaseResponse)
	err := utils.GetAndParse(url, &rate)
	if err != nil {
		return nil, err
	}

	return rate, nil
}

func (c *coinbaseResponse) toDefaultRate() (*domain.CurrencyRate, error) {
	bitSize := 64
	floatPrice, err := strconv.ParseFloat(c.Data.Amount, bitSize)
	if err != nil {
		return nil, err
	}
	return &domain.CurrencyRate{
		Price: floatPrice,
		CurrencyPair: domain.CurrencyPair{
			BaseCurrency:  c.Data.BaseCurrency,
			QuoteCurrency: c.Data.QuoteCurrency,
		},
	}, nil
}

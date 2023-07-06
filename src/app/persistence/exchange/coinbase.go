package exchange

import (
	"fmt"
	"genesis-test/src/app/application"
	"genesis-test/src/app/domain/model"
	"genesis-test/src/app/utils"
	"genesis-test/src/config"
	"strconv"
)

type CoinbaseFactory struct{}

func (f CoinbaseFactory) CreateCoinbaseFactory() application.ExchangeChain {
	return &coinbaseProvider{
		coinbaseURL: config.Get().CoinbaseURL,
	}
}

type coinbaseProvider struct {
	coinbaseURL string
	next        application.ExchangeChain
}

type coinbaseResponse struct {
	Data struct {
		Amount        string `json:"amount"`
		BaseCurrency  string `json:"base"`
		QuoteCurrency string `json:"currency"`
	} `json:"data"`
}

func (c *coinbaseProvider) GetCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error) {
	rate, err := c.getCurrencyRate(pair)
	if err != nil && c.next != nil {
		return c.next.GetCurrencyRate(pair)
	}

	return rate, nil
}

func (c *coinbaseProvider) SetNext(chain application.ExchangeChain) {
	c.next = chain
}

func (c *coinbaseProvider) getCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error) {
	resp, err := c.doRequest(pair)
	if err != nil {
		return nil, err
	}
	return resp.toDefaultRate()
}

func (c *coinbaseProvider) doRequest(pair *model.CurrencyPair) (*coinbaseResponse, error) {
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

func (c *coinbaseResponse) toDefaultRate() (*model.CurrencyRate, error) {
	bitSize := 64
	floatPrice, err := strconv.ParseFloat(c.Data.Amount, bitSize)
	if err != nil {
		return nil, err
	}
	return &model.CurrencyRate{
		Price: floatPrice,
		CurrencyPair: model.CurrencyPair{
			BaseCurrency:  model.CurrencyFactory{}.CreateCurrency(c.Data.BaseCurrency),
			QuoteCurrency: model.CurrencyFactory{}.CreateCurrency(c.Data.QuoteCurrency),
		},
	}, nil
}

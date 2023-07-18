package exchange

import (
	"fmt"
	"genesis-test/src/app/domain/model"
	"genesis-test/src/app/utils"
	"genesis-test/src/config"
	"strconv"
)

type CoinbaseFactory struct{}

type CoinbaseProvider struct {
	coinbaseURL string
}

func (f CoinbaseFactory) CreateCoinbaseFactory() *CoinbaseProvider {
	return &CoinbaseProvider{
		coinbaseURL: config.Get().CoinbaseURL,
	}
}

type coinbaseResponse struct {
	Data struct {
		Amount        string `json:"amount"`
		BaseCurrency  string `json:"base"`
		QuoteCurrency string `json:"currency"`
	} `json:"data"`
}

func (c *CoinbaseProvider) GetCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error) {
	rate, err := c.getCurrencyRate(pair)
	if err != nil {
		return nil, err
	}

	return rate, nil
}

func (c *CoinbaseProvider) getCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error) {
	resp, err := c.doRequest(pair)
	if err != nil {
		return nil, err
	}
	return resp.toDefaultRate()
}

func (c *CoinbaseProvider) doRequest(pair *model.CurrencyPair) (*coinbaseResponse, error) {
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

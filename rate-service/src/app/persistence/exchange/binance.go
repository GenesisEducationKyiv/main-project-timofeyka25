package exchange

import (
	"fmt"
	"genesis-test/src/app/domain/model"
	"genesis-test/src/app/utils"
	"genesis-test/src/config"
	"strconv"
)

type BinanceFactory struct{}

type BinanceProvider struct {
	binanceURL string
}

func (f BinanceFactory) CreateBinanceFactory() *BinanceProvider {
	return &BinanceProvider{
		binanceURL: config.Get().BinanceURL,
	}
}

func (b *BinanceProvider) GetCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error) {
	rate, err := b.getCurrencyRate(pair)
	if err != nil {
		return nil, err
	}

	return rate, nil
}

func (b *BinanceProvider) getCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error) {
	resp, err := b.doRequest(pair)
	if err != nil {
		return nil, err
	}
	return resp.toDefaultRate()
}

func (b *BinanceProvider) doRequest(pair *model.CurrencyPair) (*binanceResponse, error) {
	url := fmt.Sprintf(
		b.binanceURL,
		pair.GetBaseCurrency(),
		pair.GetQuoteCurrency(),
	)
	rate := new(binanceResponse)
	err := utils.GetAndParse(url, &rate)
	if err != nil {
		return nil, err
	}
	rate.BaseCurrency = pair.BaseCurrency.ToString()
	rate.QuoteCurrency = pair.QuoteCurrency.ToString()

	return rate, nil
}

type binanceResponse struct {
	Price         string `json:"price"`
	BaseCurrency  string
	QuoteCurrency string
}

func (b *binanceResponse) toDefaultRate() (*model.CurrencyRate, error) {
	bitSize := 64
	floatPrice, err := strconv.ParseFloat(b.Price, bitSize)
	if err != nil {
		return nil, err
	}
	return &model.CurrencyRate{
		Price: floatPrice,
		CurrencyPair: model.CurrencyPair{
			BaseCurrency:  model.CurrencyFactory{}.CreateCurrency(b.BaseCurrency),
			QuoteCurrency: model.CurrencyFactory{}.CreateCurrency(b.QuoteCurrency),
		},
	}, nil
}

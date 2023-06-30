package exchange

import (
	"fmt"
	"genesis-test/src/app/domain"
	"genesis-test/src/app/service"
	"genesis-test/src/app/utils"
	"genesis-test/src/config"
	"strconv"
)

type BinanceFactory struct{}

func (f BinanceFactory) CreateBinanceFactory() service.ExchangeChain {
	return &binanceProvider{
		binanceURL: config.Get().BinanceURL,
	}
}

type binanceProvider struct {
	binanceURL string
	next       service.ExchangeChain
}

func (b *binanceProvider) GetCurrencyRate(pair *domain.CurrencyPair) (*domain.CurrencyRate, error) {
	rate, err := b.getCurrencyRate(pair)
	if err != nil && b.next != nil {
		return b.next.GetCurrencyRate(pair)
	}

	return rate, nil
}

func (b *binanceProvider) SetNext(chain service.ExchangeChain) {
	b.next = chain
}

func (b *binanceProvider) getCurrencyRate(pair *domain.CurrencyPair) (*domain.CurrencyRate, error) {
	resp, err := b.doRequest(pair)
	if err != nil {
		return nil, err
	}
	return resp.toDefaultRate()
}

func (b *binanceProvider) doRequest(pair *domain.CurrencyPair) (*binanceResponse, error) {
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
	rate.BaseCurrency = pair.BaseCurrency
	rate.QuoteCurrency = pair.QuoteCurrency

	return rate, nil
}

type binanceResponse struct {
	Price         string `json:"price"`
	BaseCurrency  string
	QuoteCurrency string
}

func (b *binanceResponse) toDefaultRate() (*domain.CurrencyRate, error) {
	bitSize := 64
	floatPrice, err := strconv.ParseFloat(b.Price, bitSize)
	if err != nil {
		return nil, err
	}
	return &domain.CurrencyRate{
		Price: floatPrice,
		CurrencyPair: domain.CurrencyPair{
			BaseCurrency:  b.BaseCurrency,
			QuoteCurrency: b.QuoteCurrency,
		},
	}, nil
}

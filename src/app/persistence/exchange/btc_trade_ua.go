package exchange

import (
	"genesis-test/src/app/domain"
	"genesis-test/src/app/service"
	"genesis-test/src/app/utils"
	"genesis-test/src/config"
	"strconv"
)

type BTCTradeUAFactory struct{}

func (f BTCTradeUAFactory) CreateBTCTradeUAFactory() service.ExchangeChain {
	return &btcTradeUAProvider{
		btcTradeUAURL: config.Get().BTCTradeUAURL,
	}
}

type btcTradeUAProvider struct {
	btcTradeUAURL string
	next          service.ExchangeChain
}

func (b *btcTradeUAProvider) GetCurrencyRate(pair *domain.CurrencyPair) (*domain.CurrencyRate, error) {
	rate, err := b.getCurrencyRate(pair)
	if err != nil && b.next != nil {
		return b.next.GetCurrencyRate(pair)
	}

	return rate, nil
}

func (b *btcTradeUAProvider) SetNext(chain service.ExchangeChain) {
	b.next = chain
}

func (b *btcTradeUAProvider) getCurrencyRate(pair *domain.CurrencyPair) (*domain.CurrencyRate, error) {
	resp, err := b.doRequest(pair)
	if err != nil {
		return nil, err
	}
	return resp.toDefaultRate()
}

func (b *btcTradeUAProvider) doRequest(pair *domain.CurrencyPair) (*btcTradeUAResponse, error) {
	rate := new(btcTradeUAResponse)
	err := utils.GetAndParse(b.btcTradeUAURL, &rate)
	if err != nil {
		return nil, err
	}
	rate.BaseCurrency = pair.BaseCurrency
	rate.QuoteCurrency = pair.QuoteCurrency

	return rate, nil
}

type btcTradeUAResponse struct {
	BtcUah struct {
		Price string `json:"sell"`
	} `json:"btc_uah"`
	BaseCurrency  string `json:"currency_trade"`
	QuoteCurrency string `json:"currency_base"`
}

func (b *btcTradeUAResponse) toDefaultRate() (*domain.CurrencyRate, error) {
	bitSize := 64
	floatPrice, err := strconv.ParseFloat(b.BtcUah.Price, bitSize)
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

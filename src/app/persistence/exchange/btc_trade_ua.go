package exchange

import (
	"genesis-test/src/app/application"
	"genesis-test/src/app/domain/model"
	"genesis-test/src/app/utils"
	"genesis-test/src/config"
	"strconv"
)

type BTCTradeUAFactory struct{}

func (f BTCTradeUAFactory) CreateBTCTradeUAFactory() application.ExchangeChain {
	return &btcTradeUAProvider{
		btcTradeUAURL: config.Get().BTCTradeUAURL,
	}
}

type btcTradeUAProvider struct {
	btcTradeUAURL string
	next          application.ExchangeChain
}

func (b *btcTradeUAProvider) GetCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error) {
	rate, err := b.getCurrencyRate(pair)
	if err != nil && b.next != nil {
		return b.next.GetCurrencyRate(pair)
	}

	return rate, nil
}

func (b *btcTradeUAProvider) SetNext(chain application.ExchangeChain) {
	b.next = chain
}

func (b *btcTradeUAProvider) getCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error) {
	resp, err := b.doRequest(pair)
	if err != nil {
		return nil, err
	}
	return resp.toDefaultRate()
}

func (b *btcTradeUAProvider) doRequest(pair *model.CurrencyPair) (*btcTradeUAResponse, error) {
	rate := new(btcTradeUAResponse)
	err := utils.GetAndParse(b.btcTradeUAURL, &rate)
	if err != nil {
		return nil, err
	}
	rate.BaseCurrency = pair.BaseCurrency.ToString()
	rate.QuoteCurrency = pair.QuoteCurrency.ToString()

	return rate, nil
}

type btcTradeUAResponse struct {
	BtcUah struct {
		Price string `json:"sell"`
	} `json:"btc_uah"`
	BaseCurrency  string `json:"currency_trade"`
	QuoteCurrency string `json:"currency_base"`
}

func (b *btcTradeUAResponse) toDefaultRate() (*model.CurrencyRate, error) {
	bitSize := 64
	floatPrice, err := strconv.ParseFloat(b.BtcUah.Price, bitSize)
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

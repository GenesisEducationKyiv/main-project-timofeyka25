package exchange

import (
	"genesis-test/src/app/domain/model"
	"genesis-test/src/app/utils"
	"genesis-test/src/config"
	"strconv"
)

type BTCTradeUAFactory struct{}

type BTCTradeUAProvider struct {
	btcTradeUAURL string
}

func (f BTCTradeUAFactory) CreateBTCTradeUAFactory() *BTCTradeUAProvider {
	return &BTCTradeUAProvider{
		btcTradeUAURL: config.Get().BTCTradeUAURL,
	}
}

func (b *BTCTradeUAProvider) GetCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error) {
	rate, err := b.getCurrencyRate(pair)
	if err != nil {
		return nil, err
	}

	return rate, nil
}

func (b *BTCTradeUAProvider) getCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error) {
	resp, err := b.doRequest(pair)
	if err != nil {
		return nil, err
	}
	return resp.toDefaultRate()
}

func (b *BTCTradeUAProvider) doRequest(pair *model.CurrencyPair) (*btcTradeUAResponse, error) {
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

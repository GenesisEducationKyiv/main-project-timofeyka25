package exchange

import (
	"genesis-test/src/app/application"
	"genesis-test/src/app/domain/model"
)

type providerNode struct {
	provider application.ExchangeProvider
	next     application.ExchangeProvider
}

func NewProviderNode(provider application.ExchangeProvider) application.ExchangeProviderNode {
	return &providerNode{
		provider: provider,
	}
}

func (c *providerNode) GetCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error) {
	rate, err := c.provider.GetCurrencyRate(pair)
	if err != nil && c.next != nil {
		return c.next.GetCurrencyRate(pair)
	}

	return rate, nil
}

func (c *providerNode) SetNext(provider application.ExchangeProvider) {
	c.next = provider
}

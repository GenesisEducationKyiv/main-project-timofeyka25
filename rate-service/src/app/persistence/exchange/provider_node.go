package exchange

import (
	"genesis-test/src/app/domain/model"
)

type ProviderNode struct {
	provider ExchangeProvider
	next     ExchangeProvider
}

func NewProviderNode(provider ExchangeProvider) *ProviderNode {
	return &ProviderNode{
		provider: provider,
	}
}

func (c *ProviderNode) GetCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error) {
	rate, err := c.provider.GetCurrencyRate(pair)
	if err != nil {
		if c.next != nil {
			return c.next.GetCurrencyRate(pair)
		}
		return nil, err
	}

	return rate, nil
}

func (c *ProviderNode) SetNext(provider ExchangeProvider) {
	c.next = provider
}

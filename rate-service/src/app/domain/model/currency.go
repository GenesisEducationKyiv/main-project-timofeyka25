package model

type Currency string

type CurrencyPair struct {
	BaseCurrency  Currency
	QuoteCurrency Currency
}

type CurrencyFactory struct{}

func (cf CurrencyFactory) CreateCurrency(currency string) Currency {
	return Currency(currency)
}

func (c Currency) ToString() string {
	return string(c)
}

func (c CurrencyPair) GetBaseCurrency() Currency {
	return c.BaseCurrency
}

func (c CurrencyPair) GetQuoteCurrency() Currency {
	return c.QuoteCurrency
}

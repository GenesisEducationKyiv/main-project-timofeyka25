package domain

//go:generate mockgen -destination=mocks/mock_exchange.go genesis-test/src/app/domain ExchangeRepository,ExchangeService

type CurrencyPair struct {
	BaseCurrency  string
	QuoteCurrency string
}

type CurrencyRate struct {
	CurrencyPair
	Price float64
}

func (c CurrencyPair) GetBaseCurrency() string {
	return c.BaseCurrency
}

func (c CurrencyPair) GetQuoteCurrency() string {
	return c.QuoteCurrency
}

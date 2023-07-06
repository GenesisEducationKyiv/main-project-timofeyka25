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

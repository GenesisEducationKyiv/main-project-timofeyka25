package domain

import "genesis-test/src/config"

//go:generate mockgen -destination=mocks/mock_exchange.go genesis-test/src/app/domain ExchangeRepository,ExchangeService

type CurrencyRate struct {
	Price         string `json:"amount"`
	BaseCurrency  string `json:"base"`
	QuoteCurrency string `json:"currency"`
}

type ExchangeService interface {
	GetCurrencyRate(cfg *config.Config) (int, error)
}

type ExchangeRepository interface {
	GetCurrencyRate(base string, quoted string) (*CurrencyRate, error)
}

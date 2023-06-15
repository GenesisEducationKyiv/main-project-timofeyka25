package domain

type CurrencyRate struct {
	Price         string `json:"amount"`
	BaseCurrency  string `json:"base"`
	QuoteCurrency string `json:"currency"`
}

type ExchangeService interface {
	GetCurrencyRate() (int, error)
}

type ExchangeRepository interface {
	GetCurrencyRate(base string, quoted string) (*CurrencyRate, error)
}

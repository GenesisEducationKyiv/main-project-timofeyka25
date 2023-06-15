package service

import (
	"genesis-test/src/app/domain"
	"genesis-test/src/app/repository"
	"genesis-test/src/config"
	"strconv"
)

type exchangeService struct {
	repos *repository.Repositories
}

func NewExchangeService(r *repository.Repositories) domain.ExchangeService {
	return &exchangeService{
		repos: r,
	}
}

func (c exchangeService) GetCurrencyRate() (int, error) {
	cfg := config.Get()
	rate, err := c.repos.Exchange.GetCurrencyRate(cfg.BaseCurrency, cfg.QuoteCurrency)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(rate.Price)
}

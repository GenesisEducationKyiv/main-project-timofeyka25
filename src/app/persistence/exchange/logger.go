package exchange

import (
	"genesis-test/src/app/domain"
	"genesis-test/src/app/service"
	"genesis-test/src/logger"
)

type exchangeLogger struct {
	logger logger.Logger
}

func NewExchangeLogger(logger logger.Logger) service.ExchangeLogger {
	return &exchangeLogger{
		logger: logger,
	}
}

func (e exchangeLogger) LogExchangeRate(provider string, rate *domain.CurrencyRate) {
	e.logger.Infow(
		"received rate",
		"provider", provider,
		"price", rate.Price,
		"base", rate.GetBaseCurrency(),
		"quote", rate.GetQuoteCurrency(),
	)
}

package exchange

import (
	"genesis-test/src/app/application"
	"genesis-test/src/app/domain/model"
	"genesis-test/src/logger"
)

type exchangeLogger struct {
	logger logger.Logger
}

func NewExchangeLogger(logger logger.Logger) application.ExchangeLogger {
	return &exchangeLogger{
		logger: logger,
	}
}

func (e exchangeLogger) LogExchangeRate(provider string, rate *model.CurrencyRate) {
	e.logger.Infow(
		"received rate",
		"provider", provider,
		"price", rate.Price,
		"base", rate.GetBaseCurrency(),
		"quote", rate.GetQuoteCurrency(),
	)
}

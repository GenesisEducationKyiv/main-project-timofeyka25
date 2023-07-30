package logger

import (
	"encoding/json"
	"genesis-test/src/app/domain/model"
	"log"
)

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Error(msg string)
}

type ExchangeLogger struct {
	logger Logger
}

func NewExchangeLogger(logger Logger) *ExchangeLogger {
	return &ExchangeLogger{
		logger: logger,
	}
}

func (e ExchangeLogger) LogExchangeRate(provider string, rate *model.CurrencyRate) {
	marshal, err := json.Marshal(map[string]any{
		"provider": provider,
		"price":    rate.Price,
		"base":     rate.GetBaseCurrency(),
		"quote":    rate.GetQuoteCurrency(),
	})
	if err != nil {
		log.Println("failed to marshal exchange log")
	}
	e.logger.Info(string(marshal))
}

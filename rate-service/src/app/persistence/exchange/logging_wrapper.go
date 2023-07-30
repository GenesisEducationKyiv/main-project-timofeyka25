package exchange

import (
	"genesis-test/src/app/domain/model"
	"reflect"
)

type ExchangeProvider interface {
	GetCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error)
}

type ExchangeLogger interface {
	LogExchangeRate(provider string, rate *model.CurrencyRate)
}

type LoggingWrapper struct {
	provider ExchangeProvider
	logger   ExchangeLogger
}

func NewLoggingWrapper(
	provider ExchangeProvider,
	logger ExchangeLogger,
) *LoggingWrapper {
	return &LoggingWrapper{
		provider: provider,
		logger:   logger,
	}
}

func (l *LoggingWrapper) GetCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error) {
	rate, err := l.provider.GetCurrencyRate(pair)
	if err != nil {
		return nil, err
	}
	l.logger.LogExchangeRate(l.getProviderName(), rate)

	return rate, nil
}

func (l *LoggingWrapper) getProviderName() string {
	return reflect.TypeOf(l.provider).Elem().Name()
}

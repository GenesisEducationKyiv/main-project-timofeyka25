package exchange

import (
	"genesis-test/src/app/application"
	"genesis-test/src/app/domain/model"
	"reflect"
)

type loggingWrapper struct {
	provider application.ExchangeProvider
	logger   application.ExchangeLogger
}

func NewLoggingWrapper(provider application.ExchangeProvider, logger application.ExchangeLogger) application.ExchangeProvider {
	return &loggingWrapper{
		provider: provider,
		logger:   logger,
	}
}

func (l *loggingWrapper) GetCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error) {
	rate, err := l.provider.GetCurrencyRate(pair)
	if err != nil {
		return nil, err
	}
	l.logger.LogExchangeRate(l.getProviderName(), rate)

	return rate, nil
}

func (l *loggingWrapper) getProviderName() string {
	return reflect.TypeOf(l.provider).Elem().Name()
}

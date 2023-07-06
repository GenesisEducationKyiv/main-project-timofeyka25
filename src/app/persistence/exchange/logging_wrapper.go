package exchange

import (
	"genesis-test/src/app/application"
	"genesis-test/src/app/domain/model"
	"reflect"
)

type loggingWrapper struct {
	chain  application.ExchangeChain
	logger application.ExchangeLogger
}

func NewLoggingWrapper(chain application.ExchangeChain, logger application.ExchangeLogger) application.ExchangeChain {
	return &loggingWrapper{
		chain:  chain,
		logger: logger,
	}
}

func (l *loggingWrapper) GetCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error) {
	rate, err := l.chain.GetCurrencyRate(pair)
	if err != nil {
		return nil, err
	}
	l.logger.LogExchangeRate(l.getProviderName(), rate)

	return rate, nil
}

func (l *loggingWrapper) SetNext(chain application.ExchangeChain) {
	l.chain.SetNext(chain)
}

func (l *loggingWrapper) getProviderName() string {
	return reflect.TypeOf(l.chain).Elem().Name()
}

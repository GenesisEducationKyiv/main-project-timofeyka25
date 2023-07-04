package exchange

import (
	"genesis-test/src/app/domain"
	"genesis-test/src/app/service"
	"reflect"
)

type loggingWrapper struct {
	chain  service.ExchangeChain
	logger service.ExchangeLogger
}

func NewLoggingWrapper(chain service.ExchangeChain, logger service.ExchangeLogger) service.ExchangeChain {
	return &loggingWrapper{
		chain:  chain,
		logger: logger,
	}
}

func (l *loggingWrapper) GetCurrencyRate(pair *domain.CurrencyPair) (*domain.CurrencyRate, error) {
	rate, err := l.chain.GetCurrencyRate(pair)
	if err != nil {
		return nil, err
	}
	l.logger.LogExchangeRate(l.getProviderName(), rate)

	return rate, nil
}

func (l *loggingWrapper) SetNext(chain service.ExchangeChain) {
	l.chain.SetNext(chain)
}

func (l *loggingWrapper) getProviderName() string {
	return reflect.TypeOf(l.chain).Elem().Name()
}

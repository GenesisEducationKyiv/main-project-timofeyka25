package service

import "genesis-test/src/app/domain"

//go:generate mockgen -destination=../domain/mocks/mock_persistence.go genesis-test/src/app/service NewsletterSender,EmailStorage,ExchangeChain,ExchangeLogger

type NewsletterSender interface {
	MultipleSending(subscribers []string, message *domain.EmailMessage) ([]string, error)
}

type EmailStorage interface {
	GetAllEmails() ([]string, error)
	AddEmail(newEmail string) error
}

type ExchangeChain interface {
	GetCurrencyRate(pair *domain.CurrencyPair) (*domain.CurrencyRate, error)
	SetNext(provider ExchangeChain)
}

type ExchangeLogger interface {
	LogExchangeRate(provider string, rate *domain.CurrencyRate)
}

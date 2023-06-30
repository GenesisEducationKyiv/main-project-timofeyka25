package service

import "genesis-test/src/app/domain"

//go:generate mockgen -destination=../domain/mocks/mock_repositories.go genesis-test/src/app/service NewsletterRepository,EmailStorage,ExchangeChain,ExchangeLogger

type NewsletterRepository interface {
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

type Repositories struct {
	Newsletter NewsletterRepository
	Storage    EmailStorage
	Exchange   ExchangeChain
}

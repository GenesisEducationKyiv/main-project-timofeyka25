package application

import (
	"genesis-test/src/app/domain/model"
)

//go:generate mockgen -destination=../domain/mocks/mock_persistence.go genesis-test/src/app/service NewsletterSender,EmailStorage,ExchangeChain,ExchangeLogger

type NewsletterSender interface {
	MultipleSending(subscribers []string, message *model.EmailMessage) ([]string, error)
}

type EmailStorage interface {
	GetAllEmails() ([]string, error)
	AddEmail(newEmail string) error
}

type ExchangeChain interface {
	GetCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error)
	SetNext(provider ExchangeChain)
}

type ExchangeLogger interface {
	LogExchangeRate(provider string, rate *model.CurrencyRate)
}

type Persistence struct {
	Sender    NewsletterSender
	Storage   EmailStorage
	Providers ExchangeChain
}

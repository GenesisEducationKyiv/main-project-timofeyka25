package service

import (
	"genesis-test/src/app/domain"
	"genesis-test/src/app/repository"
)

type Services struct {
	Newsletter domain.NewsletterService
	Exchange   domain.ExchangeService
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Newsletter: NewNewsletterService(repos),
		Exchange:   NewExchangeService(repos),
	}
}

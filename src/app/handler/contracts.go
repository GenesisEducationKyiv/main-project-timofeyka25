package handler

import "genesis-test/src/app/domain"

//go:generate mockgen -destination=../domain/mocks/mock_services.go genesis-test/src/app/handler ExchangeService,NewsletterService,SubscriptionService

type ExchangeService interface {
	GetCurrencyRate() (float64, error)
}

type NewsletterService interface {
	SendCurrencyRate() ([]string, error)
}

type SubscriptionService interface {
	Subscribe(subscriber *domain.Subscriber) error
}

type Services struct {
	Subscription SubscriptionService
	Newsletter   NewsletterService
	Exchange     ExchangeService
}

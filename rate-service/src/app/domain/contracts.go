package domain

import (
	"genesis-test/src/app/domain/model"
)

type ExchangeService interface {
	GetCurrencyRate(pair *model.CurrencyPair) (float64, error)
}

type NewsletterService interface {
	SendCurrencyRate() ([]string, error)
}

type SubscriptionService interface {
	Subscribe(subscriber *model.Subscriber) error
}

type Services struct {
	Exchange     ExchangeService
	Newsletter   NewsletterService
	Subscription SubscriptionService
}

package service

import (
	"fmt"
	"genesis-test/src/app/domain"
	"genesis-test/src/app/handler"

	"github.com/pkg/errors"
)

type newsletterService struct {
	repos *Repositories
	pair  *domain.CurrencyPair
}

func NewNewsletterService(
	repos *Repositories,
	pair *domain.CurrencyPair,
) handler.NewsletterService {
	return &newsletterService{
		repos: repos,
		pair:  pair,
	}
}

func (s *newsletterService) SendCurrencyRate() ([]string, error) {
	rate, err := s.repos.Exchange.GetCurrencyRate(s.pair)
	if err != nil {
		return nil, errors.Wrap(err, "get rate")
	}

	return s.sendToSubscribed(s.buildMessage(rate))
}

func (s *newsletterService) sendToSubscribed(message *domain.EmailMessage) ([]string, error) {
	subscribers, err := s.repos.Storage.GetAllEmails()
	if err != nil {
		return nil, errors.Wrap(err, "send to subscribed")
	}
	return s.repos.Newsletter.MultipleSending(subscribers, message)
}

func (s *newsletterService) buildMessage(rate *domain.CurrencyRate) *domain.EmailMessage {
	return &domain.EmailMessage{
		Subject: "Crypto Exchange Newsletter",
		Body: fmt.Sprintf("The current exchange rate of %s to %s is %f %s",
			rate.BaseCurrency,
			rate.QuoteCurrency,
			rate.Price,
			rate.QuoteCurrency),
	}
}

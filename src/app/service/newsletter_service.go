package service

import (
	"fmt"
	"genesis-test/src/app/customerror"
	"genesis-test/src/app/domain"
	"genesis-test/src/app/repository"
	"genesis-test/src/config"

	"github.com/pkg/errors"
)

type newsletterService struct {
	repos *repository.Repositories
}

func NewNewsletterService(r *repository.Repositories) domain.NewsletterService {
	return newsletterService{repos: r}
}

func (m newsletterService) SendEmails() ([]string, error) {
	cfg := config.Get()

	rate, err := m.repos.Exchange.GetCurrencyRate(cfg.BaseCurrency, cfg.QuoteCurrency)
	if err != nil {
		return nil, errors.Wrap(err, "get rate")
	}

	body := fmt.Sprintf("The current exchange rate of %s to %s is %s %s",
		rate.BaseCurrency,
		rate.QuoteCurrency,
		rate.Price,
		rate.QuoteCurrency)

	return m.repos.Newsletter.SendToSubscribedEmails(body)
}

func (m newsletterService) Subscribe(subscriber *domain.Subscriber) error {
	if subscriber == nil {
		return customerror.ErrNoDataProvided
	}
	subscribed, err := m.repos.Newsletter.GetSubscribedEmails()
	if err != nil && !errors.Is(err, customerror.ErrNoSubscribers) {
		return errors.Wrap(err, "get emails")
	}

	err = m.repos.Newsletter.AddNewEmail(subscribed, subscriber.Email)
	if err != nil {
		return errors.Wrap(err, "add email")
	}

	return nil
}

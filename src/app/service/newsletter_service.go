package service

import (
	"fmt"
	"genesis-test/src/app/domain"
	"genesis-test/src/app/repository"
	"genesis-test/src/config"
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
		return nil, err
	}

	// Prepare email body with the currency rate information
	body := fmt.Sprintf("The current exchange rate of %s to %s is %s %s",
		rate.BaseCurrency,
		rate.QuoteCurrency,
		rate.Price,
		rate.QuoteCurrency)

	unsent, err := m.repos.Newsletter.SendToSubscribedEmails("Exchange Currency Newsletter", body)

	return unsent, err
}

func (m newsletterService) Subscribe(recipient *domain.Subscriber) error {
	subscribed, err := m.repos.Newsletter.GetSubscribedEmails()
	if err != nil {
		return err
	}

	err = m.repos.Newsletter.AddNewEmail(subscribed, recipient.Email)
	if err != nil {
		return err
	}

	return nil
}

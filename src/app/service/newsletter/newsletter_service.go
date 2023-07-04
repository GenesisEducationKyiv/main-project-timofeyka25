package newsletter

import (
	"fmt"
	"genesis-test/src/app/domain"
	"genesis-test/src/app/handler"
	"genesis-test/src/app/service"

	"github.com/pkg/errors"
)

type newsletterService struct {
	exchangeChain service.ExchangeChain
	storage       service.EmailStorage
	sender        service.NewsletterSender
	pair          *domain.CurrencyPair
}

func NewNewsletterService(
	exchangeChain service.ExchangeChain,
	storage service.EmailStorage,
	sender service.NewsletterSender,
	pair *domain.CurrencyPair,
) handler.NewsletterService {
	return &newsletterService{
		exchangeChain: exchangeChain,
		storage:       storage,
		sender:        sender,
		pair:          pair,
	}
}

func (s *newsletterService) SendCurrencyRate() ([]string, error) {
	rate, err := s.exchangeChain.GetCurrencyRate(s.pair)
	if err != nil {
		return nil, errors.Wrap(err, "get rate")
	}

	return s.sendToSubscribed(s.buildMessage(rate))
}

func (s *newsletterService) sendToSubscribed(message *domain.EmailMessage) ([]string, error) {
	subscribers, err := s.storage.GetAllEmails()
	if err != nil {
		return nil, errors.Wrap(err, "send to subscribed")
	}
	return s.sender.MultipleSending(subscribers, message)
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

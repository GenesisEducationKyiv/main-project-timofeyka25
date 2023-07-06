package newsletter

import (
	"fmt"
	"genesis-test/src/app/application"
	"genesis-test/src/app/domain"
	"genesis-test/src/app/domain/model"
	"github.com/pkg/errors"
)

type newsletterService struct {
	exchangeChain application.ExchangeChain
	storage       application.EmailStorage
	sender        application.NewsletterSender
	pair          *model.CurrencyPair
}

func NewNewsletterService(
	exchangeChain application.ExchangeChain,
	storage application.EmailStorage,
	sender application.NewsletterSender,
	pair *model.CurrencyPair,
) domain.NewsletterService {
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

func (s *newsletterService) sendToSubscribed(message *model.EmailMessage) ([]string, error) {
	subscribers, err := s.storage.GetAllEmails()
	if err != nil {
		return nil, errors.Wrap(err, "send to subscribed")
	}
	return s.sender.MultipleSending(subscribers, message)
}

func (s *newsletterService) buildMessage(rate *model.CurrencyRate) *model.EmailMessage {
	return &model.EmailMessage{
		Subject: "Crypto Exchange Newsletter",
		Body: fmt.Sprintf("The current exchange rate of %s to %s is %f %s",
			rate.BaseCurrency,
			rate.QuoteCurrency,
			rate.Price,
			rate.QuoteCurrency),
	}
}

package newsletter

import (
	"fmt"
	"genesis-test/src/app/domain/model"

	"github.com/pkg/errors"
)

type ExchangeProvider interface {
	GetCurrencyRate(pair *model.CurrencyPair) (*model.CurrencyRate, error)
}

type EmailStorage interface {
	GetAllEmails() ([]string, error)
	AddEmail(newEmail string) error
}

type NewsletterSender interface {
	MultipleSending(subscribers []string, message *model.EmailMessage) ([]string, error)
}

type NewsletterService struct {
	exchangeProvider ExchangeProvider
	storage          EmailStorage
	sender           NewsletterSender
	pair             *model.CurrencyPair
}

func NewNewsletterService(
	exchangeProvider ExchangeProvider,
	storage EmailStorage,
	sender NewsletterSender,
	pair *model.CurrencyPair,
) *NewsletterService {
	return &NewsletterService{
		exchangeProvider: exchangeProvider,
		storage:          storage,
		sender:           sender,
		pair:             pair,
	}
}

func (s *NewsletterService) SendCurrencyRate() ([]string, error) {
	rate, err := s.exchangeProvider.GetCurrencyRate(s.pair)
	if err != nil {
		return nil, errors.Wrap(err, "get rate")
	}

	return s.sendToSubscribed(s.buildMessage(rate))
}

func (s *NewsletterService) sendToSubscribed(message *model.EmailMessage) ([]string, error) {
	subscribers, err := s.storage.GetAllEmails()
	if err != nil {
		return nil, errors.Wrap(err, "send to subscribed")
	}
	return s.sender.MultipleSending(subscribers, message)
}

func (s *NewsletterService) buildMessage(rate *model.CurrencyRate) *model.EmailMessage {
	return &model.EmailMessage{
		Subject: "Crypto Exchange Newsletter",
		Body: fmt.Sprintf("The current exchange rate of %s to %s is %f %s",
			rate.BaseCurrency,
			rate.QuoteCurrency,
			rate.Price,
			rate.QuoteCurrency),
	}
}

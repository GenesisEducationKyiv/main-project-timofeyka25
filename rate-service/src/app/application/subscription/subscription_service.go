package subscription

import (
	"genesis-test/src/app/customerror"
	"genesis-test/src/app/domain/model"
)

type EmailStorage interface {
	GetAllEmails() ([]string, error)
	AddEmail(newEmail string) error
}

type SubscriptionService struct {
	storage EmailStorage
}

func NewSubscriptionService(storage EmailStorage) *SubscriptionService {
	return &SubscriptionService{
		storage: storage,
	}
}

func (s SubscriptionService) Subscribe(subscriber *model.Subscriber) error {
	if subscriber == nil {
		return customerror.ErrNoDataProvided
	}
	return s.storage.AddEmail(subscriber.Email)
}

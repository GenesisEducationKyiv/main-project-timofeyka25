package subscription

import (
	"genesis-test/src/app/application"
	"genesis-test/src/app/customerror"
	"genesis-test/src/app/domain"
	"genesis-test/src/app/domain/model"
)

type subscriptionService struct {
	storage application.EmailStorage
}

func NewSubscriptionService(storage application.EmailStorage) domain.SubscriptionService {
	return &subscriptionService{
		storage: storage,
	}
}

func (s subscriptionService) Subscribe(subscriber *model.Subscriber) error {
	if subscriber == nil {
		return customerror.ErrNoDataProvided
	}
	return s.storage.AddEmail(subscriber.Email)
}

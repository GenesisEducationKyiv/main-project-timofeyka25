package service

import (
	"genesis-test/src/app/customerror"
	"genesis-test/src/app/domain"
	"genesis-test/src/app/handler"
)

type subscriptionService struct {
	storage EmailStorage
}

func NewSubscriptionService(storage EmailStorage) handler.SubscriptionService {
	return &subscriptionService{
		storage: storage,
	}
}

func (s subscriptionService) Subscribe(subscriber *domain.Subscriber) error {
	if subscriber == nil {
		return customerror.ErrNoDataProvided
	}
	return s.storage.AddEmail(subscriber.Email)
}

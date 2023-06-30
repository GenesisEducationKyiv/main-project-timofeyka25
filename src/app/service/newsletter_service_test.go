package service

import (
	"fmt"
	"genesis-test/src/app/customerror"
	"genesis-test/src/app/domain"
	mocks "genesis-test/src/app/domain/mocks"
	"genesis-test/src/app/repository"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestNewsletterService_Subscribe(t *testing.T) {
	type mockBehavior func(r *mocks.MockNewsletterRepository, m []string, s *domain.Subscriber)

	cases := []struct {
		name          string
		subscriber    *domain.Subscriber
		mockBehavior  mockBehavior
		mockResponse  []string
		expectedError error
	}{
		{
			name: "Subscribe successful",
			subscriber: &domain.Subscriber{
				Email: "test@testexample.com",
			},
			mockResponse: []string{
				"123@test.com",
				"abc@example.com",
			},
			mockBehavior: func(
				r *mocks.MockNewsletterRepository,
				m []string,
				s *domain.Subscriber,
			) {
				r.EXPECT().GetSubscribedEmails().Return(m, nil)
				r.EXPECT().AddNewEmail(m, s.Email).Return(nil)
			},
		},
		{
			name:          "Subscribe error (no data)",
			mockBehavior:  func(r *mocks.MockNewsletterRepository, m []string, s *domain.Subscriber) {},
			expectedError: customerror.ErrNoDataProvided,
		},
		{
			name: "Subscribe error (already exists)",
			subscriber: &domain.Subscriber{
				Email: "test@testexample.com",
			},
			mockResponse: []string{
				"123@test.com",
				"abc@example.com",
			},
			mockBehavior: func(r *mocks.MockNewsletterRepository, m []string, s *domain.Subscriber) {
				r.EXPECT().GetSubscribedEmails().Return(m, nil)
				r.EXPECT().AddNewEmail(m, s.Email).Return(customerror.ErrAlreadyExists)
			},
			expectedError: customerror.ErrAlreadyExists,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockNewsletterRepo := mocks.NewMockNewsletterRepository(ctrl)
			c.mockBehavior(mockNewsletterRepo, c.mockResponse, c.subscriber)

			repos := &repository.Repositories{
				Newsletter: mockNewsletterRepo,
			}

			newsletterTestService := NewNewsletterService(repos)
			err := newsletterTestService.Subscribe(c.subscriber)
			require.ErrorIs(t, err, c.expectedError)
		})
	}
}

func TestNewsletterService_SendEmails(t *testing.T) {
	type mockBehavior func(mockExchangeRepo *mocks.MockExchangeRepository,
		mockNewsletterRepo *mocks.MockNewsletterRepository,
		mockExchangeResp *domain.CurrencyRate,
		mockNewsletterResp []string,
	)

	cases := []struct {
		name                   string
		mockExchangeResponse   *domain.CurrencyRate
		mockNewsletterResponse []string
		mockBehavior           mockBehavior
		isErrorExpected        bool
	}{
		{
			name: "OK",
			mockExchangeResponse: &domain.CurrencyRate{
				Price:         "",
				BaseCurrency:  "",
				QuoteCurrency: "",
			},
			mockNewsletterResponse: []string{"abc@test.com"},
			mockBehavior: func(
				mockExchangeRepo *mocks.MockExchangeRepository,
				mockNewsletterRepo *mocks.MockNewsletterRepository,
				mockExchangeResp *domain.CurrencyRate,
				mockNewsletterResp []string,
			) {
				mockExchangeRepo.EXPECT().GetCurrencyRate("", "").Return(mockExchangeResp, nil)
				mockNewsletterRepo.EXPECT().SendToSubscribedEmails(fmt.Sprintf("The current exchange rate of %s to %s is %s %s",
					mockExchangeResp.BaseCurrency,
					mockExchangeResp.QuoteCurrency,
					mockExchangeResp.Price,
					mockExchangeResp.QuoteCurrency)).Return(mockNewsletterResp, nil)
			},
			isErrorExpected: false,
		},
		{
			name: "any error case",
			mockBehavior: func(
				mockExchangeRepo *mocks.MockExchangeRepository,
				mockNewsletterRepo *mocks.MockNewsletterRepository,
				mockExchangeResp *domain.CurrencyRate,
				mockNewsletterResp []string,
			) {
				mockExchangeRepo.EXPECT().GetCurrencyRate("", "").Return(nil,
					errors.New("any error"))
			},
			isErrorExpected: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockNewsletterRepo := mocks.NewMockNewsletterRepository(ctrl)
			mockExchangeRepo := mocks.NewMockExchangeRepository(ctrl)

			repos := &repository.Repositories{
				Newsletter: mockNewsletterRepo,
				Exchange:   mockExchangeRepo,
			}

			newsletterTestService := NewNewsletterService(repos)

			c.mockBehavior(mockExchangeRepo, mockNewsletterRepo, c.mockExchangeResponse, c.mockNewsletterResponse)

			if !c.isErrorExpected {
				unsent, err := newsletterTestService.SendEmails()
				require.NoError(t, err)
				require.Equal(t, unsent, c.mockNewsletterResponse)
			} else {
				_, err := newsletterTestService.SendEmails()
				require.Error(t, err)
			}
		})
	}
}

package newsletter

import (
	"genesis-test/src/app/domain"
	"genesis-test/src/app/domain/mocks"
	"genesis-test/src/app/service"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestNewsletterService_SendEmails(t *testing.T) {
	type mockBehavior func(mockExchangeChain *mocks.MockExchangeChain,
		mockNewsletterRepo *mocks.MockNewsletterRepository,
		mockEmailStorage *mocks.MockEmailStorage,
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
				Price: 123,
				CurrencyPair: domain.CurrencyPair{
					BaseCurrency:  "BTC",
					QuoteCurrency: "UAH",
				},
			},
			mockNewsletterResponse: []string{"abc@test.com"},
			mockBehavior: func(
				mockExchangeChain *mocks.MockExchangeChain,
				mockNewsletterRepo *mocks.MockNewsletterRepository,
				mockEmailStorage *mocks.MockEmailStorage,
				mockExchangeResp *domain.CurrencyRate,
				mockNewsletterResp []string,
			) {
				mockExchangeChain.EXPECT().GetCurrencyRate(&domain.CurrencyPair{
					BaseCurrency:  "BTC",
					QuoteCurrency: "UAH",
				}).Return(mockExchangeResp, nil)
				mockEmailStorage.EXPECT().GetAllEmails().Return([]string{"abc@test.com"}, nil)
				mockNewsletterRepo.EXPECT().MultipleSending([]string{"abc@test.com"}, &domain.EmailMessage{
					Subject: "Crypto Exchange Newsletter",
					Body:    "The current exchange rate of BTC to UAH is 123.000000 UAH",
				}).Return(mockNewsletterResp, nil)
			},
			isErrorExpected: false,
		},
		{
			name: "any error case",
			mockBehavior: func(
				mockExchangeChain *mocks.MockExchangeChain,
				mockNewsletterRepo *mocks.MockNewsletterRepository,
				mockEmailStorage *mocks.MockEmailStorage,
				mockExchangeResp *domain.CurrencyRate,
				mockNewsletterResp []string,
			) {
				mockExchangeChain.EXPECT().GetCurrencyRate(&domain.CurrencyPair{
					BaseCurrency:  "BTC",
					QuoteCurrency: "UAH",
				}).Return(nil,
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
			mockExchangeChain := mocks.NewMockExchangeChain(ctrl)
			mockEmailStorage := mocks.NewMockEmailStorage(ctrl)

			repos := &service.Repositories{
				Newsletter: mockNewsletterRepo,
				Storage:    mockEmailStorage,
				Exchange:   mockExchangeChain,
			}
			newsletterTestService := NewNewsletterService(repos, &domain.CurrencyPair{
				BaseCurrency:  "BTC",
				QuoteCurrency: "UAH",
			})

			c.mockBehavior(
				mockExchangeChain,
				mockNewsletterRepo,
				mockEmailStorage,
				c.mockExchangeResponse,
				c.mockNewsletterResponse,
			)

			if !c.isErrorExpected {
				unsent, err := newsletterTestService.SendCurrencyRate()
				require.NoError(t, err)
				require.Equal(t, unsent, c.mockNewsletterResponse)
			} else {
				_, err := newsletterTestService.SendCurrencyRate()
				require.Error(t, err)
			}
		})
	}
}
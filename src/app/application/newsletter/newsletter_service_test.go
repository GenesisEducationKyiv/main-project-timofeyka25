package newsletter

import (
	"genesis-test/src/app/application/mocks"
	"genesis-test/src/app/domain/model"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestNewsletterService_SendEmails(t *testing.T) {
	type mockBehavior func(mockExchangeChain *mocks.MockExchangeChain,
		mockNewsletterSender *mocks.MockNewsletterSender,
		mockEmailStorage *mocks.MockEmailStorage,
		mockExchangeResp *model.CurrencyRate,
		mockNewsletterResp []string,
	)

	cases := []struct {
		name                   string
		mockExchangeResponse   *model.CurrencyRate
		mockNewsletterResponse []string
		mockBehavior           mockBehavior
		isErrorExpected        bool
	}{
		{
			name: "OK",
			mockExchangeResponse: &model.CurrencyRate{
				Price: 123,
				CurrencyPair: model.CurrencyPair{
					BaseCurrency:  "BTC",
					QuoteCurrency: "UAH",
				},
			},
			mockNewsletterResponse: []string{"abc@test.com"},
			mockBehavior: func(
				mockExchangeChain *mocks.MockExchangeChain,
				mockNewsletterSender *mocks.MockNewsletterSender,
				mockEmailStorage *mocks.MockEmailStorage,
				mockExchangeResp *model.CurrencyRate,
				mockNewsletterResp []string,
			) {
				mockExchangeChain.EXPECT().GetCurrencyRate(&model.CurrencyPair{
					BaseCurrency:  "BTC",
					QuoteCurrency: "UAH",
				}).Return(mockExchangeResp, nil)
				mockEmailStorage.EXPECT().GetAllEmails().Return([]string{"abc@test.com"}, nil)
				mockNewsletterSender.EXPECT().MultipleSending([]string{"abc@test.com"}, &model.EmailMessage{
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
				mockNewsletterSender *mocks.MockNewsletterSender,
				mockEmailStorage *mocks.MockEmailStorage,
				mockExchangeResp *model.CurrencyRate,
				mockNewsletterResp []string,
			) {
				mockExchangeChain.EXPECT().GetCurrencyRate(&model.CurrencyPair{
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

			mockNewsletterSender := mocks.NewMockNewsletterSender(ctrl)
			mockExchangeChain := mocks.NewMockExchangeChain(ctrl)
			mockEmailRepository := mocks.NewMockEmailStorage(ctrl)

			newsletterTestService := NewNewsletterService(
				mockExchangeChain,
				mockEmailRepository,
				mockNewsletterSender,
				&model.CurrencyPair{
					BaseCurrency:  "BTC",
					QuoteCurrency: "UAH",
				})

			c.mockBehavior(
				mockExchangeChain,
				mockNewsletterSender,
				mockEmailRepository,
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

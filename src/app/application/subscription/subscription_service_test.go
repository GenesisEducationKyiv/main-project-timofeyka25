package subscription

import (
	"genesis-test/src/app/application/mocks"
	"genesis-test/src/app/customerror"
	"genesis-test/src/app/domain/model"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionService_Subscribe(t *testing.T) {
	type mockBehavior func(r *mocks.MockEmailStorage, s *model.Subscriber)

	cases := []struct {
		name          string
		subscriber    *model.Subscriber
		mockBehavior  mockBehavior
		expectedError error
	}{
		{
			name: "Subscribe successful",
			subscriber: &model.Subscriber{
				Email: "test@testexample.com",
			},

			mockBehavior: func(
				r *mocks.MockEmailStorage,
				s *model.Subscriber,
			) {
				r.EXPECT().AddEmail(s.Email).Return(nil)
			},
		},
		{
			name:          "Subscribe error (no data)",
			mockBehavior:  func(r *mocks.MockEmailStorage, s *model.Subscriber) {},
			expectedError: customerror.ErrNoDataProvided,
		},
		{
			name: "Subscribe error (already exists)",
			subscriber: &model.Subscriber{
				Email: "test@testexample.com",
			},
			mockBehavior: func(r *mocks.MockEmailStorage, s *model.Subscriber) {
				r.EXPECT().AddEmail(s.Email).Return(customerror.ErrAlreadyExists)
			},
			expectedError: customerror.ErrAlreadyExists,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockEmailStorage := mocks.NewMockEmailStorage(ctrl)
			c.mockBehavior(mockEmailStorage, c.subscriber)

			testSubscriptionService := NewSubscriptionService(mockEmailStorage)

			err := testSubscriptionService.Subscribe(c.subscriber)
			require.ErrorIs(t, err, c.expectedError)
		})
	}
}

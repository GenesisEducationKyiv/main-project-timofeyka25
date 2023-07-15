package utils

import (
	"genesis-test/src/app/customerror"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInsertToSortedSlice(t *testing.T) {
	cases := []struct {
		name          string
		inputSlice    []string
		inputToInsert string
		requireFunc   func(t *testing.T, slice []string, err error)
	}{
		{
			name:          "Insert successful",
			inputSlice:    []string{"abc@test.com", "qwe@test.com"},
			inputToInsert: "zxcvbn@test.com",
			requireFunc: func(t *testing.T, slice []string, err error) {
				require.NoError(t, err)
				require.IsIncreasing(t, slice)
			},
		},
		{
			name:          "Insert error (duplicate)",
			inputSlice:    []string{"abc@test.com", "qwe@test.com"},
			inputToInsert: "abc@test.com",
			requireFunc: func(t *testing.T, slice []string, err error) {
				require.ErrorIs(t, err, customerror.ErrAlreadyExists)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s, err := InsertToSortedSlice(c.inputSlice, c.inputToInsert)
			c.requireFunc(t, s, err)
		})
	}
}

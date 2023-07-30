package utils

import (
	"genesis-test/src/app/customerror"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteToCsv(t *testing.T) {
	cases := []struct {
		name          string
		filepath      string
		data          []string
		expectedError error
	}{
		{
			name:          "Write successful",
			filepath:      "test.csv",
			data:          []string{"abc@example.com", "qwe@example.com", "123@example.com"},
			expectedError: nil,
		},
		{
			name:          "Write error (no path)",
			filepath:      "",
			data:          []string{"abc@example.com", "qwe@example.com", "123@example.com"},
			expectedError: customerror.ErrInvalidInput,
		},
		{
			name:          "Write error (no data)",
			filepath:      "test.csv",
			data:          nil,
			expectedError: customerror.ErrInvalidInput,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := WriteToCsv(c.filepath, c.data)
			require.ErrorIs(t, err, c.expectedError)
			defer func(filepath string) {
				err = os.Remove(filepath)
			}(c.filepath)
		})
	}
}

func TestReadAllFromCsvToSlice(t *testing.T) {
	cases := []struct {
		name          string
		filepath      string
		expectedData  []string
		expectedError error
	}{
		{
			name:          "Read successful",
			filepath:      "../../storage/csv/data.csv",
			expectedData:  []string(nil),
			expectedError: nil,
		},
		{
			name:          "Read error (file not found)",
			filepath:      "nonexistent.csv",
			expectedData:  nil,
			expectedError: customerror.ErrFileNotFound,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			data, err := ReadAllFromCsvToSlice(c.filepath)

			require.Equal(t, c.expectedData, data)
			require.ErrorIs(t, err, c.expectedError)
		})
	}
}

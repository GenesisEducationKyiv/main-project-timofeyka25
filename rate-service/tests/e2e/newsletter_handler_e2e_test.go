package e2e

import (
	"bytes"
	"fmt"
	"genesis-test/src/app/utils"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestNewsletterHandler_Subscribe(t *testing.T) {
	loadEnvironment(t)
	if err := os.Setenv("STORAGE_FILE_PATH", "../../src/storage/csv/data_test.csv"); err != nil {
		t.Fatalf("failed to set STORAGE_FILE_PATH: %v", err)
	}
	if err := clearFile(os.Getenv("STORAGE_FILE_PATH")); err != nil {
		t.Fatalf("failed to clear file: %v", err)
	}

	cases := []struct {
		name               string
		expectedStatusCode int
		body               string
	}{
		{
			name:               "Subscribe successful",
			expectedStatusCode: fiber.StatusOK,
			body:               `{"email": "abc@example.com"}`,
		},
		{
			name:               "Invalid request body",
			expectedStatusCode: fiber.StatusBadRequest,
			body:               ``,
		},
		{
			name:               "Invalid email address",
			expectedStatusCode: fiber.StatusBadRequest,
			body:               `{"email": "invalid-email"}`,
		},
		{
			name:               "Already subscribed",
			expectedStatusCode: fiber.StatusConflict,
			body:               `{"email": "abc@example.com"}`,
		},
	}

	url := fmt.Sprintf("http://%s/api/subscribe", os.Getenv("SERVER_URL"))
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			client := &http.Client{}
			req, err := http.NewRequest("POST", url, bytes.NewBufferString(c.body)) //nolint:noctx
			req.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)

			res, err := client.Do(req) //nolint:bodyclose
			defer func(Body io.ReadCloser) {
				if err = Body.Close(); err != nil {
					t.Fatal(err)
				}
			}(res.Body)
			require.NoError(t, err)
			require.Equal(t, c.expectedStatusCode, res.StatusCode, "Unexpected status code: %d", res.StatusCode)
		})
	}
	if err := clearFile(os.Getenv("STORAGE_FILE_PATH")); err != nil {
		t.Fatalf("failed to clear file: %v", err)
	}
}

func TestNewsletterHandler_SendEmails(t *testing.T) {
	loadEnvironment(t)
	if err := os.Setenv("STORAGE_FILE_PATH", "../../src/storage/csv/data_test.csv"); err != nil {
		t.Fatalf("failed to set STORAGE_FILE_PATH: %v", err)
	}

	cases := []struct {
		name               string
		expectedStatusCode int
	}{
		{
			name:               "Send emails successful",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Send emails error (no subscribers)",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	if err := utils.WriteToCsv(os.Getenv("STORAGE_FILE_PATH"), []string{"test@test.com"}); err != nil {
		t.Fatalf("failed to write in csv: %v", err)
	}

	url := fmt.Sprintf("http://%s/api/sendEmails", os.Getenv("SERVER_URL"))
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			loadEnvironment(t)
			client := &http.Client{}
			req, err := http.NewRequest("POST", url, nil) //nolint:noctx
			require.NoError(t, err)

			res, err := client.Do(req) //nolint:bodyclose
			defer func(Body io.ReadCloser) {
				if err = Body.Close(); err != nil {
					t.Fatal(err)
				}
			}(res.Body)
			require.NoError(t, err)
			require.Equal(t, c.expectedStatusCode, res.StatusCode, "Unexpected status code: %d", res.StatusCode)
		})
		if err := clearFile(os.Getenv("STORAGE_FILE_PATH")); err != nil {
			t.Fatalf("failed to clear file: %v", err)
		}
	}
}

func clearFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0o666)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Fatal("failed to close file")
		}
	}(file)

	err = file.Truncate(0)
	return err
}

func loadEnvironment(t *testing.T) {
	if err := godotenv.Load("../../test.env"); err != nil {
		t.Fatal("Failed to load .env file")
	}
}

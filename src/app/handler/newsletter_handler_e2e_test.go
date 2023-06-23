package handler

import (
	"bytes"
	"genesis-test/src/app/repository"
	"genesis-test/src/app/service"
	"genesis-test/src/config"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func setupEnvironment(t *testing.T) {
	if err := os.Setenv("STORAGE_FILE_PATH", "../../storage/csv/data_test.csv"); err != nil {
		t.Fatal("Failed to set STORAGE_FILE_PATH")
	}
	if err := os.Setenv("SMTP_SERVER", "sandbox.smtp.mailtrap.io"); err != nil {
		t.Fatal("Failed to set SMTP_SERVER")
	}
	if err := os.Setenv("SMTP_PORT", "2525"); err != nil {
		t.Fatal("Failed to set SMTP_PORT")
	}
	if err := os.Setenv("SMTP_USERNAME", "baac2a76689b33"); err != nil {
		t.Fatal("Failed to set SMTP_USERNAME")
	}
	if err := os.Setenv("SMTP_PASSWORD", "3b9561ea1b84ff"); err != nil {
		t.Fatal("Failed to set SMTP_PASSWORD")
	}
}

func loadEnvironment(t *testing.T) {
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Fatal("Failed to load .env file")
	}
}

func TestNewsletterHandler_Subscribe(t *testing.T) {
	setupEnvironment(t)
	loadEnvironment(t)
	repos := repository.NewRepositories()
	services := service.NewServices(repos)
	newsletterHandler := NewNewsletterHandler(services)

	app := fiber.New(config.FiberConfig())
	api := app.Group("/api")
	api.Post("/subscribe", newsletterHandler.Subscribe)

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
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/subscribe", bytes.NewBufferString(c.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req) //nolint:bodyclose
			defer func(Body io.ReadCloser) {
				if err = Body.Close(); err != nil {
					t.Fatal(err)
				}
			}(resp.Body)

			require.NoError(t, err)
			require.Equal(t, c.expectedStatusCode, resp.StatusCode)
		})
	}
}

func TestNewsletterHandler_SendEmails(t *testing.T) {
	cases := []struct {
		name               string
		filepath           string
		expectedStatusCode int
	}{
		{
			name:               "Send emails successful",
			filepath:           "../../storage/csv/data_test.csv",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Send emails error (no subscribers)",
			filepath:           "../../storage/csv/data_test_empty.csv",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			setupEnvironment(t)
			if err := os.Setenv("STORAGE_FILE_PATH", c.filepath); err != nil {
				t.Fatal("Failed to set STORAGE_FILE_PATH")
			}
			loadEnvironment(t)

			repos := repository.NewRepositories()
			services := service.NewServices(repos)
			newsletterHandler := NewNewsletterHandler(services)

			app := fiber.New(config.FiberConfig())
			api := app.Group("/api")
			api.Post("/sendEmails", newsletterHandler.SendEmails)

			req := httptest.NewRequest(http.MethodPost, "/api/sendEmails", nil)

			resp, err := app.Test(req, 5000) //nolint:bodyclose
			defer func(Body io.ReadCloser) {
				if err = Body.Close(); err != nil {
					t.Fatal(err)
				}
			}(resp.Body)

			require.NoError(t, err)
			require.Equal(t, c.expectedStatusCode, resp.StatusCode)
			if _, err := os.Create(os.Getenv("STORAGE_FILE_PATH")); err != nil {
				t.Fatal("error cleaning test file")
			}
		})
	}
}

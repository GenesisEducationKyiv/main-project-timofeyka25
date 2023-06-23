package handler

import (
	"genesis-test/src/app/repository"
	"genesis-test/src/app/service"
	"genesis-test/src/config"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gofiber/fiber/v2"
)

func TestExchangeHandler_GetCurrencyRate(t *testing.T) {
	cases := []struct {
		name               string
		apiURL             string
		expectedStatusCode int
	}{
		{
			name:               "Get rate successful",
			apiURL:             "https://api.coinbase.com/v2/prices/%s-%s/spot",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Get rate error (invalid url)",
			apiURL:             "",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := os.Setenv("CRYPTO_API_FORMAT_URL", c.apiURL); err != nil {
				t.Fatal("Failed to set CRYPTO_API_FORMAT_URL")
			}
			loadEnvironment(t)
			repos := repository.NewRepositories()
			services := service.NewServices(repos)
			exchangeHandler := NewExchangeHandler(services)

			app := fiber.New(config.FiberConfig())
			api := app.Group("/api")
			api.Get("/rate", exchangeHandler.GetCurrencyRate)

			req := httptest.NewRequest(http.MethodGet, "/api/rate", nil)
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

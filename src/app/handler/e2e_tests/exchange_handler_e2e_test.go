package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExchangeHandler_GetCurrencyRate(t *testing.T) {
	expectedStatusCode := http.StatusOK
	loadEnvironment(t)
	url := fmt.Sprintf("http://%s/api/rate", os.Getenv("SERVER_URL"))

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil) //nolint:noctx
	require.NoError(t, err)

	res, err := client.Do(req) //nolint:bodyclose
	require.NoError(t, err)
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			t.Fatal(err)
		}
	}(res.Body)
	require.Equal(t, expectedStatusCode, res.StatusCode, "Unexpected status code: %d", res.StatusCode)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.NotEmpty(nil, body, "Response body is empty")
	t.Log(string(body))
}

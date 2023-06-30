package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"genesis-test/src/app/domain"
	"genesis-test/src/config"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type exchangeRepository struct{}

func NewExchangeRepository() domain.ExchangeRepository {
	return &exchangeRepository{}
}

func (e exchangeRepository) GetCurrencyRate(base string, quoted string) (*domain.CurrencyRate, error) {
	cfg := config.Get()
	url := fmt.Sprintf(cfg.CryptoAPIFormatURL, base, quoted)

	client := http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HTTP request")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make HTTP request")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	var data struct {
		Rate domain.CurrencyRate `json:"data"`
	}
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal JSON")
	}

	data.Rate.Price = strings.Split(data.Rate.Price, ".")[0]
	return &data.Rate, nil
}

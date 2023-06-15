package repository

import (
	"encoding/json"
	"fmt"
	"genesis-test/src/app/domain"
	"genesis-test/src/config"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strings"
)

type exchangeRepository struct{}

func NewExchangeRepository() domain.ExchangeRepository {
	return &exchangeRepository{}
}

func (e exchangeRepository) GetCurrencyRate(base string, quoted string) (*domain.CurrencyRate, error) {
	cfg := config.Get()
	url := fmt.Sprintf(cfg.CryptoApiFormatUrl, base, quoted)

	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make HTTP request")
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
	}(resp.Body)

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

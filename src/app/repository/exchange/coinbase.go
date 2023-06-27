package exchange

import (
	"context"
	"encoding/json"
	"fmt"
	"genesis-test/src/app/domain"
	"genesis-test/src/app/service"
	"genesis-test/src/config"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

type coinbaseCurrencyRate struct {
	Amount        string `json:"amount"`
	BaseCurrency  string `json:"base"`
	QuoteCurrency string `json:"currency"`
}

type coinbaseExchangerResponse struct {
	coinbaseCurrencyRate `json:"data"`
}

type exchangeCoinbaseRepository struct {
	CoinbaseEndpoint string
}

func NewExchangeCoinbaseRepository(coinbaseEndpoint string) service.ExchangeRepository {
	return &exchangeCoinbaseRepository{
		CoinbaseEndpoint: coinbaseEndpoint,
	}
}

func (e exchangeCoinbaseRepository) GetCurrencyRate(pair *domain.CurrencyPair) (*domain.CurrencyRate, error) {
	cfg := config.Get()
	url := fmt.Sprintf(cfg.CryptoAPIFormatURL, pair.BaseCurrency, pair.QuoteCurrency)

	client := http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create HTTP request")
	}

	resp, err := client.Do(req) //nolint:bodyclose
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Fatal(err)
		}
	}(resp.Body)

	if err != nil {
		return nil, errors.Wrap(err, "failed to make HTTP request")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read response body")
	}

	rate := new(coinbaseExchangerResponse)

	if err = json.Unmarshal(body, &rate); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal JSON")
	}

	return rate.toDefaultRate()
}

func (c *coinbaseCurrencyRate) toDefaultRate() (*domain.CurrencyRate, error) {
	bitSize := 64
	floatPrice, err := strconv.ParseFloat(c.Amount, bitSize)
	if err != nil {
		return nil, err
	}
	return &domain.CurrencyRate{
		Price: floatPrice,
		CurrencyPair: domain.CurrencyPair{
			BaseCurrency:  c.BaseCurrency,
			QuoteCurrency: c.QuoteCurrency,
		},
	}, nil
}

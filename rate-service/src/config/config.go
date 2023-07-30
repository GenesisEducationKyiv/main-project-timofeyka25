package config

import (
	"os"
	"sync"
)

type Config struct {
	ServerURL         string
	ServerReadTimeout string
	BaseCurrency      string
	QuoteCurrency     string
	CoinbaseURL       string
	BinanceURL        string
	BTCTradeUAURL     string
	CoingeckoURL      string
	StorageFile       string
	LogPath           string
	SMTPServer        string
	SMTPPort          string
	SMTPUsername      string
	SMTPPassword      string
}

func Get() *Config {
	var cfg Config
	var once sync.Once
	once.Do(func() {
		cfg = Config{
			ServerURL:         os.Getenv("SERVER_URL"),
			ServerReadTimeout: os.Getenv("SERVER_READ_TIMEOUT"),
			CoinbaseURL:       os.Getenv("COINBASE_URL"),
			BinanceURL:        os.Getenv("BINANCE_URL"),
			BTCTradeUAURL:     os.Getenv("BTC_TRADE_UA_URL"),
			CoingeckoURL:      os.Getenv("COINGECKO_URL"),
			BaseCurrency:      os.Getenv("BASE_CURRENCY"),
			QuoteCurrency:     os.Getenv("QUOTED_CURRENCY"),
			StorageFile:       os.Getenv("STORAGE_FILE_PATH"),
			LogPath:           os.Getenv("LOG_PATH"),
			SMTPServer:        os.Getenv("SMTP_SERVER"),
			SMTPPort:          os.Getenv("SMTP_PORT"),
			SMTPUsername:      os.Getenv("SMTP_USERNAME"),
			SMTPPassword:      os.Getenv("SMTP_PASSWORD"),
		}
	})
	return &cfg
}

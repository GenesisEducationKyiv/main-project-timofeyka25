package config

import (
	"os"
	"sync"
)

type Config struct {
	ServerURL          string
	ServerReadTimeout  string
	BaseCurrency       string
	QuoteCurrency      string
	CryptoApiFormatUrl string
	StorageFile        string
	SmtpServer         string
	SmtpPort           string
	SmtpUsername       string
	SmtpPassword       string
}

var (
	cfg  Config
	once sync.Once
)

func Get() *Config {
	once.Do(func() {
		cfg = Config{
			ServerURL:          os.Getenv("SERVER_URL"),
			ServerReadTimeout:  os.Getenv("SERVER_READ_TIMEOUT"),
			CryptoApiFormatUrl: os.Getenv("CRYPTO_API_FORMAT_URL"),
			BaseCurrency:       os.Getenv("BASE_CURRENCY"),
			QuoteCurrency:      os.Getenv("QUOTED_CURRENCY"),
			StorageFile:        os.Getenv("STORAGE_FILE_PATH"),
			SmtpServer:         os.Getenv("SMTP_SERVER"),
			SmtpPort:           os.Getenv("SMTP_PORT"),
			SmtpUsername:       os.Getenv("SMTP_USERNAME"),
			SmtpPassword:       os.Getenv("SMTP_PASSWORD"),
		}
	})
	return &cfg
}

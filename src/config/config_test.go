package config

import (
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		t.Error("Error loading .env file")
	}
	cfg := Get()
	require.NotEmpty(t, cfg.ServerURL)
	require.NotEmpty(t, cfg.ServerReadTimeout)
	require.NotEmpty(t, cfg.BaseCurrency)
	require.NotEmpty(t, cfg.QuoteCurrency)
	require.NotEmpty(t, cfg.CryptoAPIFormatURL)
	require.NotEmpty(t, cfg.StorageFile)
	require.NotEmpty(t, cfg.SMTPServer)
	require.NotEmpty(t, cfg.SMTPPort)
	require.NotEmpty(t, cfg.SMTPUsername)
	require.NotEmpty(t, cfg.SMTPPassword)
}

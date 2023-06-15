package repository

import (
	"genesis-test/src/app/domain"
	"genesis-test/src/config"
)

type Repositories struct {
	Exchange   domain.ExchangeRepository
	Newsletter domain.NewsletterRepository
}

func NewRepositories() *Repositories {
	cfg := config.Get()

	return &Repositories{
		Exchange: NewExchangeRepository(),
		Newsletter: NewNewsletterRepository(
			cfg.SmtpServer,
			cfg.SmtpPort,
			cfg.SmtpUsername,
			cfg.SmtpPassword,
		),
	}
}

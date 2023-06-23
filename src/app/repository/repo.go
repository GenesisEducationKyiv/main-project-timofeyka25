package repository

import (
	"genesis-test/src/app/domain"
	"genesis-test/src/config"
	mailer "genesis-test/src/pkg/mailer"
)

type Repositories struct {
	Exchange   domain.ExchangeRepository
	Newsletter domain.NewsletterRepository
}

func NewRepositories() *Repositories {
	cfg := config.Get()
	smtpMailer := mailer.NewSMTPMailer(
		cfg.SMTPServer,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword)

	return &Repositories{
		Exchange: NewExchangeRepository(),
		Newsletter: NewNewsletterRepository(
			smtpMailer,
			cfg.StorageFile,
		),
	}
}

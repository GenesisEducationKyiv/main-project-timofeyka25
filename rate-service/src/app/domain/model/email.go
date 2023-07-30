package model

//go:generate mockgen -destination=mocks/mock_newsletter.go genesis-test/src/app/domain NewsletterRepository,NewsletterService

type EmailMessage struct {
	Subject string
	Body    string
}

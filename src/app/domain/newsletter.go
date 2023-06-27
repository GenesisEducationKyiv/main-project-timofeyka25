package domain

//go:generate mockgen -destination=mocks/mock_newsletter.go genesis-test/src/app/domain NewsletterRepository,NewsletterService

type Subscriber struct {
	Email string `json:"email" validate:"required,email"`
}

type EmailMessage struct {
	Subject string
	Body    string
}

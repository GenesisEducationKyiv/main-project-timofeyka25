package domain

//go:generate mockgen -destination=mocks/mock_newsletter.go genesis-test/src/app/domain NewsletterRepository,NewsletterService

type Subscriber struct {
	Email string `json:"email" validate:"required,email"`
}

type NewsletterService interface {
	SendEmails() ([]string, error)
	Subscribe(subscriber *Subscriber) error
}

type NewsletterRepository interface {
	GetSubscribedEmails() ([]string, error)
	SendToSubscribedEmails(body string) ([]string, error)
	SendEmail(to, body string) error
	AddNewEmail(emails []string, emailToInsert string) error
}

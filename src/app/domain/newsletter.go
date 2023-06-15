package domain

type Subscriber struct {
	Email string `json:"email" validate:"required,email"`
}

type NewsletterService interface {
	SendEmails() ([]string, error)
	Subscribe(recipient *Subscriber) error
}

type NewsletterRepository interface {
	GetSubscribedEmails() ([]string, error)
	SendToSubscribedEmails(subject, body string) ([]string, error)
	SendEmail(to, subject, body string) error
	AddNewEmail(emails []string, emailToInsert string) error
}

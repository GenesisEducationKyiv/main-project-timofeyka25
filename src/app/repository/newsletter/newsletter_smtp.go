package newsletter

import (
	"fmt"
	"genesis-test/src/app/domain"
	"genesis-test/src/app/service"
	"genesis-test/src/pkg/mailer"
	"log"

	"github.com/pkg/errors"
)

type newsletterSMTPRepository struct {
	mailer *mailer.SMTPMailer
}

func NewNewsletterRepository(mailer *mailer.SMTPMailer) service.NewsletterRepository {
	return &newsletterSMTPRepository{
		mailer: mailer,
	}
}

func (r *newsletterSMTPRepository) MultipleSending(subscribers []string, msg *domain.EmailMessage) ([]string, error) {
	var unsent []string

	body := fmt.Sprintf("Subject: %s\r\n%s", msg.Subject, msg.Body)

	for _, subscriber := range subscribers {
		if err := r.sendEmail(subscriber, fmt.Sprintf("To: %s\r\n", subscriber)+body); err != nil {
			log.Printf("error sending: %v\n", err)
			unsent = append(unsent, subscriber)
		} else {
			log.Printf("email sent to: %s\n", subscriber)
		}
	}
	return unsent, nil
}

func (r *newsletterSMTPRepository) sendEmail(to, body string) error {
	err := r.mailer.SendEmail(to, body)
	if err != nil {
		return errors.Wrap(err, "send email")
	}

	return nil
}

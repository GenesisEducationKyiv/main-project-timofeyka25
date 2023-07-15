package newsletter

import (
	"fmt"
	"genesis-test/src/app/application"
	"genesis-test/src/app/domain/model"
	"genesis-test/src/pkg/mailer"
	"log"

	"github.com/pkg/errors"
)

type newsletterSMTPSender struct {
	mailer *mailer.SMTPMailer
}

func NewNewsletterSender(mailer *mailer.SMTPMailer) application.NewsletterSender {
	return &newsletterSMTPSender{
		mailer: mailer,
	}
}

func (r *newsletterSMTPSender) MultipleSending(subscribers []string, msg *model.EmailMessage) ([]string, error) {
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

func (r *newsletterSMTPSender) sendEmail(to, body string) error {
	err := r.mailer.SendEmail(to, body)
	if err != nil {
		return errors.Wrap(err, "send email")
	}

	return nil
}

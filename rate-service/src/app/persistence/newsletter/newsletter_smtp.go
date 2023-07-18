package newsletter

import (
	"fmt"
	"genesis-test/src/app/domain/model"
	"genesis-test/src/pkg/mailer"

	"github.com/pkg/errors"
)

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Error(msg string)
}

type NewsletterSMTPSender struct {
	mailer *mailer.SMTPMailer
	logger Logger
}

func NewNewsletterSender(
	mailer *mailer.SMTPMailer,
	logger Logger,
) *NewsletterSMTPSender {
	return &NewsletterSMTPSender{
		mailer: mailer,
		logger: logger,
	}
}

func (r *NewsletterSMTPSender) MultipleSending(subscribers []string, msg *model.EmailMessage) ([]string, error) {
	var unsent []string

	body := fmt.Sprintf("Subject: %s\r\n%s", msg.Subject, msg.Body)

	for _, subscriber := range subscribers {
		if err := r.sendEmail(subscriber, fmt.Sprintf("To: %s\r\n", subscriber)+body); err != nil {
			r.logger.Info(fmt.Sprintf("error sending: %v\n", err))
			unsent = append(unsent, subscriber)
		} else {
			r.logger.Info(fmt.Sprintf("email sent to: %s\n", subscriber))
		}
	}
	return unsent, nil
}

func (r *NewsletterSMTPSender) sendEmail(to, body string) error {
	err := r.mailer.SendEmail(to, body)
	if err != nil {
		return errors.Wrap(err, "send email")
	}

	return nil
}

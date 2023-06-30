package repository

import (
	"genesis-test/src/app/customerror"
	"genesis-test/src/app/domain"
	"genesis-test/src/app/utils"
	mailer "genesis-test/src/pkg/mailer"
	"log"

	"github.com/pkg/errors"
)

type newsletterRepository struct {
	mailer      *mailer.SMTPMailer
	storageFile string
}

func NewNewsletterRepository(mailer *mailer.SMTPMailer, storageFile string) domain.NewsletterRepository {
	return &newsletterRepository{
		mailer:      mailer,
		storageFile: storageFile,
	}
}

func (r newsletterRepository) GetSubscribedEmails() ([]string, error) {
	subscribed, err := utils.ReadAllFromCsvToSlice(r.storageFile)
	if err != nil {
		return nil, errors.Wrap(err, "read csv file")
	}
	if len(subscribed) < 1 {
		return nil, customerror.ErrNoSubscribers
	}

	return subscribed, nil
}

func (r newsletterRepository) SendToSubscribedEmails(body string) ([]string, error) {
	subscribed, err := r.GetSubscribedEmails()
	if err != nil {
		return nil, errors.Wrap(err, "get emails")
	}

	unsent := make([]string, 0)

	for _, email := range subscribed {
		if err = r.SendEmail(email, body); err != nil {
			log.Printf("Sending error: %v\n", err)
			unsent = append(unsent, email)
		} else {
			log.Printf("Message sent to %s\n", email)
		}
	}

	return unsent, nil
}

func (r newsletterRepository) SendEmail(to, body string) error {
	msg := "To: " + to + "\r\n" +
		"Subject: Exchange Currency Newsletter" + "\r\n" +
		"\r\n" + body

	err := r.mailer.SendEmail(to, msg)
	if err != nil {
		return errors.Wrap(err, "send email")
	}

	return nil
}

func (r newsletterRepository) AddNewEmail(emails []string, emailToInsert string) error {
	emails, err := utils.InsertToSortedSlice(emails, emailToInsert)
	if err != nil {
		return errors.Wrap(err, "insert email")
	}

	return utils.WriteToCsv(r.storageFile, emails)
}

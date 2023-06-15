package repository

import (
	"genesis-test/src/app/domain"
	"genesis-test/src/app/utils"
	"genesis-test/src/config"
	"github.com/pkg/errors"
	"log"
	"net/smtp"
)

type newsletterRepository struct {
	smtpServer string
	smtpPort   string
	username   string
	password   string
}

func NewNewsletterRepository(
	smtpServer,
	smtpPort,
	username,
	password string) domain.NewsletterRepository {
	return &newsletterRepository{
		smtpServer: smtpServer,
		smtpPort:   smtpPort,
		username:   username,
		password:   password,
	}
}

func (r newsletterRepository) GetSubscribedEmails() ([]string, error) {
	cfg := config.Get()

	subscribed, err := utils.ReadAllFromCsvToSlice(cfg.StorageFile)
	if err != nil {
		return nil, err
	}

	return subscribed, nil
}

func (r newsletterRepository) SendToSubscribedEmails(subject, body string) ([]string, error) {
	subscribed, err := r.GetSubscribedEmails()
	if err != nil {
		return nil, err
	}

	unsent := make([]string, 0)

	for _, email := range subscribed {
		if err = r.SendEmail(email, subject, body); err != nil {
			log.Printf("Sending error: %v\n", err)
			unsent = append(unsent, email)
		} else {
			log.Printf("Message sent to %s\n", email)
		}
	}

	return unsent, nil
}

func (r newsletterRepository) SendEmail(to, subject, body string) error {
	// Construct the email message
	msg := "To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body

	// Connect to the SMTP server
	auth := smtp.PlainAuth("", r.username, r.password, r.smtpServer)
	address := r.smtpServer + ":" + r.smtpPort
	err := smtp.SendMail(address, auth, r.username, []string{to}, []byte(msg))
	if err != nil {
		return errors.Wrap(err, "sending failed")
	}

	return nil
}

func (r newsletterRepository) AddNewEmail(emails []string, emailToInsert string) error {
	cfg := config.Get()

	emails, err := utils.InsertToSortedSlice(emails, emailToInsert)
	if err != nil {
		return err
	}

	return utils.WriteToCsv(cfg.StorageFile, emails)
}

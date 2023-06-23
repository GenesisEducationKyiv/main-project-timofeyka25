package mailer

import (
	"log"
	"net/smtp"

	"github.com/pkg/errors"
)

type SMTPMailer struct {
	smtpServer   string
	smtpPort     string
	smtpUsername string
	smtpPassword string
}

func NewSMTPMailer(
	smtpServer,
	smtpPort,
	smtpUsername,
	smtpPassword string,
) *SMTPMailer {
	return &SMTPMailer{
		smtpServer:   smtpServer,
		smtpPort:     smtpPort,
		smtpUsername: smtpUsername,
		smtpPassword: smtpPassword,
	}
}

func (m *SMTPMailer) SendEmail(toEmail, message string) error {
	auth := smtp.PlainAuth("", m.smtpUsername, m.smtpPassword, m.smtpServer)
	address := m.smtpServer + ":" + m.smtpPort
	err := smtp.SendMail(address, auth, m.smtpUsername, []string{toEmail}, []byte(message))
	log.Println(err, auth)
	if err != nil {
		return errors.Wrap(err, "send email")
	}
	return nil
}

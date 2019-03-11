package mail

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strings"

	gomail "github.com/go-mail/mail"
)

// SMTPOptions contains the options for establishing a SMTP connection
// to a mail server

var (
	MailErrorConnection       = "Something happened on our end, please contact the support"
	MailErrorUnknownRecipient = "The specified E-Mail address seems to contain a typo"
)

type SMTPOptions struct {
	Host          string
	Port          int
	User          string
	Pass          string
	From          string
	AllowInsecure bool
}

// SMTPMailer is a Emailer that uses an SMTP Server to send mails
type SMTPMailer struct {
	Options *SMTPOptions
}

// Email satisfies the Emailer interface
func (m SMTPMailer) Email(recepient, subject, html, text string) error {

	mail := gomail.NewDialer(m.Options.Host, m.Options.Port, m.Options.User, m.Options.Pass)

	if m.Options.AllowInsecure {
		// Seems that unencrypted connection are not really supported
		// https://github.com/go-mail/mail/issues/52
		mail.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		mail.SSL = false
	}

	message := gomail.NewMessage()

	message.SetHeader("From", m.Options.From)
	message.SetAddressHeader("To", recepient, recepient)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", text)

	err := mail.DialAndSend(message)
	message.Reset()

	if strings.Contains(err.Error(), "unencrypted connection") {

		return errors.New(MailErrorConnection)
	} else if strings.Contains(err.Error(), "Recipient address rejected") {
		return errors.New(MailErrorUnknownRecipient)
	} else {
		fmt.Println("Error seding email message")
		fmt.Println(recepient)
		fmt.Println(err)
		return errors.New(MailErrorConnection)
	}
}

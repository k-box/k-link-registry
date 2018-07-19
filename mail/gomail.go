package mail

import (
	"github.com/go-gomail/gomail"
)

// SMTPOptions contains the options for establishing a SMTP connection
// to a mail server
type SMTPOptions struct {
	Host string
	Port int
	User string
	Pass string
	From string
}

// SMTPMailer is a Emailer that uses an SMTP Server to send mails
type SMTPMailer struct {
	Options *SMTPOptions
}

// Email satisfies the Emailer interface
func (m SMTPMailer) Email(recepient, subject, html, text string) error {
	mail := gomail.NewDialer(m.Options.Host, m.Options.Port, m.Options.User, m.Options.Pass)
	message := gomail.NewMessage()

	message.SetHeader("From", m.Options.From)
	message.SetAddressHeader("To", recepient, recepient)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", text)

	err := mail.DialAndSend(message)
	//message.Reset() // optional here: prepare for re-use

	return err
}

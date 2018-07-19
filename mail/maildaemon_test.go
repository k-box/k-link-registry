package mail

import (
	"log"
	"testing"
	"time"

	"github.com/go-gomail/gomail"
	"github.com/pkg/errors"
)

type MailDaemon struct {
	From    string
	running bool
	dialer  *gomail.Dialer
	queue   chan *gomail.Message
}

func NewMailDaemon() *MailDaemon {
	return &MailDaemon{}
}

func (md *MailDaemon) Init(host string, port int, from, user, pass string) error {
	md.From = from
	dialer := gomail.NewDialer(host, port, user, pass)

	// test connection
	conn, err := dialer.Dial()
	if err != nil {
		return errors.Wrap(err, "Could not establish connection to mail server")
	}
	conn.Close()

	md.dialer = dialer
	return nil
}

func (md *MailDaemon) Run() {
	// do not run more than once
	if md.running == true {
		return
	}
	// reset running state, if this function quits
	defer func() {
		md.running = false
	}()

	var s gomail.SendCloser
	var err error
	open := false
	for {
		select {
		case m, ok := <-md.queue:
			if !ok {
				return
			}
			if !open {
				if s, err = md.dialer.Dial(); err != nil {
					panic(err)
				}
				open = true
			}
			if err := gomail.Send(s, m); err != nil {
				log.Print(err)
			}
		// Close the connection to the SMTP server if no email was sent in
		// the last 30 seconds.
		case <-time.After(30 * time.Second):
			if open {
				open = false
				if err := s.Close(); err != nil {
					panic(err)
				}
			}
		}
	}
}

func (md *MailDaemon) Email(recepient, subject, html, text string) error {
	message := gomail.NewMessage()

	message.SetHeader("From", md.From)
	message.SetAddressHeader("To", recepient, recepient)
	message.SetHeader("Subject", subject)
	message.SetBody("text/plain", text)

	if html != "" {
		message.AddAlternative("text/html", html)
	}

	if !md.running {
		// send mail directly, since no queue is running
		err := md.dialer.DialAndSend(message)
		if err != nil {
			err = errors.Wrap(err, "Mail: Failed sending directly")
		}
		return err
	}

	// put the message on the queue.
	md.queue <- message
	return nil
}

func TestMailSending(t *testing.T) {
	md := NewMailDaemon()

	err := md.Init("example.com", 587, "test@example.com", "test", "password")
	if err != nil {
		log.Fatalf("Error dialing to server: %s", err)
	}

	go md.Run()

	md.Email("paul@zom.bi", "Muh mail", "", "test content")

	time.Sleep(4 * time.Second)
}

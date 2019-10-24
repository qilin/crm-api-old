package invites

import (
	"crypto/tls"
	"log"

	"gopkg.in/gomail.v2"
)

type Mailer interface {
	Send(to, subject, body string) error
}

type mailerImpl struct {
	replyTo string
	from    string
	dialer  *gomail.Dialer
}

func NewMailer(config MailConfig) (mailer Mailer) {
	dialer := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: config.InsecureSkipVerify}
	mailer = &mailerImpl{config.ReplyTo, config.From, dialer}
	return
}

func (mailer *mailerImpl) Send(to, subject, body string) (err error) {

	m := gomail.NewMessage()
	m.SetHeader("From", mailer.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	if mailer.replyTo != "" {
		m.SetHeader("ReplyTo", mailer.replyTo)
	}
	m.SetBody("text/html", body)

	err = mailer.dialer.DialAndSend(m)
	if err != nil {
		log.Println(err)
	}

	return
}

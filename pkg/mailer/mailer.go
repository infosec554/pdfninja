package mailer

import (
	"fmt"
	"net/smtp"
)

type Mailer struct {
	Host       string
	Port       string
	User       string
	Pass       string
	SenderName string
}

func New(host, port, user, pass, sender string) *Mailer {
	return &Mailer{
		Host:       host,
		Port:       port,
		User:       user,
		Pass:       pass,
		SenderName: sender,
	}
}

func (m *Mailer) Send(to string, subject string, body string) error {
	addr := fmt.Sprintf("%s:%s", m.Host, m.Port)

	msg := "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/html; charset=\"UTF-8\"\r\n"
	msg += fmt.Sprintf("From: %s <%s>\r\n", m.SenderName, m.User)
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n\r\n", subject)
	msg += body

	auth := smtp.PlainAuth("", m.User, m.Pass, m.Host)
	return smtp.SendMail(addr, auth, m.User, []string{to}, []byte(msg))
}

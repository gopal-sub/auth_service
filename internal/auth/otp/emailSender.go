package otp

import (
	"fmt"
	"net/smtp"
)

type EmailSender interface{
	SendEmail(to, subject, body string) error
}


type SMTPEmailSender struct {
	host string
	port string
	username string
	password string
}

func NewSMTPEmailSender(host string, port string, username string, password string) *SMTPEmailSender{
	return &SMTPEmailSender{
		host: host,
		port: port,
		username: username,
		password: password,
	}
}


func (s *SMTPEmailSender) SendEmail(to, subject, body, from string) error{
	auth := smtp.PlainAuth(
		"",
		s.username,
		s.password,
		s.host,
	)

	message := fmt.Sprintf("Subject: %s\r\n\r\n%s", subject, body)

	err := smtp.SendMail(
		s.host+":"+s.port,
		auth,
		s.username,
		[]string{to},
		[]byte(message),
	)
	return err
}


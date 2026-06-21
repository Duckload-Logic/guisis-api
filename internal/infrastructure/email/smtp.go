package email

import (
	"context"
	"fmt"
	"net"
	"net/smtp"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
)

type SMTPMailer struct {
	host     string
	port     int
	username string
	password string
}

func NewSMTPMailer(
	host string,
	port int,
	username, password string,
) *SMTPMailer {
	return &SMTPMailer{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (s *SMTPMailer) Send(
	ctx context.Context,
	email audit.EmailEntry,
) error {
	body, err := renderTemplate(email)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	from := constants.FromEmail()

	toHeader := ""
	for i, to := range email.To {
		if i > 0 {
			toHeader += ", "
		}
		toHeader += to
	}

	msg := fmt.Sprintf("To: %s\r\n", toHeader) +
		fmt.Sprintf("From: %s\r\n", from) +
		fmt.Sprintf("Subject: %s\r\n", email.Subject) +
		"MIME-version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body

	var auth smtp.Auth
	if s.username != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	}

	err = smtp.SendMail(addr, auth, from, email.To, []byte(msg))
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return fmt.Errorf("[SMTP] Connection timeout: %w", err)
		}
		return fmt.Errorf("[SMTP] Protocol Error: %w", err)
	}

	return nil
}

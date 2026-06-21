package email

import (
	"context"
	"fmt"
	"net"
	"net/smtp"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
)

type MailPit struct {
	host string
	port int
}

func NewMailPit(host string, port int) (*MailPit, error) {
	return &MailPit{
		host: host,
		port: port,
	}, nil
}

func (m *MailPit) Send(
	ctx context.Context,
	email audit.EmailEntry,
) error {
	body, err := renderTemplate(email)
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%d", m.host, m.port)
	from := constants.FromEmail()

	// Construct the email message in standard SMTP format (using \r\n)
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

	err = smtp.SendMail(addr, nil, from, email.To, []byte(msg))
	if err != nil {
		// Check for Network Timeouts/Connection errors
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return fmt.Errorf("[MailPit] Connection timeout: %w", err)
		}

		// Check for SMTP Protocol Errors (e.g., 550 Invalid Recipient)
		return fmt.Errorf("[MailPit] SMTP Protocol Error: %w", err)
	}

	return nil
}

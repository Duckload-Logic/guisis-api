package email

import (
	"context"
	"fmt"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGrid struct {
	client *sendgrid.Client
}

func NewSendGrid(apiKey string) *SendGrid {
	return &SendGrid{
		client: sendgrid.NewSendClient(apiKey),
	}
}

func (s *SendGrid) Send(
	ctx context.Context,
	email audit.EmailEntry,
) error {
	body, err := renderTemplate(email)
	if err != nil {
		return err
	}

	from := constants.FromEmail()
	fromAddress := mail.NewEmail("GuiSIS", from)
	m := mail.NewV3Mail()
	m.SetFrom(fromAddress)
	m.Subject = email.Subject

	p := mail.NewPersonalization()
	for _, to := range email.To {
		p.AddTos(mail.NewEmail("", to))
	}
	m.AddPersonalizations(p)

	content := mail.NewContent("text/html", body)
	m.AddContent(content)

	response, err := s.client.Send(m)
	if err != nil {
		return fmt.Errorf("[SendGrid] Network Error: %w", err)
	}

	if response.StatusCode >= 400 {
		return fmt.Errorf(
			"[SendGrid] API Error (Status %d): %s",
			response.StatusCode,
			response.Body,
		)
	}

	return nil
}

func (s *SendGrid) SendOTP(ctx context.Context, to, otp string) error {
	return s.Send(ctx, audit.EmailEntry{
		To:      []string{to},
		Subject: "Your Verification Code",
		Body:    fmt.Sprintf("<h1>Verification Code</h1><p>Your code is: <b>%s</b></p>", otp),
	})
}

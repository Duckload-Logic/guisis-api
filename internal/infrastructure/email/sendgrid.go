package email

import (
	"context"
	"fmt"

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

func (s *SendGrid) SendEmail(
	ctx context.Context,
	to, subject, body string,
) (bool, error) {
	from := constants.FromEmail()
	fromAddress := mail.NewEmail("GuiSIS", from)
	toAddress := mail.NewEmail("", to)

	message := mail.NewSingleEmail(fromAddress, subject, toAddress, "", body)

	response, err := s.client.Send(message)
	if err != nil {
		return false, fmt.Errorf("[SendGrid] Network Error: %w", err)
	}

	if response.StatusCode >= 400 {
		return false, fmt.Errorf(
			"[SendGrid] API Error (Status %d): %s",
			response.StatusCode,
			response.Body,
		)
	}

	return true, nil
}

func (s *SendGrid) SendOTP(ctx context.Context, to, otp string) (bool, error) {
	return s.SendEmail(ctx, to, "Verification Code", OTP_TEMPLATE(otp))
}

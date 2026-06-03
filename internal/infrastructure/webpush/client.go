package webpush

import (
	"context"
	"fmt"

	webpushlib "github.com/SherClockHolmes/webpush-go"
)

type Client struct {
	publicKey  string
	privateKey string
	email      string
}

func NewClient(pub, priv, email string) *Client {
	return &Client{
		publicKey:  pub,
		privateKey: priv,
		email:      email,
	}
}

// Send sends a web push notification to a subscription.
// Returns the HTTP response status code, or an error.
func (c *Client) Send(
	ctx context.Context,
	endpoint string,
	p256dh string,
	auth string,
	payload []byte,
) (int, error) {
	sub := webpushlib.Subscription{
		Endpoint: endpoint,
		Keys: webpushlib.Keys{
			P256dh: p256dh,
			Auth:   auth,
		},
	}

	opts := &webpushlib.Options{
		Subscriber:      c.email,
		VAPIDPublicKey:  c.publicKey,
		VAPIDPrivateKey: c.privateKey,
		TTL:             30,
	}

	resp, err := webpushlib.SendNotificationWithContext(
		ctx,
		payload,
		&sub,
		opts,
	)
	if err != nil {
		return 0, fmt.Errorf("webpush library error: %w", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

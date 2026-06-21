package email

import (
	"context"

	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
)

type Emailer interface {
	Send(ctx context.Context, email audit.EmailEntry) error
}

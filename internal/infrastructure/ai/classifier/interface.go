package classifier

import (
	"context"

	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
)

type ServiceInterface interface {
	Classify(
		ctx context.Context,
		text string,
		cfg *config.Config,
	) (*ClassifyResponse, error)
}

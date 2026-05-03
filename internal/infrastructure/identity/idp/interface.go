package idp

import (
	"context"

	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
)

type IDPClientInterface interface {
	ExchangeCodeForToken(
		ctx context.Context,
		code string,
		cfg *config.Config,
	) (*IDPTokenResponse, error)
	GetUserInfo(
		ctx context.Context,
		accessToken string,
		cfg *config.Config,
	) (*IDPUserInfo, error)
	RefreshToken(
		ctx context.Context,
		refreshToken string,
		cfg *config.Config,
	) (*IDPTokenResponse, error)
	GetLogoutURL(cfg *config.Config, idpUserID string) string
}

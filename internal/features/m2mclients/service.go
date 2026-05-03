package m2mclients

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/sessions"
	"github.com/olazo-johnalbert/duckload-api/internal/core/tokens"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Service struct {
	repo           *Repository
	logService     audit.Logger
	notifService   audit.Notifier
	tokenService   *tokens.Service
	sessionService *sessions.Service
}

func NewService(
	repo *Repository,
	logService audit.Logger,
	notifService audit.Notifier,
	tokenService *tokens.Service,
	sessionService *sessions.Service,
) *Service {
	return &Service{
		repo:           repo,
		logService:     logService,
		notifService:   notifService,
		tokenService:   tokenService,
		sessionService: sessionService,
	}
}

func (s *Service) CreateClient(
	ctx context.Context,
	userID string,
	req CreateM2MClientRequest,
) (*CreateM2MClientResponse, error) {
	// Deactivate existing
	err := s.repo.WithTransaction(ctx, func(tx datastore.DB) error {
		return s.repo.DeactivateAllForUser(ctx, tx, userID)
	})
	if err != nil {
		return nil, err
	}

	clientID := uuid.NewString()
	rawSecret, _ := s.generateRandomString(32)
	hashedSecret := s.hashSecret(rawSecret)

	client := M2MClient{
		UserID:            userID,
		ClientName:        req.ClientName,
		ClientDescription: req.ClientDescription,
		ClientID:          clientID,
		ClientSecret:      hashedSecret,
		IsActive:          true,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	err = s.repo.WithTransaction(ctx, func(tx datastore.DB) error {
		return s.repo.Create(ctx, tx, client)
	})
	if err != nil {
		return nil, err
	}

	return &CreateM2MClientResponse{
		M2MClientDTO: M2MClientDTO{
			ClientID:   clientID,
			ClientName: req.ClientName,
			IsActive:   true,
			CreatedAt:  client.CreatedAt,
		},
		ClientSecret: rawSecret,
	}, nil
}

func (s *Service) Authenticate(
	ctx context.Context,
	clientID, clientSecret string,
) (*M2MTokenResponse, error) {
	client, err := s.repo.GetByClientID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("invalid client credentials")
	}

	if client.ClientSecret != s.hashSecret(clientSecret) {
		return nil, fmt.Errorf("invalid client credentials")
	}

	token, claims, err := s.tokenService.GenerateToken(
		client.ClientName,
		client.UserID,
		[]int{int(constants.DeveloperRoleID)},
		"m2m",
		constants.AccessTokenMaxAge,
	)
	if err != nil {
		return nil, err
	}

	val := map[string]string{
		"userID":    client.UserID,
		"tokenType": "m2m",
		"clientID":  clientID,
	}
	err = s.sessionService.StoreToken(
		ctx,
		sessions.NewJTI(claims.ID),
		val,
		constants.AccessTokenMaxAge,
	)
	if err != nil {
		return nil, err
	}

	return &M2MTokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   constants.AccessTokenMaxAge,
	}, nil
}

func (s *Service) ListClients(ctx context.Context, includeRevoked bool) ([]M2MClient, error) {
	return s.repo.ListClients(ctx, includeRevoked)
}

func (s *Service) GetClientByUserID(ctx context.Context, userID string) (*M2MClient, error) {
	return s.repo.GetActiveByUserID(ctx, userID)
}

func (s *Service) ResetSecret(
	ctx context.Context,
	id string,
	requestingUserID string,
	isAdmin bool,
) (string, error) {
	client, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("client not found")
	}

	if !isAdmin && client.UserID != requestingUserID {
		return "", fmt.Errorf("unauthorized: you do not own this client")
	}

	rawSecret, _ := s.generateRandomString(32)
	hashedSecret := s.hashSecret(rawSecret)

	err = s.repo.UpdateSecret(ctx, id, hashedSecret)
	if err != nil {
		return "", fmt.Errorf("failed to rotate client secret: %w", err)
	}
	return rawSecret, nil
}

func (s *Service) Deactivate(
	ctx context.Context,
	id string,
	requestingUserID string,
	isAdmin bool,
) error {
	client, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("client not found")
	}

	if !isAdmin && client.UserID != requestingUserID {
		return fmt.Errorf("unauthorized: you do not own this client")
	}

	return s.repo.DeactivateByID(ctx, id)
}

func (s *Service) Verify(ctx context.Context, id string) error {
	return s.repo.VerifyByID(ctx, id)
}

func (s *Service) generateRandomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *Service) hashSecret(secret string) string {
	hash := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(hash[:])
}

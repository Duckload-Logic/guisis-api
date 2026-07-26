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
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/datetime"
	"github.com/olazo-johnalbert/duckload-api/internal/core/sessions"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/core/tokens"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

const secretLength = 32


type Service struct {
	repo           *Repository
	logService     audit.Logger
	notifService   audit.Notifier
	emailService   audit.Emailer
	userService    *users.Service
	tokenService   *tokens.Service
	sessionService *sessions.Service
	cfg            *config.Config
}

func NewService(
	repo *Repository,
	logService audit.Logger,
	notifService audit.Notifier,
	emailService audit.Emailer,
	userService *users.Service,
	tokenService *tokens.Service,
	sessionService *sessions.Service,
	cfg *config.Config,
) *Service {
	return &Service{
		repo:           repo,
		logService:     logService,
		notifService:   notifService,
		emailService:   emailService,
		userService:    userService,
		tokenService:   tokenService,
		sessionService: sessionService,
		cfg:            cfg,
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
	rawSecret, _ := s.generateRandomString(secretLength)
	hashedSecret := s.hashSecret(rawSecret)

	client := M2MClient{
		UserID:            userID,
		ClientName:        req.ClientName,
		ClientDescription: req.ClientDescription,
		ClientID:          clientID,
		ClientSecret:      hashedSecret,
		IsActive:          true,
		HasPersonalInfoAccess: req.HasPersonalInfoAccess,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	err = s.repo.WithTransaction(ctx, func(tx datastore.DB) error {
		return s.repo.Create(ctx, tx, client)
	})
	if err != nil {
		return nil, err
	}

	// Dispatch notifications
	superadminIDs, _ := s.userService.GetUserIDsByRole(
		ctx,
		int(constants.SuperAdminRoleID),
	)
	notifications := []audit.NotificationParams{
		{
			ReceiverID: structs.StringToNullableString(userID),
			Title:      "M2M Client Request Sent",
			Message: "Your request for a new M2M client " +
				"has been notified to the Superadmin. " +
				"Please wait for approval.",
			Type: constants.SystemEntityType,
		},
	}
	for _, aid := range superadminIDs {
		notifications = append(notifications, audit.NotificationParams{
			ReceiverID: structs.StringToNullableString(aid),
			Title:      "M2M Client Pending Approval",
			Message: fmt.Sprintf(
				"New M2M client request from user %s "+
					"is pending for approval.",
				userID,
			),
			Type: constants.SystemEntityType,
		})
	}

	superadminEmails, _ := s.userService.GetEmailsByRole(
		ctx,
		int(constants.SuperAdminRoleID),
	)

	audit.Dispatch(
		ctx,
		s.logService,
		s.notifService,
		s.emailService,
		audit.DispatchParams{
			Log: &audit.LogParams{
				Level:    audit.LevelInfo,
				Category: audit.CategoryAudit,
				Action:   audit.ActionM2MClientCreated,
				Message: fmt.Sprintf(
					"M2M Client %s requested by user %s",
					clientID,
					userID,
				),
				Metadata: &audit.LogMetadata{
					EntityType: "m2m_client",
					EntityID:   clientID,
				},
			},
			Notifications: notifications,
			Email: []audit.EmailParams{
				{
					To:           superadminEmails,
					Subject:      "New M2M Client Request",
					TemplatePath: "request.html",
					TemplateData: map[string]any{
						"EntityType":  "M2M Client",
						"StudentName": userID,
						"Category":    "M2M Access",
						"Reason":      req.ClientDescription,
						"TimeSlot": datetime.FormatDateTime(time.Now()),
						"UrgencyLevel": "HIGH",
						"RequestURL": fmt.Sprintf(
							"%s/superadmin/m2m-management",
							s.cfg.BaseURL,
						),
					},
				},
			},
		},
	)

	return &CreateM2MClientResponse{
		M2MClientDTO: M2MClientDTO{
			ClientID:   clientID,
			ClientName: req.ClientName,
			IsActive:   true,
			HasPersonalInfoAccess: client.HasPersonalInfoAccess,
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

	return s.issueTokens(ctx, client)
}

func (s *Service) RefreshM2MToken(
	ctx context.Context,
	refreshToken string,
) (*M2MTokenResponse, error) {
	claims, err := s.tokenService.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired refresh token")
	}

	if claims.TokenType != "m2m_refresh" {
		return nil, fmt.Errorf("invalid token type")
	}

	session, err := s.sessionService.GetToken(ctx, sessions.NewJTI(claims.ID))
	if err != nil {
		return nil, fmt.Errorf("refresh session expired or revoked")
	}

	clientID := session["clientID"]
	if clientID == "" {
		return nil, fmt.Errorf("invalid session data")
	}

	// Verify client is still active
	client, err := s.repo.GetByClientID(ctx, clientID)
	if err != nil || !client.IsActive {
		return nil, fmt.Errorf("client is inactive or revoked")
	}

	return s.issueTokens(ctx, client)
}

func (s *Service) issueTokens(
	ctx context.Context,
	client *M2MClient,
) (*M2MTokenResponse, error) {
	// Access Token
	token, claims, err := s.tokenService.GenerateM2MToken(
		client.ClientName,
		client.UserID,
		[]int{int(constants.DeveloperRoleID)},
		"m2m",
		client.ClientID,
		client.HasPersonalInfoAccess,
		client.IsVerified,
		constants.M2MAccessTokenMaxAge,
	)
	if err != nil {
		return nil, err
	}

	isVerifiedStr := "false"
	if client.IsVerified {
		isVerifiedStr = "true"
	}

	val := map[string]string{
		"userID":     client.UserID,
		"tokenType":  "m2m",
		"clientID":   client.ClientID,
		"isVerified": isVerifiedStr,
	}
	err = s.sessionService.StoreToken(
		ctx,
		sessions.NewJTI(claims.ID),
		val,
		constants.M2MAccessTokenMaxAge,
	)
	if err != nil {
		return nil, err
	}

	// Refresh Token
	refreshToken, rClaims, err := s.tokenService.GenerateM2MToken(
		client.ClientName,
		client.UserID,
		[]int{int(constants.DeveloperRoleID)},
		"m2m_refresh",
		client.ClientID,
		client.HasPersonalInfoAccess,
		client.IsVerified,
		constants.M2MRefreshTokenMaxAge,
	)
	if err != nil {
		return nil, err
	}

	rVal := map[string]string{
		"userID":     client.UserID,
		"tokenType":  "m2m_refresh",
		"clientID":   client.ClientID,
		"isVerified": isVerifiedStr,
	}
	err = s.sessionService.StoreToken(
		ctx,
		sessions.NewJTI(rClaims.ID),
		rVal,
		constants.M2MRefreshTokenMaxAge,
	)
	if err != nil {
		return nil, err
	}

	return &M2MTokenResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    constants.M2MAccessTokenMaxAge,
	}, nil
}

func (s *Service) ListClients(
	ctx context.Context,
	params ListM2MClientsRequest,
) ([]M2MClient, error) {
	return s.repo.ListClients(ctx, params)
}

func (s *Service) GetClientByUserID(
	ctx context.Context,
	userID string,
) (*M2MClient, error) {
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

	rawSecret, _ := s.generateRandomString(secretLength)
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

func (s *Service) Verify(
	ctx context.Context,
	id string,
	hasPersonalInfoAccess bool,
) error {
	err := s.repo.VerifyByID(ctx, id, hasPersonalInfoAccess)
	if err != nil {
		return err
	}

	client, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf(
			"failed to fetch client after verification: %w",
			err,
		)
	}

	user, _ := s.userService.GetUserByID(ctx, client.UserID)

	audit.Dispatch(
		ctx,
		s.logService,
		s.notifService,
		s.emailService,
		audit.DispatchParams{
			Log: &audit.LogParams{
				Level:    audit.LevelInfo,
				Category: audit.CategoryAudit,
				Action:   audit.ActionM2MClientVerified,
				Message:  fmt.Sprintf("M2M Client %s verified", id),
			},
			Notifications: []audit.NotificationParams{
				{
					ReceiverID: structs.StringToNullableString(
						client.UserID,
					),
					Title: "M2M Client Approved",
					Message: fmt.Sprintf(
						"Your M2M client '%s' has " +
							"been approved and is ready for use.",
						client.ClientName,
					),
					Type: constants.SystemEntityType,
				},
			},
			Email: []audit.EmailParams{
				{
					To:           []string{user.Email},
					Subject:      "M2M Client Approved",
					TemplatePath: "m2m.html",
					TemplateData: map[string]interface{}{
						"Status":     "Approved",
						"ClientName": client.ClientName,
						"ClientID":   client.ClientID,
					},
				},
			},
		},
	)

	return nil
}

func (s *Service) Reject(ctx context.Context, id string) error {
	client, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf(
			"failed to fetch client before rejection: %w",
			err,
		)
	}

	err = s.repo.DeactivateByID(ctx, id)
	if err != nil {
		return err
	}

	user, _ := s.userService.GetUserByID(ctx, client.UserID)

	audit.Dispatch(
		ctx,
		s.logService,
		s.notifService,
		s.emailService,
		audit.DispatchParams{
			Log: &audit.LogParams{
				Level:    audit.LevelInfo,
				Category: audit.CategoryAudit,
				Action:   audit.ActionM2MClientRevoked,
				Message:  fmt.Sprintf("M2M Client %s rejected", id),
			},
			Notifications: []audit.NotificationParams{
				{
					ReceiverID: structs.StringToNullableString(
						client.UserID,
					),
					Title: "M2M Client Rejected",
					Message: fmt.Sprintf(
						"Your M2M client request '%s' " +
							"has been rejected.",
						client.ClientName,
					),
					Type: constants.SystemEntityType,
				},
			},
			Email: []audit.EmailParams{
				{
					To:           []string{user.Email},
					Subject:      "M2M Client Rejected",
					TemplatePath: "m2m.html",
					TemplateData: map[string]interface{}{
						"Status":     "Rejected",
						"ClientName": client.ClientName,
						"ClientID":   client.ClientID,
					},
				},
			},
		},
	)

	return nil
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

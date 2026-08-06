package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/olazo-johnalbert/duckload-api/internal/core/audit"
	"github.com/olazo-johnalbert/duckload-api/internal/core/config"
	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/core/sessions"
	"github.com/olazo-johnalbert/duckload-api/internal/core/structs"
	"github.com/olazo-johnalbert/duckload-api/internal/core/tokens"
	"github.com/olazo-johnalbert/duckload-api/internal/features/users"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/identity/idp"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo           *users.Repository
	idpClient      idp.IDPClient
	redis          *datastore.RedisClient
	sessionService *sessions.Service
	emailer        audit.Emailer
	logger         audit.Logger
}

func NewService(
	repo *users.Repository,
	redis *datastore.RedisClient,
	sessionService *sessions.Service,
	emailer audit.Emailer,
	logger audit.Logger,
) *Service {
	return &Service{
		repo:           repo,
		idpClient:      *idp.NewIDPClient(),
		redis:          redis,
		sessionService: sessionService,
		emailer:        emailer,
		logger:         logger,
	}
}

func (s *Service) validateEmailDomain(email string) error {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return fmt.Errorf("invalid email format")
	}
	domain := parts[1]
	mx, err := net.LookupMX(domain)
	if err != nil || len(mx) == 0 {
		return fmt.Errorf("email domain %s is unreachable", domain)
	}
	return nil
}

// AuthenticateUser handles native email/password authentication.
// TODO: This will be depracated
func (s *Service) AuthenticateUser(
	ctx context.Context, email, password, ipAddress, userAgent string,
) (string, string, string, error) {
	const (
		lockoutPrefix   = "lockout:"
		failurePrefix   = "failed_attempts:"
		maxFailures     = 5
		lockoutDuration = 15 * time.Minute
	)

	// Check if user is locked out
	// TODO: Remove this implementation in the future
	lockoutKey := fmt.Sprintf("%s%s", lockoutPrefix, email)
	if locked, _ := s.redis.Get(ctx, lockoutKey); locked != "" {
		log.Printf("Account locked for user: %s", email)
		return "", "", "", fmt.Errorf(
			"account locked due to too many failed attempts. " +
				"Please try again in 15 minutes",
		)
	}

	// Fetch user from database (Native only)
	user, err := s.repo.GetUserByEmail(
		ctx,
		email,
		string(constants.AuthTypeNative),
	)
	if err != nil {
		s.logger.Record(ctx, nil, audit.LogEntry{
			Level:     audit.LevelWarning,
			Category:  audit.CategorySecurity,
			Action:    audit.ActionLoginFailed,
			Message:   fmt.Sprintf("Failed login attempt: %s", email),
			UserEmail: structs.StringToNullableString(email),
			IPAddress: structs.StringToNullableString(ipAddress),
			UserAgent: structs.StringToNullableString(userAgent),
		})
		return "", "", "", fmt.Errorf("invalid credentials")
	}

	if !user.IsActive {
		s.logger.Record(ctx, nil, audit.LogEntry{
			Level:    audit.LevelWarning,
			Category: audit.CategorySecurity,
			Action:   audit.ActionLoginFailed,
			Message: fmt.Sprintf(
				"Failed login attempt (Inactive): %s",
				email,
			),
			UserEmail: structs.StringToNullableString(email),
			IPAddress: structs.StringToNullableString(ipAddress),
			UserAgent: structs.StringToNullableString(userAgent),
		})
		return "", "", "", fmt.Errorf("invalid credentials")
	}

	// Compare hashed password
	if !user.PasswordHash.Valid {
		s.logger.Record(ctx, nil, audit.LogEntry{
			Level:     audit.LevelWarning,
			Category:  audit.CategorySecurity,
			Action:    audit.ActionLoginFailed,
			Message:   fmt.Sprintf("Failed login attempt (No hash): %s", email),
			UserEmail: structs.StringToNullableString(email),
			IPAddress: structs.StringToNullableString(ipAddress),
			UserAgent: structs.StringToNullableString(userAgent),
		})
		return "", "", "", fmt.Errorf("invalid credentials")
	}

	failureKey := fmt.Sprintf("%s%s", failurePrefix, email)
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash.String),
		[]byte(password),
	)
	if err != nil {
		// Increment failures
		failuresStr, _ := s.redis.Get(ctx, failureKey)
		failures := 0
		if failuresStr != "" {
			fmt.Sscanf(failuresStr, "%d", &failures)
		}

		failures++

		if failures >= maxFailures {
			err = s.redis.Set(
				ctx,
				lockoutKey,
				"true",
				lockoutDuration,
			)
			if err != nil {
				return "", "", "", fmt.Errorf(
					"[REDIS:SET-LOCKOUT]:%v", err,
				)
			}
			err = s.redis.Del(ctx, failureKey)
			if err != nil {
				return "", "", "", fmt.Errorf(
					"[REDIS:DEL-FAILURES]:%v", err,
				)
			}

			s.logger.Record(ctx, nil, audit.LogEntry{
				Level:    audit.LevelCritical,
				Category: audit.CategorySecurity,
				Action:   audit.ActionLoginFailed,
				Message: fmt.Sprintf(
					"Account locked due to too many failed attempts: %s",
					email,
				),
				UserEmail: structs.StringToNullableString(email),
				IPAddress: structs.StringToNullableString(ipAddress),
				UserAgent: structs.StringToNullableString(userAgent),
			})

			return "", "", "", fmt.Errorf(
				"account locked due to too many failed attempts. " +
					"please try again in 15 minutes",
			)
		}

		err = s.redis.Set(
			ctx,
			failureKey,
			fmt.Sprintf("%d", failures),
			lockoutDuration, // lockout period
		)
		if err != nil {
			return "", "", "", fmt.Errorf("[REDIS:SET-FAILURES]:%v", err)
		}

		s.logger.Record(ctx, nil, audit.LogEntry{
			Level:     audit.LevelWarning,
			Category:  audit.CategorySecurity,
			Action:    audit.ActionLoginFailed,
			Message:   fmt.Sprintf("Failed login attempt: %s", email),
			UserEmail: structs.StringToNullableString(email),
			IPAddress: structs.StringToNullableString(ipAddress),
			UserAgent: structs.StringToNullableString(userAgent),
		})

		return "", "", "", fmt.Errorf("invalid credentials")
	}

	// Success: Reset failures and lockout
	err = s.redis.Del(ctx, failureKey)
	if err != nil {
		return "", "", "", fmt.Errorf("[REDIS:DEL-FAILURES]:%v", err)
	}
	err = s.redis.Del(ctx, lockoutKey)
	if err != nil {
		return "", "", "", fmt.Errorf("[REDIS:DEL-LOCKOUT]:%v", err)
	}

	// Look up student IDs if they have Student role
	var iirID, corID string
	for _, r := range user.Roles {
		if r.ID == int(constants.StudentRoleID) {
			_ = s.repo.GetDB().QueryRow(`
				SELECT id FROM iir_records WHERE user_id = ?
			`, user.ID).Scan(&iirID)
			_ = s.repo.GetDB().QueryRow(`
				SELECT file_id FROM student_cors WHERE student_id = ?
			`, user.ID).Scan(&corID)
			break
		}
	}

	// Generate the token
	roleIDs := make([]int, len(user.Roles))
	for i, r := range user.Roles {
		roleIDs[i] = r.ID
	}
	token, accessClaims, err := tokens.NewService().GenerateSessionToken(
		user.Email,
		user.ID,
		roleIDs,
		string(constants.AuthTypeNative),
		constants.AccessTokenMaxAge,
		func(c *tokens.Claims) {
			c.IsVerified = user.IsActive
			c.IIRID = iirID
			c.CORID = corID
		},
	)
	if err != nil {
		return "", "", "", fmt.Errorf(
			"failed to generate session: %v",
			err,
		)
	}

	// Generate refresh token
	refreshToken, refreshClaims, err := tokens.NewService().
		GenerateSessionToken(
			user.Email,
			user.ID,
			roleIDs,
			string(constants.AuthTypeNative),
			constants.RefreshTokenMaxAge,
			func(c *tokens.Claims) {
				c.IsVerified = user.IsActive
				c.IIRID = iirID
				c.CORID = corID
			},
		)
	if err != nil {
		return "", "", "", fmt.Errorf(
			"failed to generate refresh token: %v",
			err,
		)
	}

	// Whitelist the session in Redis
	err = s.sessionService.WhitelistSession(
		ctx,
		user.ID,
		accessClaims.ID,
		refreshClaims.ID,
		constants.RefreshTokenMaxAge,
	)
	if err != nil {
		return "", "", "", fmt.Errorf(
			"failed to whitelist session: %w",
			err,
		)
	}
	// Log successful login
	s.logger.Record(ctx, nil, audit.LogEntry{
		Level:     audit.LevelInfo,
		Category:  audit.CategorySecurity,
		Action:    audit.ActionLoginSuccess,
		Message:   fmt.Sprintf("User %s successfully logged in", user.Email),
		UserID:    structs.StringToNullableString(user.ID),
		UserEmail: structs.StringToNullableString(user.Email),
		IPAddress: structs.StringToNullableString(ipAddress),
		UserAgent: structs.StringToNullableString(userAgent),
	})

	return user.ID, token, refreshToken, nil
}

// RefreshToken generates a new access token using a valid session handle.
func (s *Service) RefreshToken(
	ctx context.Context,
	refreshToken string,
	cfg *config.Config,
	ipAddress, userAgent string,
) (string, string, error) {
	// 1. Validate the App Refresh Token
	claims, err := tokens.NewService().ValidateToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %v", err)
	}
	// 2. Check if the refresh token JTI is whitelisted
	whitelist, err := s.sessionService.GetWhitelistedSession(
		ctx,
		claims.UserID,
	)
	if err != nil || len(whitelist) == 0 {
		return "", "", fmt.Errorf("refresh token has been revoked")
	}

	whitelistedRefreshJTI := whitelist[constants.RedisSessionRefreshJTIField]
	if whitelistedRefreshJTI != claims.ID {
		return "", "", fmt.Errorf("refresh token has been revoked")
	}
	var iirID, corID string
	_ = s.repo.GetDB().QueryRowContext(ctx, `
		SELECT id FROM iir_records WHERE user_id = ?
	`, claims.UserID).Scan(&iirID)
	_ = s.repo.GetDB().QueryRowContext(ctx, `
		SELECT file_id FROM student_cors WHERE student_id = ?
	`, claims.UserID).Scan(&corID)

	// Check token type
	if claims.TokenType == string(constants.AuthTypeIDP) {
		// Get IDP refresh token from Redis using the OLD refresh token's ID
		idpRefreshKey := sessions.NewJTI(claims.ID).ToIDPRefreshKey()
		idpRefreshToken, err := s.redis.Get(ctx, idpRefreshKey)
		if err != nil {
			return "", "", fmt.Errorf(
				"[AuthService] {Get IDP Refresh Token}: idp token missing",
			)
		}

		// Call IDP refresh endpoint
		tokenResp, err := s.idpClient.RefreshToken(ctx, idpRefreshToken, cfg)
		if err != nil {
			return "", "", fmt.Errorf("[AuthService] {IDP Refresh}: %w", err)
		}

		newAppAccessToken, newAppAccessClaims, err := tokens.NewService().
			GenerateSessionToken(
				claims.UserEmail,
				claims.UserID,
				claims.RoleIDs,
				string(constants.AuthTypeIDP),
				constants.AccessTokenMaxAge,
				func(c *tokens.Claims) {
					c.IsVerified = true
					c.IDPAccessToken = tokenResp.AccessToken
					c.IIRID = iirID
					c.CORID = corID
				},
			)
		if err != nil {
			return "", "", err
		}

		newAppRefreshToken, refreshClaims, err := tokens.NewService().
			GenerateSessionToken(
				claims.UserEmail,
				claims.UserID,
				claims.RoleIDs,
				string(constants.AuthTypeIDP),
				constants.RefreshTokenMaxAge,
				func(c *tokens.Claims) {
					c.IsVerified = true
					c.IDPAccessToken = tokenResp.AccessToken
					c.IIRID = iirID
					c.CORID = corID
				},
			)
		if err != nil {
			return "", "", err
		}

		// Whitelist the new session in Redis
		err = s.sessionService.WhitelistSession(
			ctx,
			claims.UserID,
			newAppAccessClaims.ID,
			refreshClaims.ID,
			constants.RefreshTokenMaxAge,
		)
		if err != nil {
			return "", "", err
		}

		// Update Redis: IDP Refresh linked to NEW App Refresh Token's ID
		newIdpRefreshKey := sessions.NewJTI(refreshClaims.ID).ToIDPRefreshKey()
		idpRefreshTokenToStore := tokenResp.RefreshToken
		if idpRefreshTokenToStore == "" {
			idpRefreshTokenToStore = idpRefreshToken
		}
		err = s.redis.Set(
			ctx,
			newIdpRefreshKey,
			idpRefreshTokenToStore,
			time.Duration(constants.RefreshTokenMaxAge)*time.Second,
		)
		if err != nil {
			return "", "", err
		}

		_ = s.redis.Del(ctx, idpRefreshKey)

		return newAppAccessToken, newAppRefreshToken, nil
	}

	// Native flow
	newToken, newAccessClaims, err := tokens.NewService().GenerateSessionToken(
		claims.UserEmail,
		claims.UserID,
		claims.RoleIDs,
		string(constants.AuthTypeNative),
		constants.AccessTokenMaxAge,
		func(c *tokens.Claims) {
			c.IsVerified = true
			c.IIRID = iirID
			c.CORID = corID
		},
	)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, newRefreshClaims, err := tokens.NewService().
		GenerateSessionToken(
			claims.UserEmail,
			claims.UserID,
			claims.RoleIDs,
			string(constants.AuthTypeNative),
			constants.RefreshTokenMaxAge,
			func(c *tokens.Claims) {
				c.IsVerified = true
				c.IIRID = iirID
				c.CORID = corID
			},
		)
	if err != nil {
		return "", "", err
	}

	// Whitelist the new session in Redis
	err = s.sessionService.WhitelistSession(
		ctx,
		claims.UserID,
		newAccessClaims.ID,
		newRefreshClaims.ID,
		constants.RefreshTokenMaxAge,
	)
	if err != nil {
		return "", "", err
	}

	return newToken, newRefreshToken, nil
}

// GetMe retrieves the currently authenticated user's profile information.
func (s *Service) GetMe(
	ctx context.Context,
	userID, tokenType string,
) (*MeResponse, error) {
	// only fetch user info for native tokens
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	resp := &MeResponse{
		ID:         user.ID,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		SuffixName: user.SuffixName.String,
		MiddleName: user.MiddleName.String,
		CreatedAt:  user.CreatedAt.Time,
		Roles:      user.Roles,
		Type:       tokenType,
	}

	profilePicture, err := s.repo.GetProfilePictureURLByUserID(ctx, userID)
	if err == nil {
		resp.ProfilePicture = profilePicture
	} else if err != sql.ErrNoRows {
		log.Printf("[GetMe] {Profile Picture Fetch}: %v", err)
	}

	// Fetch COR for students
	for _, role := range user.Roles {
		if role.ID == int(constants.StudentRoleID) {
			corURL, err := s.repo.GetStudentCORURLByUserID(ctx, userID)
			if err == nil {
				resp.StudentCORURL = corURL
				valid, err := s.repo.CheckStudentCORValidByUserID(ctx, userID)
				if err == nil {
					resp.IsStudentCORValid = valid
				} else {
					log.Printf("[GetMe] {Student COR Validity Fetch}: %v", err)
				}
			} else if err != sql.ErrNoRows {
				log.Printf("[GetMe] {Student COR Fetch}: %v", err)
			}
			break
		}
	}

	return resp, nil
}

// Logout invalidates the user's session in Redis and optionally the IDP.
func (s *Service) Logout(
	ctx context.Context,
	token string,
	refreshToken string,
	tokenType string,
	cfg *config.Config,
) (string, error) {
	var idpToken string

	// Identify the session using the Access Token JTI
	claims, err := tokens.NewService().ParseTokenUnverified(token)
	if err == nil && claims.ID != "" {
		_ = s.sessionService.DeleteUserToken(
			ctx,
			claims.UserID,
			sessions.NewJTI(claims.ID),
		)
		idpToken = claims.IDPAccessToken
	}

	// Identify and blacklist the Refresh Token JTI
	rClaims, err := tokens.NewService().ParseTokenUnverified(refreshToken)
	if err == nil && rClaims.ID != "" {
		_ = s.sessionService.DeleteUserToken(
			ctx,
			rClaims.UserID,
			sessions.NewJTI(rClaims.ID),
		)
		idpKey := sessions.NewJTI(rClaims.ID).ToIDPRefreshKey()
		_ = s.redis.Del(ctx, idpKey)
		if idpToken == "" {
			idpToken = rClaims.IDPAccessToken
		}
	}

	userInfo, _ := s.idpClient.GetUserInfo(ctx, idpToken, cfg)

	// Construct logout URL for front-channel redirect
	if tokenType == string(constants.AuthTypeIDP) && userInfo != nil {
		return s.idpClient.GetLogoutURL(cfg, userInfo.ID), nil
	}

	// Fallback redirect for native logout or incomplete IDP sessions
	return "", nil
}

// IDP integration methods

// GetAuthorizeURL generates the complete OAuth 2.0 authorization URL
func (s *Service) GetAuthorizeURL(
	cfg *config.Config,
) (string, error) {
	// Build authorization URL with all required parameters
	params := url.Values{}
	params.Set("client_id", cfg.IDPClientID)

	authURL := fmt.Sprintf(
		"%s?%s",
		fmt.Sprintf("%s/auth/authorize", cfg.IDPBaseUrl),
		params.Encode(),
	)

	return authURL, nil
}

// PostIDPTokenExchange orchestrates the complete IDP login flow
func (s *Service) PostIDPTokenExchange(
	ctx context.Context,
	code string,
	cfg *config.Config,
	ipAddress, userAgent string,
) (string, string, error) {
	// Exchange authorization code for IDP tokens
	idpTokenResp, err := s.idpClient.ExchangeCodeForToken(ctx, code, cfg)
	if err != nil {
		return "", "", fmt.Errorf(
			"[AuthService] {Token Exchange}: %w",
			err,
		)
	}

	// Fetch User Info from IDP
	userInfo, err := s.GetIDPUserInfo(ctx, idpTokenResp.AccessToken, cfg)
	if err != nil {
		return "", "", fmt.Errorf(
			"[AuthService] {Fetch User Info}: %w",
			err,
		)
	}

	if userInfo.ID == "" {
		return "", "", fmt.Errorf(
			"[AuthService] {IDP Login}: IDP user ID is empty",
		)
	}

	// Whitelist Check: Authoritative role source for IDP users.
	// We check this on every login to support dynamic role changes (promotions).
	whitelistRoleIDs, whitelistErr := s.repo.CheckUserWhitelist(
		ctx,
		userInfo.Email,
	)

	// User Existence Check by IDP UUID
	localUser, err := s.repo.GetUserByIDPUUID(ctx, userInfo.ID)
	if err == sql.ErrNoRows {
		// Fallback to checking by email for pre-existing native accounts
		localUser, err = s.repo.GetUserByEmail(
			ctx,
			userInfo.Email,
			string(constants.AuthTypeNative),
		)
	}

	if err == nil {
		var changed bool
		if localUser.Email != userInfo.Email {
			localUser.Email = userInfo.Email
			changed = true
		}
		if localUser.FirstName != userInfo.FirstName {
			localUser.FirstName = userInfo.FirstName
			changed = true
		}
		if localUser.LastName != userInfo.LastName {
			localUser.LastName = userInfo.LastName
			changed = true
		}

		newMiddle := structs.StringToNullableString(
			userInfo.MiddleName,
		)
		if localUser.MiddleName.String != newMiddle.String ||
			localUser.MiddleName.Valid != newMiddle.Valid {
			localUser.MiddleName = newMiddle
			changed = true
		}

		newSuffix := structs.StringToNullableString(
			userInfo.SuffixName,
		)
		if localUser.SuffixName.String != newSuffix.String ||
			localUser.SuffixName.Valid != newSuffix.Valid {
			localUser.SuffixName = newSuffix
			changed = true
		}

		if localUser.AuthType != string(constants.AuthTypeIDP) {
			localUser.AuthType = string(constants.AuthTypeIDP)
			changed = true
		}

		if localUser.IDPUUID.String != userInfo.ID ||
			!localUser.IDPUUID.Valid {
			localUser.IDPUUID = structs.StringToNullableString(
				userInfo.ID,
			)
			changed = true
		}

		if changed {
			err = s.repo.WithTransaction(
				ctx,
				func(tx datastore.DB) error {
					return s.repo.CreateUser(ctx, tx, *localUser)
				},
			)
			if err != nil {
				return "", "", fmt.Errorf(
					"[PostIDPTokenExchange] {Update User Info}: %w",
					err,
				)
			}
		}
	}

	// Determine target roles from whitelist or defaults
	var targetRoleIDs []int
	if whitelistErr != nil {
		return "", "", fmt.Errorf(
			"[AuthService] {Whitelist Check}: %w",
			whitelistErr,
		)
	}

	if len(whitelistRoleIDs) > 0 {
		targetRoleIDs = whitelistRoleIDs
	}

	if err == sql.ErrNoRows {
		if len(targetRoleIDs) == 0 {
			targetRoleIDs = []int{int(constants.StudentRoleID)}
		}

		// JIT Provisioning (First Login Only)
		roles := make([]users.Role, len(targetRoleIDs))
		for i, id := range targetRoleIDs {
			roles[i] = users.Role{ID: id}
		}

		localUser = &users.User{
			ID:           uuid.NewString(),
			Email:        userInfo.Email,
			Roles:        roles,
			FirstName:    userInfo.FirstName,
			LastName:     userInfo.LastName,
			MiddleName:   structs.StringToNullableString(userInfo.MiddleName),
			SuffixName:   structs.StringToNullableString(userInfo.SuffixName),
			IDPUUID:      structs.StringToNullableString(userInfo.ID),
			AuthType:     string(constants.AuthTypeIDP),
			PasswordHash: structs.NullableString{Valid: false},
			IsActive:     true,
		}

		err = s.repo.WithTransaction(
			ctx,
			func(tx datastore.DB) error {
				if err := s.repo.CreateUser(ctx, tx, *localUser); err != nil {
					return err
				}
				if len(whitelistRoleIDs) > 0 {
					return s.repo.RemoveUserFromWhitelist(
						ctx,
						tx,
						userInfo.Email,
					)
				}
				return nil
			},
		)
		if err != nil {
			return "", "", fmt.Errorf(
				"[AuthService] {Provision IDP User}: %w",
				err,
			)
		}
	} else if err != nil {
		return "", "", fmt.Errorf(
			"[AuthService] {Anchor Check}: %w",
			err,
		)
	}

	roleIDs := make([]int, len(localUser.Roles))
	for i, r := range localUser.Roles {
		roleIDs[i] = r.ID
	}

	// Hydrate Student IIR and COR IDs
	var iirID, corID string
	for _, r := range localUser.Roles {
		if r.ID == int(constants.StudentRoleID) {
			_ = s.repo.GetDB().QueryRow(`
				SELECT id FROM iir_records WHERE user_id = ?
			`, localUser.ID).Scan(&iirID)
			_ = s.repo.GetDB().QueryRow(`
				SELECT file_id FROM student_cors WHERE student_id = ?
			`, localUser.ID).Scan(&corID)
			break
		}
	}
	appAccessToken, accessClaims, err := tokens.NewService().
		GenerateSessionToken(
			userInfo.Email,
			localUser.ID,
			roleIDs,
			string(constants.AuthTypeIDP),
			constants.AccessTokenMaxAge,
			func(c *tokens.Claims) {
				c.IsVerified = localUser.IsActive
				c.IDPAccessToken = idpTokenResp.AccessToken
				c.IIRID = iirID
				c.CORID = corID
			},
		)
	if err != nil {
		return "", "", err
	}

	appRefreshToken, refreshClaims, err := tokens.NewService().
		GenerateSessionToken(
			userInfo.Email,
			localUser.ID,
			roleIDs,
			string(constants.AuthTypeIDP),
			constants.RefreshTokenMaxAge,
			func(c *tokens.Claims) {
				c.IsVerified = localUser.IsActive
				c.IDPAccessToken = idpTokenResp.AccessToken
				c.IIRID = iirID
				c.CORID = corID
			},
		)
	if err != nil {
		return "", "", err
	}

	// Whitelist the session in Redis
	err = s.sessionService.WhitelistSession(
		ctx,
		localUser.ID,
		accessClaims.ID,
		refreshClaims.ID,
		constants.RefreshTokenMaxAge,
	)
	if err != nil {
		return "", "", err
	}

	idpRefreshKey := sessions.NewJTI(refreshClaims.ID).ToIDPRefreshKey()
	err = s.redis.Set(
		ctx,
		idpRefreshKey,
		idpTokenResp.RefreshToken,
		time.Duration(constants.RefreshTokenMaxAge)*time.Second,
	)
	if err != nil {
		return "", "", err
	}
	// Log successful login
	s.logger.Record(ctx, nil, audit.LogEntry{
		Level:    audit.LevelInfo,
		Category: audit.CategorySecurity,
		Action:   audit.ActionLoginSuccess,
		Message: fmt.Sprintf(
			"User %s successfully logged in via IDP",
			userInfo.Email,
		),
		UserID:    structs.StringToNullableString(localUser.ID),
		UserEmail: structs.StringToNullableString(userInfo.Email),
		IPAddress: structs.StringToNullableString(ipAddress),
		UserAgent: structs.StringToNullableString(userAgent),
	})

	return appAccessToken, appRefreshToken, nil
}

// GetIDPUserInfo fetches user information from the IDP
func (s *Service) GetIDPUserInfo(
	ctx context.Context,
	accessToken string,
	cfg *config.Config,
) (*idp.IDPUserInfo, error) {
	userInfo, err := s.idpClient.GetUserInfo(ctx, accessToken, cfg)
	if err != nil {
		return nil, fmt.Errorf(
			"[AuthService] {Get IDP User Info}: %w",
			err,
		)
	}
	return userInfo, nil
}

// IsIDPUp checks if the IDP is up, utilizing a short-lived cache in Redis.
func (s *Service) IsIDPUp(
	ctx context.Context,
	cfg *config.Config,
) bool {
	const cacheKey = "idp_health_status"
	const cacheTTL = 30 * time.Second

	if s.redis != nil {
		status, err := s.redis.Get(ctx, cacheKey)
		if err == nil && status != "" {
			return status == "up"
		}
	}

	err := s.idpClient.PingIDP(ctx, cfg)
	if err != nil {
		log.Printf("[AuthService] {IsIDPUp}: IDP down: %v", err)
		if s.redis != nil {
			_ = s.redis.Set(ctx, cacheKey, "down", cacheTTL)
		}
		return false
	}

	if s.redis != nil {
		_ = s.redis.Set(ctx, cacheKey, "up", cacheTTL)
	}
	return true
}

// GenerateAndSendOTP generates a cryptographically secure OTP,
// saves it to Redis, and sends it to the user's email.
func (s *Service) GenerateAndSendOTP(
	ctx context.Context,
	email string,
) error {
	// First verify that the user exists in the system
	// We support native or idp fallback login for registered users
	_, err := s.repo.GetUserByEmail(
		ctx,
		email,
		string(constants.AuthTypeIDP),
	)
	if err != nil {
		// Fallback to checking native user existence
		_, err = s.repo.GetUserByEmail(
			ctx,
			email,
			string(constants.AuthTypeNative),
		)
		if err != nil {
			return fmt.Errorf(
				"[AuthService] {Verify User}: user not found",
			)
		}
	}

	// Generate 6-digit OTP
	maxVal := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, maxVal)
	if err != nil {
		return fmt.Errorf("[AuthService] {Generate OTP}: %w", err)
	}
	otp := fmt.Sprintf("%06d", n.Int64())

	// Store OTP in Redis with 5 minutes TTL
	redisKey := fmt.Sprintf("otp:%s", email)
	err = s.redis.Set(ctx, redisKey, otp, 5*time.Minute)
	if err != nil {
		return fmt.Errorf("[AuthService] {Store OTP}: %w", err)
	}

	// Send email
	emailEntry := audit.EmailEntry{
		To:           []string{email},
		Subject:      "PUPT GuiSIS - Verification Code",
		TemplatePath: "otp.html",
		TemplateData: map[string]string{
			"OTP": otp,
		},
	}
	err = s.emailer.Send(ctx, emailEntry)
	if err != nil {
		return fmt.Errorf("[AuthService] {Send OTP Email}: %w", err)
	}

	return nil
}

// AuthenticateUserOTP authenticates a user using a verified OTP.
func (s *Service) AuthenticateUserOTP(
	ctx context.Context,
	email, otp string,
	ipAddress, userAgent string,
) (string, string, string, error) {
	redisKey := fmt.Sprintf("otp:%s", email)
	cachedOTP, err := s.redis.Get(ctx, redisKey)
	if err != nil || cachedOTP == "" {
		return "", "", "", fmt.Errorf(
			"verification code expired or invalid",
		)
	}

	if cachedOTP != otp {
		return "", "", "", fmt.Errorf("invalid verification code")
	}

	// Consume OTP immediately to prevent reuse
	_ = s.redis.Del(ctx, redisKey)

	// Fetch user (either IDP or native fallback type)
	user, err := s.repo.GetUserByEmail(
		ctx,
		email,
		string(constants.AuthTypeIDP),
	)
	authTypeUsed := string(constants.AuthTypeIDP)
	if err != nil {
		user, err = s.repo.GetUserByEmail(
			ctx,
			email,
			string(constants.AuthTypeNative),
		)
		authTypeUsed = string(constants.AuthTypeNative)
		if err != nil {
			return "", "", "", fmt.Errorf("user not found")
		}
	}

	if !user.IsActive {
		return "", "", "", fmt.Errorf("account is inactive")
	}

	// Retrieve student IDs if applicable
	var iirID, corID string
	for _, r := range user.Roles {
		if r.ID == int(constants.StudentRoleID) {
			_ = s.repo.GetDB().QueryRowContext(ctx, `
				SELECT id FROM iir_records WHERE user_id = ?
			`, user.ID).Scan(&iirID)
			_ = s.repo.GetDB().QueryRowContext(ctx, `
				SELECT file_id FROM student_cors WHERE student_id = ?
			`, user.ID).Scan(&corID)
			break
		}
	}

	// Generate session tokens
	roleIDs := make([]int, len(user.Roles))
	for i, r := range user.Roles {
		roleIDs[i] = r.ID
	}

	token, accessClaims, err := tokens.NewService().GenerateSessionToken(
		user.Email,
		user.ID,
		roleIDs,
		authTypeUsed,
		constants.AccessTokenMaxAge,
		func(c *tokens.Claims) {
			c.IsVerified = user.IsActive
			c.IIRID = iirID
			c.CORID = corID
		},
	)
	if err != nil {
		return "", "", "", fmt.Errorf(
			"[AuthService] {Gen Token}: %w",
			err,
		)
	}

	refreshToken, refreshClaims, err := tokens.NewService().
		GenerateSessionToken(
			user.Email,
			user.ID,
			roleIDs,
			authTypeUsed,
			constants.RefreshTokenMaxAge,
			func(c *tokens.Claims) {
				c.IsVerified = user.IsActive
				c.IIRID = iirID
				c.CORID = corID
			},
		)
	if err != nil {
		return "", "", "", fmt.Errorf(
			"[AuthService] {Gen Refresh Token}: %w",
			err,
		)
	}

	// Whitelist session in Redis
	err = s.sessionService.WhitelistSession(
		ctx,
		user.ID,
		accessClaims.ID,
		refreshClaims.ID,
		constants.RefreshTokenMaxAge,
	)
	if err != nil {
		return "", "", "", fmt.Errorf(
			"[AuthService] {Whitelist Session}: %w",
			err,
		)
	}

	// Log successful fallback login
	s.logger.Record(ctx, nil, audit.LogEntry{
		Level:    audit.LevelInfo,
		Category: audit.CategorySecurity,
		Action:   audit.ActionLoginSuccess,
		Message: fmt.Sprintf(
			"User %s successfully logged in via OTP fallback",
			user.Email,
		),
		UserID:    structs.StringToNullableString(user.ID),
		UserEmail: structs.StringToNullableString(user.Email),
		IPAddress: structs.StringToNullableString(ipAddress),
		UserAgent: structs.StringToNullableString(userAgent),
	})

	return user.ID, token, refreshToken, nil
}

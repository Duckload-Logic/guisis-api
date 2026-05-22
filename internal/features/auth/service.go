package auth

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo           *users.Repository
	idpClient      idp.IDPClientInterface
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
		idpClient:      idp.NewIDPClient(),
		redis:          redis,
		sessionService: sessionService,
		emailer:        emailer,
		logger:         logger,
	}
}

// RegisterUser handles native user registration.
func (s *Service) RegisterUser(
	ctx context.Context,
	req RegisterRequest,
) (string, error) {
	// Check cooldown
	cooldownKey := fmt.Sprintf("cooldown:%s", req.Email)
	if val, _ := s.redis.Get(ctx, cooldownKey); val != "" {
		return "", fmt.Errorf(
			"please wait before requesting another verification email",
		)
	}

	// Check if user already exists
	existingUser, _ := s.repo.GetUserByEmail(
		ctx,
		req.Email,
		string(constants.AuthTypeNative),
	)
	if existingUser != nil {
		return "", fmt.Errorf("user already exists: %s", req.Email)
	}

	// Validate password length
	if len(req.Password) < 8 {
		return "", fmt.Errorf("password must be at least 8 characters long")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}

	transactionID := uuid.NewString()

	// Create user object in Domain format
	user := users.User{
		ID:           uuid.NewString(),
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		MiddleName:   structs.StringToNullableString(req.MiddleName),
		SuffixName:   structs.StringToNullableString(req.SuffixName),
		PasswordHash: structs.StringToNullableString(string(hashedPassword)),
		Roles:        []users.Role{{ID: int(constants.DeveloperRoleID)}},
		AuthType:     string(constants.AuthTypeNative),
		IsActive:     false,
	}

	verificationOTP, err := s.get6DigitOTP()
	if err != nil {
		return "", fmt.Errorf("failed to generate verification token: %v", err)
	}
	hashedOTP, err := bcrypt.GenerateFromPassword(
		[]byte(verificationOTP),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate verification token: %v", err)
	}

	// Use a map for storage to include the hidden PasswordHash
	storageData := map[string]interface{}{
		"user":         user,
		"passwordHash": user.PasswordHash,
	}

	userJSON, err := json.Marshal(storageData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal user: %v", err)
	}

	// Store in Redis using the UUID (transactionID)
	val := map[string]string{
		"registrationID":    transactionID,
		"verificationToken": string(hashedOTP),
		"user":              string(userJSON),
	}
	err = s.sessionService.StoreToken(
		ctx,
		sessions.NewJTI(transactionID),
		val,
		300,
	)
	if err != nil {
		return "", fmt.Errorf("failed to store token in redis: %v", err)
	}

	if err := s.validateEmailDomain(req.Email); err != nil {
		return "", err
	}

	err = s.sendVerificationEmail(ctx, verificationOTP, req.Email)
	if err != nil {
		return "", err
	}

	// Set cooldown
	_ = s.redis.Set(ctx, cooldownKey, "1", 120*time.Second)

	return transactionID, nil
}

func (s *Service) ResendVerification(
	ctx context.Context,
	registrationID string,
) error {
	val, err := s.sessionService.GetToken(ctx, sessions.NewJTI(registrationID))
	if err != nil {
		return fmt.Errorf("failed to get token from redis: %v", err)
	}

	// Extract email for cooldown check
	var storageData struct {
		User users.User `json:"user"`
	}
	_ = json.Unmarshal([]byte(val["user"]), &storageData)
	userEmail := storageData.User.Email

	// Check cooldown
	cooldownKey := fmt.Sprintf("cooldown:%s", userEmail)
	if cval, _ := s.redis.Get(ctx, cooldownKey); cval != "" {
		return fmt.Errorf(
			"please wait before requesting another verification email",
		)
	}

	// Resend limit
	var resends int
	if r, ok := val["resends"]; ok {
		fmt.Sscanf(r, "%d", &resends)
	}
	resends++
	val["resends"] = fmt.Sprintf("%d", resends)

	if resends > 3 {
		return fmt.Errorf(
			"too many resends. please try again later or contact support",
		)
	}

	verificationOTP, err := s.get6DigitOTP()
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %v", err)
	}

	hashedOTP, err := bcrypt.GenerateFromPassword(
		[]byte(verificationOTP),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("failed to generate verification token: %v", err)
	}

	val["verificationToken"] = string(hashedOTP)
	err = s.sessionService.StoreToken(
		ctx,
		sessions.NewJTI(registrationID),
		val,
		300,
	)
	if err != nil {
		return fmt.Errorf("failed to store token in redis: %v", err)
	}

	err = json.Unmarshal([]byte(val["user"]), &storageData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user: %v", err)
	}
	user := storageData.User

	if err := s.validateEmailDomain(user.Email); err != nil {
		return err
	}

	err = s.sendVerificationEmail(ctx, verificationOTP, user.Email)
	if err != nil {
		return err
	}

	// Set cooldown
	_ = s.redis.Set(ctx, userEmail, "1", 120*time.Second)

	return nil
}

func (s *Service) get6DigitOTP() (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      constants.ClaimsIssuer,
		AccountName: constants.FromEmail(),
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate verification token: %v", err)
	}

	verificationOTP, err := totp.GenerateCode(key.Secret(), time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to generate verification token: %v", err)
	}

	return verificationOTP, nil
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

func (s *Service) sendVerificationEmail(
	ctx context.Context,
	verificationOTP, receiverEmail string,
) error {
	err := s.emailer.SendOTP(ctx, receiverEmail, verificationOTP)
	if err != nil {
		return fmt.Errorf("failed to send verification email: %v", err)
	}

	return nil
}

func (s *Service) VerifyUser(
	ctx context.Context,
	registrationID string,
	verificationOTP string,
) (string, string, error) {
	// Get user from Redis
	jti := sessions.NewJTI(registrationID)
	val, err := s.sessionService.GetToken(ctx, jti)
	if err != nil {
		return "", "", fmt.Errorf("failed to get token from redis: %v", err)
	}

	// Trials management
	var trials int
	if t, ok := val["trials"]; ok {
		fmt.Sscanf(t, "%d", &trials)
	}

	trials++
	val["trials"] = fmt.Sprintf("%d", trials)

	if trials > 5 {
		_ = s.sessionService.DeleteToken(ctx, jti)
		return "", "", fmt.Errorf(
			"too many verification attempts. registration expired",
		)
	}

	// Update trials in Redis
	_ = s.sessionService.StoreToken(ctx, jti, val, 300)

	// Verify OTP
	err = bcrypt.CompareHashAndPassword(
		[]byte(val["verificationToken"]),
		[]byte(verificationOTP),
	)
	if err != nil {
		return "", "", fmt.Errorf("invalid credentials: %v", err)
	}

	var storageData struct {
		User         users.User             `json:"user"`
		PasswordHash structs.NullableString `json:"passwordHash"`
	}

	err = json.Unmarshal([]byte(val["user"]), &storageData)
	if err != nil {
		return "", "", fmt.Errorf("failed to unmarshal user data: %v", err)
	}

	user := storageData.User
	user.PasswordHash = storageData.PasswordHash

	user.IsActive = true

	// Run in transaction
	err = s.repo.WithTransaction(
		ctx,
		func(tx datastore.DB) error {
			return s.repo.CreateUser(ctx, tx, user)
		},
	)
	if err != nil {
		return "", "", fmt.Errorf("failed to create user: %v", err)
	}

	return user.ID, user.Email, nil
}

// AuthenticateUser handles native email/password authentication.
func (s *Service) AuthenticateUser(
	ctx context.Context, email, password, ipAddress, userAgent string,
) (string, string, string, error) {
	// Check lockout
	lockoutKey := fmt.Sprintf("lockout:%s", email)
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
			Level:     audit.LevelWarning,
			Category:  audit.CategorySecurity,
			Action:    audit.ActionLoginFailed,
			Message:   fmt.Sprintf("Failed login attempt (Inactive): %s", email),
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

	failureKey := fmt.Sprintf("failed_attempts:%s", email)
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

		if failures >= 5 {
			err = s.redis.Set(ctx, lockoutKey, "true", 15*time.Minute)
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
			15*time.Minute, // 15 minutes lockout period
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

	// Generate the token
	roleIDs := make([]int, len(user.Roles))
	for i, r := range user.Roles {
		roleIDs[i] = r.ID
	}

	token, claims, err := tokens.NewService().GenerateToken(
		user.Email,
		user.ID,
		roleIDs,
		string(constants.AuthTypeNative),
		constants.AccessTokenMaxAge,
	)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate session: %v", err)
	}

	// Generate refresh token
	refreshToken, _, err := tokens.NewService().GenerateToken(
		user.Email,
		user.ID,
		roleIDs,
		string(constants.AuthTypeNative),
		constants.RefreshTokenMaxAge,
	)
	if err != nil {
		return "", "", "", fmt.Errorf(
			"failed to generate refresh token: %v", err,
		)
	}

	// Store in Redis using the Token ID (jti)
	val := map[string]string{
		"userID":          user.ID,
		"tokenType":       string(constants.AuthTypeNative),
		"appRefreshToken": refreshToken,
		"ipAddress":       ipAddress,
		"userAgent":       userAgent,
	}
	err = s.sessionService.StoreUserToken(
		ctx,
		user.ID,
		sessions.NewJTI(claims.ID),
		val,
		constants.RefreshTokenMaxAge,
	)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to store token in redis: %v", err)
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
	accessTokenJTI sessions.JTIDTO,
	cfg *config.Config,
	ipAddress, userAgent string,
) (string, string, error) {
	// Get session from Redis
	session, err := s.sessionService.GetToken(ctx, accessTokenJTI)
	if err != nil {
		return "", "", fmt.Errorf("session expired or invalid: %v", err)
	}

	refreshToken := session["appRefreshToken"]
	if refreshToken == "" {
		return "", "", fmt.Errorf("refresh token missing from session: %v", err)
	}

	// Validate App Refresh Token
	claims, err := tokens.NewService().ValidateToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %v", err)
	}

	// Check token type
	if claims.TokenType == string(constants.AuthTypeIDP) {
		// Get IDP refresh token from Redis
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

		// Generate NEW App Tokens
		newAppAccessToken, accessClaims, err := tokens.NewService().
			GenerateToken(
				claims.UserEmail,
				claims.UserID,
				claims.RoleIDs,
				string(constants.AuthTypeIDP),
				constants.AccessTokenMaxAge,
			)
		if err != nil {
			return "", "", err
		}

		newAppRefreshToken, refreshClaims, err := tokens.NewService().
			GenerateToken(
				claims.UserEmail,
				claims.UserID,
				claims.RoleIDs,
				string(constants.AuthTypeIDP),
				constants.RefreshTokenMaxAge,
			)
		if err != nil {
			return "", "", err
		}

		// Update Redis: App Access session
		val := map[string]string{
			"userID":          claims.UserID,
			"tokenType":       string(constants.AuthTypeIDP),
			"appRefreshToken": newAppRefreshToken,
			"idpAccessToken":  tokenResp.AccessToken,
		}
		err = s.sessionService.StoreToken(
			ctx,
			sessions.NewJTI(accessClaims.ID),
			val,
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

		// Clean up OLD session and IDP refresh keys
		_ = s.sessionService.DeleteUserToken(ctx, claims.UserID, accessTokenJTI)
		_ = s.redis.Del(ctx, idpRefreshKey)

		return newAppAccessToken, newAppRefreshToken, nil
	}

	// Native flow
	newToken, newClaims, err := tokens.NewService().GenerateToken(
		claims.UserEmail,
		claims.UserID,
		claims.RoleIDs,
		string(constants.AuthTypeNative),
		constants.AccessTokenMaxAge,
	)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, _, err := tokens.NewService().GenerateToken(
		claims.UserEmail,
		claims.UserID,
		claims.RoleIDs,
		string(constants.AuthTypeNative),
		constants.RefreshTokenMaxAge,
	)
	if err != nil {
		return "", "", err
	}

	// Update Redis using new jti
	val := map[string]string{
		"userID":          claims.UserID,
		"tokenType":       string(constants.AuthTypeNative),
		"appRefreshToken": newRefreshToken,
		"ipAddress":       ipAddress,
		"userAgent":       userAgent,
	}
	err = s.sessionService.StoreUserToken(
		ctx,
		claims.UserID,
		sessions.NewJTI(newClaims.ID),
		val,
		constants.RefreshTokenMaxAge,
	)
	if err != nil {
		return "", "", err
	}

	// Clean up OLD session key
	_ = s.sessionService.DeleteUserToken(ctx, claims.UserID, accessTokenJTI)

	return newToken, newRefreshToken, nil
}

// RefreshIDPToken handles token refresh for IDP-authenticated sessions.
func (s *Service) RefreshIDPToken(
	ctx context.Context, refreshToken string, cfg *config.Config,
) (string, string, error) {
	// Call IDP refresh endpoint
	tokenResp, err := s.idpClient.RefreshToken(ctx, refreshToken, cfg)
	if err != nil {
		return "", "", fmt.Errorf("[AuthService] {IDP Refresh}: %w", err)
	}

	return tokenResp.AccessToken, tokenResp.RefreshToken, nil
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
	tokenType string,
	cfg *config.Config,
) (string, error) {
	// Identify the session using the Access Token JTI
	claims, err := tokens.NewService().ParseTokenUnverified(token)
	if err != nil {
		log.Printf("[AuthService:Logout] {Parse Error}: %v", err)
		return "", nil // Move on since we can't identify the session
	}
	accessJTI := claims.ID

	// Fetch the session data to find linked refresh tokens
	sessionData, _ := s.sessionService.GetToken(ctx, sessions.NewJTI(accessJTI))
	var idpToken string
	if sessionData != nil {
		idpToken = sessionData["idpAccessToken"]
	}

	// Delete any linked IDP refresh tokens
	if sessionData != nil {
		if appRefreshToken := sessionData["appRefreshToken"]; appRefreshToken != "" {
			// Get refresh token claims to identify IDP refresh key
			rClaims, err := tokens.NewService().
				ParseTokenUnverified(appRefreshToken)
			if err == nil {
				idpKey := sessions.NewJTI(rClaims.ID).ToIDPRefreshKey()
				_ = s.redis.Del(ctx, idpKey)
			}
		}
	}

	// Delete the primary session key
	if userID := claims.UserID; userID != "" {
		if err := s.sessionService.DeleteUserToken(
			ctx,
			userID,
			sessions.NewJTI(accessJTI),
		); err != nil {
			log.Printf("[AuthService:Logout] {Redis Error}: %v", err)
		}
	} else {
		// Fallback for M2M or cases where userID missing
		if err := s.sessionService.DeleteToken(
			ctx,
			sessions.NewJTI(accessJTI),
		); err != nil {
			log.Printf("[AuthService:Logout] {Redis Error}: %v", err)
		}
	}

	userInfo, _ := s.idpClient.GetUserInfo(ctx, idpToken, cfg)

	// Construct logout URL for front-channel redirect
	if tokenType == string(constants.AuthTypeIDP) && userInfo != nil {
		return s.idpClient.GetLogoutURL(cfg, userInfo.ID), nil
	}

	// Fallback redirect for native logout or incomplete IDP sessions
	return "/", nil
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
) (string, string, string, string, string, error) {
	// Exchange authorization code for IDP tokens
	tokenResp, err := s.idpClient.ExchangeCodeForToken(ctx, code, cfg)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf(
			"[AuthService] {Token Exchange}: %w",
			err,
		)
	}

	// Fetch User Info from IDP
	userInfo, err := s.GetIDPUserInfo(ctx, tokenResp.AccessToken, cfg)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf(
			"[AuthService] {Fetch User Info}: %w",
			err,
		)
	}

	// Parse Tokens
	idpAccessToken := tokenResp.AccessToken
	idpRefreshToken := tokenResp.RefreshToken

	// Whitelist Check: Authoritative role source for IDP users.
	// We check this on every login to support dynamic role changes (promotions).
	whitelistRoleIDs, whitelistErr := s.repo.CheckUserWhitelist(
		ctx,
		userInfo.Email,
	)

	// User Existence Check
	localUser, err := s.repo.GetUserByEmail(
		ctx,
		userInfo.Email,
		string(constants.AuthTypeIDP),
	)

	// Determine target roles from whitelist or defaults
	var targetRoleIDs []int
	if whitelistErr != nil {
		return "", "", "", "", "", fmt.Errorf(
			"[AuthService] {Whitelist Check}: %w",
			whitelistErr,
		)
	}

	if len(whitelistRoleIDs) > 0 {
		targetRoleIDs = whitelistRoleIDs
	} else if err == sql.ErrNoRows {
		// New users not in whitelist default to Student.
		targetRoleIDs = []int{int(constants.StudentRoleID)}
	}

	switch err {
	case sql.ErrNoRows:
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
			SuffixName:   structs.StringToNullableString(userInfo.SuffixName),
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
			return "", "", "", "", "", fmt.Errorf(
				"[AuthService] {Provision IDP User}: %w",
				err,
			)
		}
	case nil:
		// User exists. Sync roles and consume whitelist if present
		if len(whitelistRoleIDs) > 0 && len(targetRoleIDs) > 0 {
			addedAny := false
			err = s.repo.WithTransaction(
				ctx,
				func(tx datastore.DB) error {
					for _, targetID := range targetRoleIDs {
						hasRole := false
						for _, r := range localUser.Roles {
							if r.ID == targetID {
								hasRole = true
								break
							}
						}

						if !hasRole {
							addedAny = true
							// Promotion or sync from whitelist
							err = s.repo.AssignRole(
								ctx,
								tx,
								users.RoleAssignment{
									UserID: localUser.ID,
									RoleID: targetID,
									Reason: structs.StringToNullableString(
										"Whitelist synchronization",
									),
									AssignedBy: structs.StringToNullableString(
										constants.SystemEntityType,
									),
								},
							)
							if err != nil {
								return err
							}
						}
					}
					return s.repo.RemoveUserFromWhitelist(
						ctx,
						tx,
						userInfo.Email,
					)
				},
			)
			if err != nil {
				return "", "", "", "", "", fmt.Errorf(
					"[AuthService] {Sync Roles & Clear Whitelist}: %w",
					err,
				)
			}

			if addedAny {
				// Refresh roles
				updatedRoles, _ := s.repo.GetRolesByUserID(ctx, localUser.ID)
				localUser.Roles = updatedRoles
			}
		}
	default:
		return "", "", "", "", "", fmt.Errorf(
			"[AuthService] {Anchor Check}: %w",
			err,
		)
	}

	appUserID := localUser.ID

	roleName := "student"
	if len(localUser.Roles) > 0 {
		roleName = strings.ToLower(localUser.Roles[0].Name)
	}

	roleIDs := make([]int, len(localUser.Roles))
	for i, r := range localUser.Roles {
		roleIDs[i] = r.ID
	}

	appAccessToken, accessClaims, err := tokens.NewService().GenerateToken(
		userInfo.Email,
		appUserID,
		roleIDs,
		string(constants.AuthTypeIDP),
		constants.AccessTokenMaxAge,
	)
	if err != nil {
		return "", "", "", "", "", err
	}

	appRefreshToken, refreshClaims, err := tokens.NewService().GenerateToken(
		userInfo.Email,
		appUserID,
		roleIDs,
		string(constants.AuthTypeIDP),
		constants.RefreshTokenMaxAge,
	)
	if err != nil {
		return "", "", "", "", "", err
	}

	val := map[string]string{
		"userID":          appUserID,
		"tokenType":       string(constants.AuthTypeIDP),
		"appRefreshToken": appRefreshToken,
		"idpAccessToken":  idpAccessToken,
		"ipAddress":       ipAddress,
		"userAgent":       userAgent,
	}
	if err := s.sessionService.StoreUserToken(
		ctx,
		appUserID,
		sessions.NewJTI(accessClaims.ID),
		val,
		constants.RefreshTokenMaxAge,
	); err != nil {
		return "", "", "", "", "", err
	}

	idpRefreshKey := sessions.NewJTI(refreshClaims.ID).ToIDPRefreshKey()
	err = s.redis.Set(
		ctx,
		idpRefreshKey,
		idpRefreshToken,
		time.Duration(constants.RefreshTokenMaxAge)*time.Second,
	)
	if err != nil {
		return "", "", "", "", "", err
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
		UserID:    structs.StringToNullableString(appUserID),
		UserEmail: structs.StringToNullableString(userInfo.Email),
		IPAddress: structs.StringToNullableString(ipAddress),
		UserAgent: structs.StringToNullableString(userAgent),
	})

	return appAccessToken, appRefreshToken, appUserID, userInfo.Email,
		roleName, nil
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

func (s *Service) BlockUser(ctx context.Context, userID string) error {
	return s.repo.WithTransaction(
		ctx,
		func(tx datastore.DB) error {
			return s.repo.BlockUser(ctx, tx, userID)
		},
	)
}

func (s *Service) UnblockUser(ctx context.Context, userID string) error {
	return s.repo.WithTransaction(
		ctx,
		func(tx datastore.DB) error {
			return s.repo.UnblockUser(ctx, tx, userID)
		},
	)
}

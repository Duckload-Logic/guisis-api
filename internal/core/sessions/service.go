package sessions

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/olazo-johnalbert/duckload-api/internal/core/constants"
	"github.com/olazo-johnalbert/duckload-api/internal/infrastructure/datastore"
)

type Service struct {
	redis *datastore.RedisClient
}

func NewService(redis *datastore.RedisClient) *Service {
	return &Service{redis: redis}
}

// StoreToken saves session data in Redis with a JTI-based key.
func (s *Service) StoreToken(
	ctx context.Context,
	jti JTIDTO,
	data map[string]string,
	expireSeconds int,
) error {
	key := jti.ToSessionKey()
	valJSON, _ := json.Marshal(data)

	err := s.redis.Set(
		ctx,
		key,
		string(valJSON),
		time.Duration(expireSeconds)*time.Second,
	)
	if err != nil {
		return fmt.Errorf("failed to store token in redis: %w", err)
	}

	return nil
}

// GetToken retrieves session data from Redis.
func (s *Service) GetToken(
	ctx context.Context,
	jti JTIDTO,
) (map[string]string, error) {
	key := jti.ToSessionKey()
	val, err := s.redis.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("session not found or expired: %w", err)
	}

	var data map[string]string
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return nil, fmt.Errorf("failed to parse session data: %w", err)
	}

	return data, nil
}

// DeleteToken removes session data from Redis.
func (s *Service) DeleteToken(ctx context.Context, jti JTIDTO) error {
	key := jti.ToSessionKey()
	return s.redis.Del(ctx, key)
}

// StoreUserToken saves session data and links it to a user for auditing.
func (s *Service) StoreUserToken(
	ctx context.Context,
	userID string,
	jti JTIDTO,
	data map[string]string,
	expireSeconds int,
) error {
	// Store the session data
	if err := s.StoreToken(ctx, jti, data, expireSeconds); err != nil {
		return err
	}

	// Link to user sessions set
	userKey := ToUserSessionsKey(userID)
	err := s.redis.SAdd(ctx, userKey, jti.Value)
	if err != nil {
		return fmt.Errorf("failed to link session to user: %w", err)
	}

	// Set expiration on the set if it's new (or refresh it)
	// We use the same expiration as the token for simplicity
	s.redis.Expire(ctx, userKey, time.Duration(expireSeconds)*time.Second)

	return nil
}
// WhitelistSession stores active JTIs in a Redis hash.
func (s *Service) WhitelistSession(
	ctx context.Context,
	userID string,
	accessJTI string,
	refreshJTI string,
	expireSeconds int,
) error {
	key := fmt.Sprintf("%s%s", constants.RedisUserSessionKeyPrefix, userID)

	err := s.redis.HSet(
		ctx,
		key,
		constants.RedisSessionAccessJTIField,
		accessJTI,
		constants.RedisSessionRefreshJTIField,
		refreshJTI,
	)
	if err != nil {
		return fmt.Errorf("failed to store whitelist session: %w", err)
	}

	err = s.redis.Expire(ctx, key, time.Duration(expireSeconds)*time.Second)
	if err != nil {
		return fmt.Errorf("failed to set whitelist expiration: %w", err)
	}

	return nil
}

// RevokeUserSession deletes the whitelist key for a user.
func (s *Service) RevokeUserSession(
	ctx context.Context,
	userID string,
) error {
	key := fmt.Sprintf("%s%s", constants.RedisUserSessionKeyPrefix, userID)
	return s.redis.Del(ctx, key)
}

// GetWhitelistedSession retrieves the whitelisted session for a user.
func (s *Service) GetWhitelistedSession(
	ctx context.Context,
	userID string,
) (map[string]string, error) {
	key := fmt.Sprintf("%s%s", constants.RedisUserSessionKeyPrefix, userID)
	return s.redis.HGetAll(ctx, key)
}

// DeleteUserToken deletes the user's whitelisted session from Redis.
func (s *Service) DeleteUserToken(
	ctx context.Context,
	userID string,
	jti JTIDTO,
) error {
	return s.RevokeUserSession(ctx, userID)
}

// ListUserSessions returns the whitelisted session data for a user.
func (s *Service) ListUserSessions(
	ctx context.Context,
	userID string,
) ([]map[string]string, error) {
	sessionData, err := s.GetWhitelistedSession(ctx, userID)
	if err != nil || len(sessionData) == 0 {
		return []map[string]string{}, nil
	}

	res := []map[string]string{
		{
			"jti":    sessionData[constants.RedisSessionAccessJTIField],
			"type":   "access",
			"status": "active",
			"userID": userID,
		},
		{
			"jti":    sessionData[constants.RedisSessionRefreshJTIField],
			"type":   "refresh",
			"status": "active",
			"userID": userID,
		},
	}
	return res, nil
}

// RevokeAllUserSessions invalidates all active sessions for a user.
func (s *Service) RevokeAllUserSessions(
	ctx context.Context,
	userID string,
) error {
	return s.RevokeUserSession(ctx, userID)
}

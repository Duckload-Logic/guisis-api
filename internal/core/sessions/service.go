package sessions

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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

// DeleteUserToken removes a session and its link to the user.
func (s *Service) DeleteUserToken(
	ctx context.Context,
	userID string,
	jti JTIDTO,
) error {
	// Unlink from user
	userKey := ToUserSessionsKey(userID)
	s.redis.SRem(ctx, userKey, jti.Value)

	// Delete session data
	return s.DeleteToken(ctx, jti)
}

// ListUserSessions returns all active session data for a user.
func (s *Service) ListUserSessions(
	ctx context.Context,
	userID string,
) ([]map[string]string, error) {
	userKey := ToUserSessionsKey(userID)
	jtis, err := s.redis.SMembers(ctx, userKey)
	if err != nil {
		return nil, fmt.Errorf("failed to list user sessions: %w", err)
	}

	sessions := make([]map[string]string, 0, len(jtis))
	for _, jtiVal := range jtis {
		jti := NewJTI(jtiVal)
		data, err := s.GetToken(ctx, jti)
		if err != nil {
			// Session might have expired individually, clean up the set
			s.redis.SRem(ctx, userKey, jtiVal)
			continue
		}
		// Add JTI to the data map for the frontend
		data["jti"] = jtiVal
		sessions = append(sessions, data)
	}

	return sessions, nil
}

// RevokeAllUserSessions invalidates all active sessions for a user.
func (s *Service) RevokeAllUserSessions(
	ctx context.Context,
	userID string,
) error {
	userKey := ToUserSessionsKey(userID)
	jtis, err := s.redis.SMembers(ctx, userKey)
	if err != nil {
		return fmt.Errorf("failed to get user sessions: %w", err)
	}

	for _, jtiVal := range jtis {
		_ = s.DeleteToken(ctx, NewJTI(jtiVal))
	}

	return s.redis.Del(ctx, userKey)
}

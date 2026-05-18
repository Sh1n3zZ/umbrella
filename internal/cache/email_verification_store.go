// Package cache hosts Redis-backed implementations of domain interfaces.
package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/Sh1n3zZ/umbrella/domain"
	"github.com/Sh1n3zZ/umbrella/internal/cache/cachekey"
)

// EmailVerificationStore is the Redis-backed implementation of
// domain.EmailVerificationStore. Tokens are stored under the
// cachekey.OAuthEmailVerification namespace and consumed atomically via
// GETDEL (Redis 6.2+).
type EmailVerificationStore struct {
	client *redis.Client
}

// NewEmailVerificationStore returns a store bound to the given Redis client.
func NewEmailVerificationStore(client *redis.Client) *EmailVerificationStore {
	return &EmailVerificationStore{client: client}
}

// compile-time assertion that the implementation satisfies the domain interface.
var _ domain.EmailVerificationStore = (*EmailVerificationStore)(nil)

// Save records token -> userID with the given TTL.
func (s *EmailVerificationStore) Save(
	ctx context.Context,
	token string,
	userID uuid.UUID,
	ttl time.Duration,
) error {
	if s == nil || s.client == nil {
		return errors.New("cache: email verification store is not initialised")
	}
	if token == "" {
		return errors.New("cache: token is required")
	}
	if userID == uuid.Nil {
		return errors.New("cache: user id is required")
	}
	if ttl <= 0 {
		return errors.New("cache: ttl must be positive")
	}

	key := cachekey.OAuthEmailVerification.Key(token)
	if err := s.client.Set(ctx, key, userID.String(), ttl).Err(); err != nil {
		return fmt.Errorf("cache: save verification token: %w", err)
	}
	return nil
}

// Consume atomically reads and removes the token, returning the associated
// user id. Missing or expired tokens surface as
// domain.ErrEmailVerificationTokenInvalid.
func (s *EmailVerificationStore) Consume(
	ctx context.Context,
	token string,
) (uuid.UUID, error) {
	if s == nil || s.client == nil {
		return uuid.Nil, errors.New("cache: email verification store is not initialised")
	}
	if token == "" {
		return uuid.Nil, domain.ErrEmailVerificationTokenInvalid
	}

	key := cachekey.OAuthEmailVerification.Key(token)
	raw, err := s.client.GetDel(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return uuid.Nil, domain.ErrEmailVerificationTokenInvalid
		}
		return uuid.Nil, fmt.Errorf("cache: consume verification token: %w", err)
	}

	userID, err := uuid.Parse(raw)
	if err != nil {
		// Stored value is corrupted; treat as invalid token rather than leak details.
		return uuid.Nil, domain.ErrEmailVerificationTokenInvalid
	}
	return userID, nil
}

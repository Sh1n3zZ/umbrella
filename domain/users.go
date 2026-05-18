package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// User is an OAuth end-user account stored in oauth_users.
type User struct {
	UserID        uuid.UUID
	Nickname      string
	Email         string
	Password      string
	IsActive      bool
	EmailVerified bool
	CreatedAt     time.Time
}

// UserResponse is the public representation of an OAuth user (password omitted).
type UserResponse struct {
	UserID        string `json:"user_id"`
	Nickname      string `json:"nickname"`
	Email         string `json:"email"`
	IsActive      bool   `json:"is_active"`
	EmailVerified bool   `json:"email_verified"`
	CreatedAt     string `json:"created_at"`
}

// Email verification is part of the oauth_users registration flow. The
// store interface lives here rather than in a dedicated domain file because
// it has no meaning outside of the User aggregate.
const (
	// EmailVerificationTTL is how long an issued verification token stays valid.
	EmailVerificationTTL = 24 * time.Hour
	// EmailVerificationTokenBytes is the random entropy budget for verification tokens
	// before URL-safe base64 encoding (see tokenutil.GenerateURLSafeRandomString).
	EmailVerificationTokenBytes = 32
)

// ErrEmailVerificationTokenInvalid is returned when a verification token is
// missing, expired, or already consumed. It deliberately conflates the cases
// so callers cannot probe for token existence.
var ErrEmailVerificationTokenInvalid = errors.New("verification token is invalid or expired")

// ErrOauthUserEmailTaken signals that registration was rejected because the
// email is already associated with an existing oauth_users row.
var ErrOauthUserEmailTaken = errors.New("email already registered")

// EmailVerificationStore persists single-use verification tokens issued during
// registration. Implementations MUST consume tokens atomically so a token can
// be used at most once.
type EmailVerificationStore interface {
	// Save records token -> userID with the given TTL. Existing tokens with
	// the same key are overwritten.
	Save(ctx context.Context, token string, userID uuid.UUID, ttl time.Duration) error
	// Consume atomically reads and deletes the token. It returns
	// ErrEmailVerificationTokenInvalid when the token is missing or expired.
	Consume(ctx context.Context, token string) (uuid.UUID, error)
}

// OauthUsersRepository is the persistence contract the usecase depends on.
// Implementations live under the repository package.
type OauthUsersRepository interface {
	// CreateOauthUser inserts a new user and returns the row including
	// DB-assigned fields. It returns pgx.ErrNoRows when an ON CONFLICT clause
	// suppresses the insert (i.e. the email is already registered).
	CreateOauthUser(ctx context.Context, nickname, email, password string) (User, error)
	// MarkOauthUserEmailVerified sets email_verified to true for the given user
	// and returns the updated row.
	MarkOauthUserEmailVerified(ctx context.Context, userID uuid.UUID) (User, error)
}

// OauthUsersUsecase is the business-rules contract the HTTP controller
// depends on. Implementations live under the usecase package.
type OauthUsersUsecase interface {
	// RegisterOauthUser provisions a new account, dispatches the verification
	// email, and returns the public user representation. It returns
	// ErrOauthUserEmailTaken when the email is already registered.
	RegisterOauthUser(ctx context.Context, nickname, email, password string) (UserResponse, error)
	// VerifyEmail consumes a verification token and marks the matching user
	// as email-verified. It returns ErrEmailVerificationTokenInvalid for
	// missing, expired, or already-used tokens.
	VerifyEmail(ctx context.Context, token string) (UserResponse, error)
}

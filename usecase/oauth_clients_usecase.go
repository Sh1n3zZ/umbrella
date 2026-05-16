package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Sh1n3zZ/umbrella/repository"
)

// OAuth client / redirect resolution errors for the authorization endpoint.
var (
	ErrOauthUnknownClient    = errors.New("oauth: unknown client_id")
	ErrOauthRedirectMismatch = errors.New("oauth: redirect_uri is not registered for this client")
	ErrOauthRedirectRequired = errors.New("oauth: redirect_uri is required when zero or multiple URIs are registered")
)

const (
	// n in GenerateURLSafeRandomString is random byte count before encoding.
	randomEntropyMinBytes = 1
	randomEntropyMaxBytes = 4096
)

// OauthClientsUsecase coordinates oauth_clients-related business rules.
type OauthClientsUsecase struct {
	oauthClientsRepo *repository.OauthClientsRepository
}

// NewOauthClientsUsecase builds a usecase backed by the given repository.
func NewOauthClientsUsecase(oauthClientsRepo *repository.OauthClientsRepository) *OauthClientsUsecase {
	return &OauthClientsUsecase{oauthClientsRepo: oauthClientsRepo}
}

// EffectiveRedirectURI returns the redirection URI for the authorization response (RFC 6749
// Section 3.1.2.3). If requested is non-empty it must exactly match a registered URI. If empty,
// exactly one registered URI must exist.
func (u *OauthClientsUsecase) EffectiveRedirectURI(ctx context.Context, clientID pgtype.UUID, requested string) (string, error) {
	if !clientID.Valid {
		return "", ErrOauthRedirectRequired
	}

	uris, err := u.oauthClientsRepo.GetOauthClientRedirectUrisByClientID(ctx, clientID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrOauthUnknownClient
		}
		return "", err
	}

	if requested != "" {
		for _, ru := range uris {
			if ru == requested {
				return requested, nil
			}
		}
		return "", ErrOauthRedirectMismatch
	}

	switch len(uris) {
	case 1:
		return uris[0], nil
	default:
		return "", ErrOauthRedirectRequired
	}
}

// ValidateClientRedirectURI reports whether redirectURI is registered for clientID.
// It returns (false, nil) when client_id is unknown, redirectURI is empty, clientID is invalid,
// or the URI is not among redirect_uris; (false, err) only for persistence failures.
func (u *OauthClientsUsecase) ValidateClientRedirectURI(ctx context.Context, clientID pgtype.UUID, redirectURI string) (bool, error) {
	if redirectURI == "" || !clientID.Valid {
		return false, nil
	}

	registered, err := u.oauthClientsRepo.GetOauthClientRedirectUrisByClientID(ctx, clientID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	for _, registeredURI := range registered {
		if registeredURI == redirectURI {
			return true, nil
		}
	}
	return false, nil
}

// GenerateURLSafeRandomString reads n cryptographically random bytes and returns
// base64.URLEncoding of that buffer (RFC 4648 URL-safe alphabet: '-' and '_' instead
// of '+' and '/'). Output uses only [A-Za-z0-9-_=], which is a subset of RFC 6749
// VSCHAR for values like "state" (Appendix A.5), "code" (A.11), and tokens (A.12, A.17).
// Returned string length is base64.URLEncoding.EncodedLen(n).
func (u *OauthClientsUsecase) GenerateURLSafeRandomString(n int) (string, error) {
	if n < randomEntropyMinBytes || n > randomEntropyMaxBytes {
		return "", fmt.Errorf("oauth: random entropy byte count must be between %d and %d", randomEntropyMinBytes, randomEntropyMaxBytes)
	}

	raw := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, raw); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(raw), nil
}

package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Sh1n3zZ/umbrella/domain"
)

// OauthClientsUsecase coordinates oauth_clients-related business rules.
type OauthClientsUsecase struct {
	oauthClientsRepo domain.OauthClientsRepository
}

var _ domain.OauthClientsUsecase = (*OauthClientsUsecase)(nil)

// NewOauthClientsUsecase builds a usecase backed by the given repository.
func NewOauthClientsUsecase(oauthClientsRepo domain.OauthClientsRepository) *OauthClientsUsecase {
	return &OauthClientsUsecase{oauthClientsRepo: oauthClientsRepo}
}

// EffectiveRedirectURI returns the redirection URI for the authorization response (RFC 6749
// Section 3.1.2.3). If requested is non-empty it must exactly match a registered URI. If empty,
// exactly one registered URI must exist.
func (u *OauthClientsUsecase) EffectiveRedirectURI(ctx context.Context, clientID uuid.UUID, requested string) (string, error) {
	if clientID == uuid.Nil {
		return "", domain.ErrOauthRedirectRequired
	}

	uris, err := u.oauthClientsRepo.GetOauthClientRedirectUrisByClientID(ctx, clientID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", domain.ErrOauthUnknownClient
		}
		return "", err
	}

	if requested != "" {
		for _, ru := range uris {
			if ru == requested {
				return requested, nil
			}
		}
		return "", domain.ErrOauthRedirectMismatch
	}

	switch len(uris) {
	case 1:
		return uris[0], nil
	default:
		return "", domain.ErrOauthRedirectRequired
	}
}

// ValidateClientRedirectURI reports whether redirectURI is registered for clientID.
// It returns (false, nil) when client_id is unknown, redirectURI is empty, clientID is invalid,
// or the URI is not among redirect_uris; (false, err) only for persistence failures.
func (u *OauthClientsUsecase) ValidateClientRedirectURI(ctx context.Context, clientID uuid.UUID, redirectURI string) (bool, error) {
	if redirectURI == "" || clientID == uuid.Nil {
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

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Sh1n3zZ/umbrella/domain"
)

// OauthClientsRepository handles persistence for oauth_clients.
type OauthClientsRepository struct {
	*BaseRepository
}

var _ domain.OauthClientsRepository = (*OauthClientsRepository)(nil)

// NewOauthClientsRepository returns a repository using the shared base (pool + sqlc session).
func NewOauthClientsRepository(base *BaseRepository) *OauthClientsRepository {
	return &OauthClientsRepository{BaseRepository: base}
}

// GetOauthClientRedirectUrisByClientID returns redirect_uris for the row matching client_id.
// If no row exists, the error is pgx.ErrNoRows (same as sqlc :one semantics).
func (r *OauthClientsRepository) GetOauthClientRedirectUrisByClientID(
	ctx context.Context,
	clientID uuid.UUID,
) ([]string, error) {
	return r.GetQueries(ctx).GetOauthClientRedirectUrisByClientID(
		ctx, pgtype.UUID{Bytes: clientID, Valid: true},
	)
}

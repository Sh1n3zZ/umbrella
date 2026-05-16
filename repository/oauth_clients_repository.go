package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Sh1n3zZ/umbrella/internal/sqlc"
)

// OauthClientsRepository handles persistence for oauth_clients.
type OauthClientsRepository struct {
	*BaseRepository
}

// NewOauthClientsRepository returns a repository using the shared base (pool + sqlc session).
func NewOauthClientsRepository(base *BaseRepository) *OauthClientsRepository {
	return &OauthClientsRepository{BaseRepository: base}
}

// CreateOauthClient inserts one row and returns it with DB-assigned fields (e.g. id, create_time).
func (r *OauthClientsRepository) CreateOauthClient(ctx context.Context, params sqlc.CreateOauthClientParams) (sqlc.OauthClient, error) {
	return r.GetQueries(ctx).CreateOauthClient(ctx, params)
}

// OauthClientExistsByClientID returns whether a row exists for the given client_id.
func (r *OauthClientsRepository) OauthClientExistsByClientID(ctx context.Context, clientID pgtype.UUID) (bool, error) {
	return r.GetQueries(ctx).OauthClientExistsByClientID(ctx, clientID)
}

// GetOauthClientRedirectUrisByClientID returns redirect_uris for the row matching client_id.
// If no row exists, the error is pgx.ErrNoRows (same as sqlc :one semantics).
func (r *OauthClientsRepository) GetOauthClientRedirectUrisByClientID(ctx context.Context, clientID pgtype.UUID) ([]string, error) {
	return r.GetQueries(ctx).GetOauthClientRedirectUrisByClientID(ctx, clientID)
}

package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Sh1n3zZ/umbrella/domain"
	"github.com/Sh1n3zZ/umbrella/internal/sqlc"
)

// OauthUsersRepository handles persistence for oauth_users.
type OauthUsersRepository struct {
	*BaseRepository
}

var _ domain.OauthUsersRepository = (*OauthUsersRepository)(nil)

// NewOauthUsersRepository returns a repository using the shared base (pool + sqlc session).
func NewOauthUsersRepository(base *BaseRepository) *OauthUsersRepository {
	return &OauthUsersRepository{BaseRepository: base}
}

// CreateOauthUser inserts one row and returns it with DB-assigned fields (e.g. user_id, created_at).
func (r *OauthUsersRepository) CreateOauthUser(
	ctx context.Context,
	nickname,
	email,
	password string,
) (domain.User, error) {
	row, err := r.GetQueries(ctx).CreateOauthUser(ctx, sqlc.CreateOauthUserParams{
		Nickname: nickname,
		Email:    email,
		Password: password,
	})
	if err != nil {
		return domain.User{}, err
	}
	return toUserDomain(row)
}

// MarkOauthUserEmailVerified sets email_verified to true for the given user_id and returns the updated row.
func (r *OauthUsersRepository) MarkOauthUserEmailVerified(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	row, err := r.GetQueries(ctx).MarkOauthUserEmailVerified(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return domain.User{}, err
	}
	return toUserDomain(row)
}

func toUserDomain(row sqlc.OauthUser) (domain.User, error) {
	user := domain.User{
		Nickname:      row.Nickname,
		Email:         row.Email,
		Password:      row.Password,
		IsActive:      row.IsActive,
		EmailVerified: row.EmailVerified,
	}
	if row.UserID.Valid {
		id, err := uuid.FromBytes(row.UserID.Bytes[:])
		if err != nil {
			return domain.User{}, err
		}
		user.UserID = id
	}
	if row.CreatedAt.Valid {
		user.CreatedAt = row.CreatedAt.Time
	}
	return user, nil
}

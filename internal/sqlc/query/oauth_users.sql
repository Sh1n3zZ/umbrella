-- name: CreateOauthUser :one
INSERT INTO oauth_users (
    nickname,
    email,
    password
) VALUES (
    $1,
    $2,
    $3
)
ON CONFLICT (email) DO NOTHING
RETURNING
    user_id,
    nickname,
    email,
    password,
    is_active,
    email_verified,
    created_at;

-- name: MarkOauthUserEmailVerified :one
UPDATE oauth_users
SET email_verified = true
WHERE user_id = $1
RETURNING
    user_id,
    nickname,
    email,
    password,
    is_active,
    email_verified,
    created_at;

-- +goose Up
CREATE TABLE oauth_users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nickname TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    email_verified BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX oauth_users_email_key ON oauth_users (email);

-- +goose Down
DROP TABLE IF EXISTS oauth_users;

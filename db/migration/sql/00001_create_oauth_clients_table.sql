-- +goose Up
CREATE TABLE oauth_clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL,
    client_name TEXT NOT NULL,
    client_type TEXT NOT NULL CHECK (client_type IN ('confidential', 'public')),
    client_secret TEXT,
    client_secret_expires_at TIMESTAMPTZ,
    redirect_uris TEXT[] NOT NULL DEFAULT '{}',
    create_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    status TEXT NOT NULL
);

CREATE UNIQUE INDEX oauth_clients_client_id_key ON oauth_clients (client_id);

-- +goose Down
DROP TABLE IF EXISTS oauth_clients;

-- name: CreateOauthClient :one
INSERT INTO oauth_clients (
    client_id,
    client_name,
    client_type,
    client_secret,
    client_secret_expires_at,
    redirect_uris,
    status
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING
    id,
    client_id,
    client_name,
    client_type,
    client_secret,
    client_secret_expires_at,
    redirect_uris,
    create_time,
    status;

-- name: OauthClientExistsByClientID :one
SELECT EXISTS (
    SELECT 1
    FROM oauth_clients
    WHERE client_id = $1
) AS client_exists;

-- name: GetOauthClientRedirectUrisByClientID :one
SELECT redirect_uris
FROM oauth_clients
WHERE client_id = $1;

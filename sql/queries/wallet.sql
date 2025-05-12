
-- name: CreateWallet :one

INSERT INTO wallets (id, user_id, address, encrypted_private_key, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

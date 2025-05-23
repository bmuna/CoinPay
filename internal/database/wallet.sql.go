// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: wallet.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createWallet = `-- name: CreateWallet :one
INSERT INTO wallets (id, user_id, address, encrypted_private_key, created_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, user_id, address, encrypted_private_key, created_at
`

type CreateWalletParams struct {
	ID                  uuid.UUID
	UserID              uuid.UUID
	Address             string
	EncryptedPrivateKey string
	CreatedAt           time.Time
}

func (q *Queries) CreateWallet(ctx context.Context, arg CreateWalletParams) (Wallet, error) {
	row := q.db.QueryRowContext(ctx, createWallet,
		arg.ID,
		arg.UserID,
		arg.Address,
		arg.EncryptedPrivateKey,
		arg.CreatedAt,
	)
	var i Wallet
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Address,
		&i.EncryptedPrivateKey,
		&i.CreatedAt,
	)
	return i, err
}

const getWallet = `-- name: GetWallet :one
SELECT id, user_id, address, encrypted_private_key, created_at FROM wallets WHERE user_id = $1
`

func (q *Queries) GetWallet(ctx context.Context, userID uuid.UUID) (Wallet, error) {
	row := q.db.QueryRowContext(ctx, getWallet, userID)
	var i Wallet
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Address,
		&i.EncryptedPrivateKey,
		&i.CreatedAt,
	)
	return i, err
}

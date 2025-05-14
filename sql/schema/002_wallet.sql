-- +goose Up

CREATE TABLE wallets (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id),
  address TEXT NOT NULL,
  encrypted_private_key TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL,
);

-- +goose Down
DROP TABLE IF EXISTS wallets;

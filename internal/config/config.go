package config

import "github.com/bmuna/CoinPay/backend/internal/database"

type ApiConfig struct {
	DB *database.Queries
}
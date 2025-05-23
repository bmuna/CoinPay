package models

import (
	"time"

	"github.com/bmuna/CoinPay/backend/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Token     *string    `json:"token,omitempty"`
}

func DatabaseUserToUser(dbUser database.User, token *string) User {
	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
		Password:  dbUser.Password,
		Token:     token,
	}
}

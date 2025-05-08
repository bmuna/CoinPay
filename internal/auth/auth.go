package internal

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

const (
	key    = "randomString"
	MaxAge = 84600 * 30
	isProd = false
)

func NewAuth() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error landing .env file")
	}

	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, "http://localhost:3000/auth/go"),
	)
}

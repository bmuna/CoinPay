package internal

import (
	"os"

	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
)

const (
	key    = "randomString"
	MaxAge = 84600 * 30
	isProd = false
)

func NewAuth() {

	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, "http://localhost:8081/auth/google/callback"),
	)
}

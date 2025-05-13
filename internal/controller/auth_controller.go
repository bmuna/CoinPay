package controller

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bmuna/CoinPay/backend/internal/config"
	"github.com/bmuna/CoinPay/backend/internal/database"
	"github.com/bmuna/CoinPay/backend/internal/models"
	"github.com/bmuna/CoinPay/backend/internal/security"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/markbates/goth/gothic"
	"golang.org/x/crypto/bcrypt"
)

type UserCredential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, struct{}{})
}

func Signin(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userCredential UserCredential
		err := json.NewDecoder(r.Body).Decode(&userCredential)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
			return
		}

		user, err := apiCfg.DB.GetUser(r.Context(), userCredential.Email)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				respondWithError(w, http.StatusUnauthorized, "User does not exist with this email")
				return
			}
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error fetching user: %v", err))
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userCredential.Password))
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}

		tokenString, _ := createToken(userCredential.Email)

		respondWithJSON(w, http.StatusOK, models.DatabaseUserToUser(user, &tokenString))

	}
}

func Signup(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var userCredential UserCredential

		err := json.NewDecoder(r.Body).Decode(&userCredential)
		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
			return
		}

		_, err = apiCfg.DB.GetUser(r.Context(), userCredential.Email)
		if err == nil {
			respondWithError(w, http.StatusConflict, "User already exists with this email")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userCredential.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Error when hashing the password")
		}

		user, err := apiCfg.DB.CreateUser(
			r.Context(),
			database.CreateUserParams{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Email:     userCredential.Email,
				Password:  string(hashedPassword),
			},
		)

		if err != nil {
			respondWithError(w, 400, fmt.Sprintf("coudn't create user: %v", err))
		}

		respondWithJSON(w, 200, models.DatabaseUserToUser(user, nil))
	}
}

func BiginAuthProviderCallback(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))
	gothic.BeginAuthHandler(w, r)
}

func GetAuthCallBackFuction(apiCfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provider := chi.URLParam(r, "provider")
		r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

		user, err := gothic.CompleteUserAuth(w, r)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		email := user.Email

		_, err = apiCfg.DB.GetUser(r.Context(), email)

		if err != nil {
			randomPassword := uuid.New().String()
			hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(randomPassword), bcrypt.DefaultCost)
			if hashErr != nil {
				http.Error(w, "error hashing password", http.StatusInternalServerError)
				return
			}

			userId := uuid.New()

			_, err = apiCfg.DB.CreateUser(
				r.Context(),
				database.CreateUserParams{
					ID:        userId,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Email:     email,
					Password:  string(hashedPassword),
				},
			)

			if err != nil {
				http.Error(w, fmt.Sprintf("could not create user: %v", err), http.StatusInternalServerError)
				return
			}

			privateKeyHex, address := CreateWallet()

			encryptedKey, err := security.Encrypt(privateKeyHex, os.Getenv("ENCRYPTION_SECRET"))
			if err != nil {
				log.Printf("Error when encrypting private key %v", err)
			}

			_, err = apiCfg.DB.CreateWallet(r.Context(), database.CreateWalletParams{
				ID:                  uuid.New(),
				UserID:              userId,
				Address:             address,
				EncryptedPrivateKey: encryptedKey,
				CreatedAt:           time.Now(),
			})

			if err != nil {
				log.Printf("Error when creating a wallet %v", err)
			}

		}

		tokenString, _ := createToken(email)
		//  (w, 200, models.DatabaseUserToUser(dbUser, &tokenString))

		redirectURL := fmt.Sprintf("myapp://auth_callback?email=%s&token=%s", email, tokenString)
		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)

	}
}



func handlerSendETH() {

}

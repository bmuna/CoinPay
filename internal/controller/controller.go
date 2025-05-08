package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bmuna/CoinPay/backend/internal/config"
	"github.com/bmuna/CoinPay/backend/internal/database"
	"github.com/bmuna/CoinPay/backend/internal/models"
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

func Login(w http.ResponseWriter, r *http.Request) {

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

		respondWithJSON(w, 200, models.DatabaseUserToUser(user))
	}
}

func GetAuthCallBackFuction(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	fmt.Println(user)
}

func handlerSendETH() {

}

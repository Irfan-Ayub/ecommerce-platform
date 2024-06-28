package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Irfan-Ayub/ecommerce-platform/user-service/models"
	"github.com/Irfan-Ayub/ecommerce-platform/user-service/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Signup(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Signup Request Recieved")
		var user models.User
		_ = json.NewDecoder(r.Body).Decode(&user)

		// log.Println("User Recieved", user)

		// Hash the Password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error Hasing Password", http.StatusInternalServerError)
		}

		user.Password = string(hashedPassword)

		db.Create(&user)
		json.NewEncoder(w).Encode(user)
	}
}

func Login(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Login Request Recieved")
		var user models.User
		_ = json.NewDecoder(r.Body).Decode(&user)

		var dbUser models.User
		db.Where("username = ?", user.Username).First(&dbUser)
		if dbUser.Username == "" {
			http.Error(w, "Invalid Username", http.StatusUnauthorized)
			return
		}

		// Compare the hased password with the plain text password
		err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
		if err != nil {
			http.Error(w, "Invalid Password", http.StatusUnauthorized)
			return
		}

		// Generate JWT token
		token, err := utils.GenerateJWT(user.Username)
		if err != nil {
			http.Error(w, "Error Generating Token", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "login-token",
			Value:   token,
			Expires: time.Now().Add(5 * time.Minute),
		})
	}
}

func Profile(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Profile Request Recieved")

		cookie, err := r.Cookie("login-token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "UnAuthorized", http.StatusUnauthorized)
				return
			}

			http.Error(w, "UnAuthorized", http.StatusBadRequest)
			return
		}

		tokenStr := cookie.Value
		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			http.Error(w, "UnAuthorized", http.StatusUnauthorized)
			return
		}

		var user models.User

		db.Where("username = ?", claims.Username).First(&user)
		json.NewEncoder(w).Encode(user)
	}
}

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Irfan-Ayub/ecommerce-platform/user-service/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load .env file from
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
	}

	// Database connection
	db_url := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(db_url), &gorm.Config{})
	if err != nil {
		log.Fatalf("Faled to connect to daabase: %v", err)
	}

	// Auto migrate User Model
	// db.AutoMigrate(&models.User{})
	router := mux.NewRouter()
	router.HandleFunc("/signup", handlers.Signup(db)).Methods("POST")
	router.HandleFunc("/login", handlers.Login(db)).Methods("POST")
	router.HandleFunc("/profile", handlers.Profile(db)).Methods("GET")
	log.Println("Server Started on : 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Irfan-Ayub/ecommerce-platform/product-service/handlers"
	"github.com/Irfan-Ayub/ecommerce-platform/product-service/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	// Load Environment Variabls from .env file
	if err != nil {
		log.Fatalf("Error Loading .env file: %v", err)
	}

	// Database Connection
	db_url := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(db_url), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	// Auto migrate Product model
	db.AutoMigrate(&models.Product{})

	// Initialize router
	r := mux.NewRouter()
	r.HandleFunc("/products", handlers.CreateProduct(db)).Methods("POST")
	r.HandleFunc("/products", handlers.GetProducts(db)).Methods("GET")
	r.HandleFunc("/products/{id}", handlers.GetProduct(db)).Methods("GET")

	log.Println("Server Started on: 8001")
	log.Fatal(http.ListenAndServe(":8001", r))
}

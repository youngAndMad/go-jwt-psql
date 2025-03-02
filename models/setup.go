package models

import (
	"database/sql"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
	"log"
	"os"
)

var db *gorm.DB

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load env file: %v", err)
	}

	postgresDSN := os.Getenv("POSTGRES_DSN")

	dbConn, err := sql.Open("postgres", postgresDSN)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	} else {
		log.Println("Connected to database")
		defer dbConn.Close()
	}
	db, err = gorm.Open(postgres.New(postgres.Config{Conn: dbConn}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
		return
	}
}

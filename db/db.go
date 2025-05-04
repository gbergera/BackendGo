package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {
	err := godotenv.Load("db/.env")
	if err != nil {
		log.Println("⚠️ Warning: Could not load .env file")
	}

	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	if dbHost == "" {
		dbHost = "db"
	}
	if dbPort == "" {
		dbPort = "3306"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/ualabackend?parseTime=true", dbUser, dbPassword, dbHost, dbPort)

	maxAttempts := 10
	var db *sql.DB

	for i := 1; i <= maxAttempts; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
		}

		if err == nil {
			fmt.Println("✅ Successfully connected to the database.")
			DB = db
			return DB, nil
		}

		log.Printf("⏳ Attempt %d/%d - Could not connect to DB: %v", i, maxAttempts, err)
		time.Sleep(time.Duration(i) * time.Second)
	}

	return nil, fmt.Errorf("❌ failed to connect to DB after %d attempts: %v", maxAttempts, err)
}

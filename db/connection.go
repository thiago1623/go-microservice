package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var dbUser = os.Getenv("POSTGRES_USER")
var dbPassword = os.Getenv("POSTGRES_PASSWORD")
var dbName = os.Getenv("POSTGRES_DB")
var dbHost = os.Getenv("POSTGRES_HOST")
var dbPort = os.Getenv("POSTGRES_PORT")
var dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)

func ConnectDB() {

	maxRetries := 10                 // Número máximo de intentos de conexión
	retryInterval := 5 * time.Second // Intervalo entre reintentos

	for i := 0; i < maxRetries; i++ {
		var err error
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("Attempt %d failed to connect to DB: %v\n", i+1, err)
			time.Sleep(retryInterval)
		} else {
			log.Println("DB connected..")
			return
		}
	}

	log.Fatalf("Failed to connect to DB after %d retries.", maxRetries)
}

func GetDB() *gorm.DB {
	return DB
}
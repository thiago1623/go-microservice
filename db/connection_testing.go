package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DBTesting *gorm.DB
var dbUserTesting = os.Getenv("POSTGRES_USER_TESTING")
var dbPasswordTesting = os.Getenv("POSTGRES_PASSWORD_TESTING")
var dbNameTesting = os.Getenv("POSTGRES_DB_TESTING")
var dbHostTesting = os.Getenv("POSTGRES_HOST_TESTING")
var dbPortTesting = os.Getenv("POSTGRES_PORT_TESTING")
var dsnTesting = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHostTesting, dbUserTesting, dbPasswordTesting, dbNameTesting, dbPortTesting)

func ConnectDBTesting() {
	maxRetries := 10                 // Maximum number of connection attempts
	retryInterval := 5 * time.Second // Interval between reintents

	dsnTesting := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHostTesting, dbUserTesting, dbPasswordTesting, dbNameTesting, dbPortTesting)

	for i := 0; i < maxRetries; i++ {
		var err error
		DBTesting, err = gorm.Open(postgres.Open(dsnTesting), &gorm.Config{})
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

func GetDBTesting() *gorm.DB {
	return DBTesting
}

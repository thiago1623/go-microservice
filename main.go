package main

import (
	"net/http"

	"github.com/thiago1623/go-microservice/api/routes"
	"github.com/thiago1623/go-microservice/db"
	"github.com/thiago1623/go-microservice/models"
)

func main() {
	db.ConnectDB()

	db.DB.AutoMigrate(models.EnergyConsumption{})
	router := routes.SetupRoutes()
	http.ListenAndServe(":8001", router)
}

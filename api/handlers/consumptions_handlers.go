package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/thiago1623/go-microservice/db"
	"github.com/thiago1623/go-microservice/models"
)

func GetConsumptionHandler(w http.ResponseWriter, r *http.Request) {
	kindPeriod := r.URL.Query().Get("kind_period")

	switch kindPeriod {
	case "monthly":
		GetMonthlyConsumptionHandler(w, r)
	case "weekly":
		GetWeeklyConsumptionHandler(w, r)
	case "daily":
		GetDailyConsumptionHandler(w, r)
	default:
		http.Error(w, "Invalid kind_period value", http.StatusBadRequest)
	}
}

func GetConsumptionsHandler(w http.ResponseWriter, r *http.Request) {
	var energy_consumptions []models.EnergyConsumption
	db.GetDB().Limit(15).Find(&energy_consumptions)
	json.NewEncoder(w).Encode(&energy_consumptions)
}

package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/thiago1623/go-microservice/db"
	"github.com/thiago1623/go-microservice/models"
)

func GetDailyConsumptionHandler(w http.ResponseWriter, r *http.Request) {
	meterIDs := r.URL.Query().Get("meters_ids")
	startDate, _ := time.Parse("2006-01-02", r.URL.Query().Get("start_date"))
	endDate, _ := time.Parse("2006-01-02", r.URL.Query().Get("end_date"))
	testingMode := r.URL.Query().Get("testing_mode") == "true"

	energyConsumptions := fetchEnergyConsumptions(meterIDs, startDate, endDate, testingMode)

	response := make(map[string]interface{})
	response["period"] = getDailyPeriods(startDate, endDate)
	response["data_graph"] = aggregateDailyData(energyConsumptions, startDate, endDate)

	json.NewEncoder(w).Encode(response)
}

func getDailyPeriods(startDate, endDate time.Time) []string {
	var periods []string

	currentDate := startDate
	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		period := currentDate.Format("Jan 2")
		periods = append(periods, period)

		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return periods
}

func aggregateDailyData(consumptions []models.EnergyConsumption, startDate, endDate time.Time) []map[string]interface{} {
	var dataGraph []map[string]interface{}

	// Group the energyConsumptions by meter ID
	consumptionsByMeterID := make(map[int32][]models.EnergyConsumption)
	for _, consumption := range consumptions {
		consumptionsByMeterID[consumption.MeterID] = append(consumptionsByMeterID[consumption.MeterID], consumption)
	}

	// Aggregate data for each meter based on daily consumption
	for meterID, consumptions := range consumptionsByMeterID {
		meterData := make(map[string]interface{})
		meterData["meter_id"] = meterID

		// Fetch the address for this meter (assuming one address per meter)
		var address models.Address
		db.GetDB().Where("meter_id = ? AND start_date <= ? AND (end_date >= ? OR end_date IS NULL)", meterID, startDate, endDate).First(&address)

		if address.ID != uuid.Nil {
			meterData["address"] = address.Address
		} else {
			meterData["address"] = "Address not found"
		}

		active, reactiveInductive, reactiveCapacitive, exported := aggregateDailyDataForMeter(consumptions, startDate, endDate)
		meterData["active"] = active
		meterData["reactive_inductive"] = reactiveInductive
		meterData["reactive_capacitive"] = reactiveCapacitive
		meterData["exported"] = exported

		dataGraph = append(dataGraph, meterData)
	}

	return dataGraph
}

func aggregateDailyDataForMeter(consumptions []models.EnergyConsumption, startDate, endDate time.Time) ([]float64, []float64, []float64, []float64) {
	active := make([]float64, 0)
	reactiveInductive := make([]float64, 0)
	reactiveCapacitive := make([]float64, 0)
	exported := make([]float64, 0)

	consumptionsByDay := make(map[int]models.EnergyConsumption)

	for _, consumption := range consumptions {
		consumptionsByDay[consumption.Date.Day()] = consumption
	}

	currentDate := startDate
	for currentDate.Before(endDate) || currentDate.Equal(endDate) {
		if consumption, found := consumptionsByDay[currentDate.Day()]; found {
			active = append(active, consumption.ActiveEnergy)
			reactiveInductive = append(reactiveInductive, consumption.ReactiveEnergy)
			reactiveCapacitive = append(reactiveCapacitive, consumption.CapacitiveReactive)
			exported = append(exported, consumption.Solar)
		} else {
			active = append(active, 0)
			reactiveInductive = append(reactiveInductive, 0)
			reactiveCapacitive = append(reactiveCapacitive, 0)
			exported = append(exported, 0)
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return active, reactiveInductive, reactiveCapacitive, exported
}

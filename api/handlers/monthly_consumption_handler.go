package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/thiago1623/go-microservice/db"
	"github.com/thiago1623/go-microservice/models"
)

func getPeriods(startDate, endDate time.Time, periodType string) []string {
	var periods []string

	// Determine the time interval between start and end date based on the period type
	interval := endDate.Sub(startDate)
	switch periodType {
	case "monthly":
		months := int(interval.Hours() / 24 / 30.44) // Assuming an average month has 30.44 days
		for i := 0; i <= months; i++ {
			period := startDate.AddDate(0, i, 0).Format("Jan 2006")
			periods = append(periods, period)
		}
	}
	return periods
}

func aggregateMonthlyData(consumptions []models.EnergyConsumption, startDate, endDate time.Time) []map[string]interface{} {
	var dataGraph []map[string]interface{}

	// Group the energyConsumptions by meter ID
	consumptionsByMeterID := make(map[int32][]models.EnergyConsumption)
	for _, consumption := range consumptions {
		consumptionsByMeterID[consumption.MeterID] = append(consumptionsByMeterID[consumption.MeterID], consumption)
	}

	// Aggregate data for each meter based on monthly consumption
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

		active, reactiveInductive, reactiveCapacitive, exported := aggregateMonthlyDataForMeter(consumptions)
		meterData["active"] = active
		meterData["reactive_inductive"] = reactiveInductive
		meterData["reactive_capacitive"] = reactiveCapacitive
		meterData["exported"] = exported

		dataGraph = append(dataGraph, meterData)
	}

	return dataGraph
}

func aggregateMonthlyDataForMeter(consumptions []models.EnergyConsumption) ([]float64, []float64, []float64, []float64) {
	var active []float64
	var reactiveInductive []float64
	var reactiveCapacitive []float64
	var exported []float64

	// Creamos un mapa para agrupar los datos de consumo por mes
	consumptionsByMonth := make(map[string]models.EnergyConsumption)

	// Iteramos sobre los datos de consumo para acumularlos por mes
	for _, consumption := range consumptions {
		monthYear := consumption.Date.Format("Jan 2006")
		if existingConsumption, found := consumptionsByMonth[monthYear]; found {
			// Ya existe un consumo para este mes, acumulamos los valores
			existingConsumption.ActiveEnergy += consumption.ActiveEnergy
			existingConsumption.ReactiveEnergy += consumption.ReactiveEnergy
			existingConsumption.CapacitiveReactive += consumption.CapacitiveReactive
			existingConsumption.Solar += consumption.Solar
			consumptionsByMonth[monthYear] = existingConsumption
		} else {
			// No hay consumo previo para este mes, añadimos uno nuevo
			consumptionsByMonth[monthYear] = consumption
		}
	}

	// Extraemos los valores acumulados para cada mes
	for _, consumption := range consumptionsByMonth {
		active = append(active, consumption.ActiveEnergy)
		reactiveInductive = append(reactiveInductive, consumption.ReactiveEnergy)
		reactiveCapacitive = append(reactiveCapacitive, consumption.CapacitiveReactive)
		exported = append(exported, consumption.Solar)
	}

	return active, reactiveInductive, reactiveCapacitive, exported
}

func GetMonthlyConsumptionHandler(w http.ResponseWriter, r *http.Request) {
	meterIDs := r.URL.Query().Get("meters_ids")
	startDate, _ := time.Parse("2006-01-02", r.URL.Query().Get("start_date"))
	endDate, _ := time.Parse("2006-01-02", r.URL.Query().Get("end_date"))

	energyConsumptions := fetchEnergyConsumptions(meterIDs, startDate, endDate)

	response := make(map[string]interface{})
	response["period"] = getPeriods(startDate, endDate, "monthly")
	response["data_graph"] = aggregateMonthlyData(energyConsumptions, startDate, endDate)

	json.NewEncoder(w).Encode(response)
}

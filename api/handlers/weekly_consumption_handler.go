package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/thiago1623/go-microservice/db"
	"github.com/thiago1623/go-microservice/models"
)

func GetWeeklyConsumptionHandler(w http.ResponseWriter, r *http.Request) {
	meterIDs := r.URL.Query().Get("meters_ids")
	startDate, _ := time.Parse("2006-01-02", r.URL.Query().Get("start_date"))
	endDate, _ := time.Parse("2006-01-02", r.URL.Query().Get("end_date"))

	energyConsumptions := fetchEnergyConsumptions(meterIDs, startDate, endDate)

	response := make(map[string]interface{})
	response["period"] = getWeeklyPeriods(startDate, endDate)
	response["data_graph"] = aggregateWeeklyData(energyConsumptions, startDate, endDate)

	json.NewEncoder(w).Encode(response)
}

func aggregateWeeklyData(consumptions []models.EnergyConsumption, startDate, endDate time.Time) []map[string]interface{} {
	var dataGraph []map[string]interface{}

	// Group the energyConsumptions by meter ID
	consumptionsByMeterID := make(map[int32][]models.EnergyConsumption)
	for _, consumption := range consumptions {
		consumptionsByMeterID[consumption.MeterID] = append(consumptionsByMeterID[consumption.MeterID], consumption)
	}

	// Aggregate data for each meter based on weekly consumption
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

		active, reactiveInductive, reactiveCapacitive, exported := aggregateWeeklyDataForMeter(consumptions, startDate, endDate)
		meterData["active"] = active
		meterData["reactive_inductive"] = reactiveInductive
		meterData["reactive_capacitive"] = reactiveCapacitive
		meterData["exported"] = exported

		dataGraph = append(dataGraph, meterData)
	}

	return dataGraph
}

func aggregateWeeklyDataForMeter(consumptions []models.EnergyConsumption, startDate, endDate time.Time) ([]string, []string, []string, []string) {
	var active []string
	var reactiveInductive []string
	var reactiveCapacitive []string
	var exported []string

	// Creamos un mapa para agrupar los datos de consumo por semana
	consumptionsByWeek := make(map[int]string)

	// Iteramos sobre los datos de consumo para agruparlos por semana
	for _, consumption := range consumptions {
		// Calculamos el número de semana para el consumo actual
		weekNumber := int(consumption.Date.Weekday())

		// Utilizamos el número de semana como clave para agrupar los datos de consumo
		if week, found := consumptionsByWeek[weekNumber]; found {
			// Ya existe un consumo para esta semana, acumulamos los valores
			consumptionsByWeek[weekNumber] = week + fmt.Sprintf(",%.3f", consumption.ActiveEnergy) +
				fmt.Sprintf(",%.3f", consumption.ReactiveEnergy) +
				fmt.Sprintf(",%.3f", consumption.CapacitiveReactive) +
				fmt.Sprintf(",%.3f", consumption.Solar)
		} else {
			// No hay consumo previo para esta semana, añadimos uno nuevo
			consumptionsByWeek[weekNumber] = fmt.Sprintf("%.3f", consumption.ActiveEnergy) +
				fmt.Sprintf(",%.3f", consumption.ReactiveEnergy) +
				fmt.Sprintf(",%.3f", consumption.CapacitiveReactive) +
				fmt.Sprintf(",%.3f", consumption.Solar)
		}
	}

	// Ordenamos las semanas para asegurarnos de que están en el orden correcto
	var weeks []int
	for week := range consumptionsByWeek {
		weeks = append(weeks, week)
	}
	sort.Ints(weeks)

	// Extraemos los valores acumulados para cada semana
	for _, week := range weeks {
		data := strings.Split(consumptionsByWeek[week], ",")
		active = append(active, data[0])
		reactiveInductive = append(reactiveInductive, data[1])
		reactiveCapacitive = append(reactiveCapacitive, data[2])
		exported = append(exported, data[3])
	}

	return active, reactiveInductive, reactiveCapacitive, exported
}

func getWeeklyPeriods(startDate, endDate time.Time) []string {
	var periods []string

	// Obtenemos el primer día de la primera semana
	currentDate := startDate
	for currentDate.Weekday() != time.Sunday {
		currentDate = currentDate.AddDate(0, 0, -1)
	}

	// Iteramos por las semanas y construimos los períodos
	for currentDate.Before(endDate) {
		// Obtenemos el último día de la semana
		endOfWeek := currentDate.AddDate(0, 0, 6)
		period := fmt.Sprintf("%s - %s", currentDate.Format("Jan 2"), endOfWeek.Format("Jan 2"))
		periods = append(periods, period)

		// Pasamos al próximo domingo para la siguiente semana
		currentDate = endOfWeek.AddDate(0, 0, 1)
	}

	return periods
}

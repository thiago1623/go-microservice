package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/thiago1623/go-microservice/api/handlers"
	"github.com/thiago1623/go-microservice/db"
	"github.com/thiago1623/go-microservice/models"
)

func TestGetMonthlyConsumptionHandler(t *testing.T) {
	// Configuración de la base de datos de testing
	db.ConnectDBTesting()
	defer db.GetDBTesting().Where("meter_id IN (?)", []int32{1, 2}).Delete(&models.EnergyConsumption{})

	// Insertar datos de prueba en la base de datos de testing
	consumption1 := models.EnergyConsumption{
		MeterID:            1,
		ActiveEnergy:       4401580.218640001,
		ReactiveEnergy:     475529.59263999935,
		CapacitiveReactive: 10.988000000000001,
		Solar:              5943.037877272399,
		Date:               time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
	}
	db.GetDBTesting().Create(&consumption1)

	consumption2 := models.EnergyConsumption{
		MeterID:            2,
		ActiveEnergy:       4008410.949719998,
		ReactiveEnergy:     2443377.975379999,
		CapacitiveReactive: 0,
		Solar:              153.72857775836155,
		Date:               time.Date(2023, 6, 30, 0, 0, 0, 0, time.UTC),
	}
	db.GetDBTesting().Create(&consumption2)

	consumption3 := models.EnergyConsumption{
		MeterID:            1,
		ActiveEnergy:       1563573.5172400007,
		ReactiveEnergy:     166221.74783000012,
		CapacitiveReactive: 3.365999999999999,
		Solar:              2017.5683214119479,
		Date:               time.Date(2023, 7, 5, 0, 0, 0, 0, time.UTC),
	}
	db.GetDBTesting().Create(&consumption3)

	consumption4 := models.EnergyConsumption{
		MeterID:            2,
		ActiveEnergy:       10898488.745770002,
		ReactiveEnergy:     6605406.771730003,
		CapacitiveReactive: 0,
		Solar:              468.7381462616537,
		Date:               time.Date(2023, 7, 10, 0, 0, 0, 0, time.UTC),
	}
	db.GetDBTesting().Create(&consumption4)

	// Simular petición HTTP a través de httptest
	req, err := http.NewRequest("GET", "/consumption?meters_ids=1,2&start_date=2023-06-01&end_date=2023-07-10&kind_period=monthly&testing_mode=true", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.GetMonthlyConsumptionHandler)
	handler.ServeHTTP(rr, req)

	// Verificar el código de estado HTTP
	assert.Equal(t, http.StatusOK, rr.Code, "El código de estado esperado es 200")

	// Verificar la estructura de la respuesta JSON
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Error al decodificar la respuesta JSON: %v", err)
	}

	// Verificar que la respuesta tiene la estructura esperada
	expectedPeriods := []string{"Jun 2023", "Jul 2023"}
	periods, ok := response["period"].([]interface{})
	assert.True(t, ok, "La clave 'period' no tiene el tipo de dato esperado.")
	for i, p := range periods {
		assert.Equal(t, expectedPeriods[i], p, "Periodo esperado y periodo obtenido no coinciden")
	}

	dataGraph, ok := response["data_graph"].([]interface{})
	assert.True(t, ok, "La clave 'data_graph' no tiene el tipo de dato esperado.")

	// Verificar la estructura y los valores de cada elemento en 'data_graph'
	for _, data := range dataGraph {
		meterData, ok := data.(map[string]interface{})
		assert.True(t, ok, "El elemento en 'data_graph' no tiene el tipo de dato esperado.")

		_, ok = meterData["meter_id"].(float64)
		assert.True(t, ok, "La clave 'meter_id' en 'data_graph' no tiene el tipo de dato esperado.")

		_, ok = meterData["address"].(string)
		assert.True(t, ok, "La clave 'address' en 'data_graph' no tiene el tipo de dato esperado.")

		_, ok = meterData["active"].([]interface{})
		assert.True(t, ok, "La clave 'active' en 'data_graph' no tiene el tipo de dato esperado.")

		_, ok = meterData["reactive_inductive"].([]interface{})
		assert.True(t, ok, "La clave 'reactive_inductive' en 'data_graph' no tiene el tipo de dato esperado.")

		_, ok = meterData["reactive_capacitive"].([]interface{})
		assert.True(t, ok, "La clave 'reactive_capacitive' en 'data_graph' no tiene el tipo de dato esperado.")
		_, ok = meterData["exported"].([]interface{})
		assert.True(t, ok, "La clave 'exported' en 'data_graph' no tiene el tipo de dato esperado.")
	}

	// Limpiar los datos de prueba de la base de datos de testing después de las pruebas
	db.GetDBTesting().Where("meter_id IN (?)", []int32{1, 2}).Delete(&models.EnergyConsumption{})
}

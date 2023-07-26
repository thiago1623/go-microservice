// importer/import_date.go

package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/thiago1623/go-microservice/db"
	"github.com/thiago1623/go-microservice/models"
)

// Función para leer el archivo CSV y guardar los datos en la base de datos
func ImportCSVToDB() {
	file, err := os.Open("db_complete.csv")
	if err != nil {
		fmt.Println("Error al abrir el archivo:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error al leer el archivo CSV:", err)
		return
	}

	// Conexión a la base de datos
	db.ConnectDB()
	dbInstance := db.GetDB()

	for _, record := range records {
		// Parsear los datos del registro CSV
		meterID, _ := strconv.Atoi(record[1]) // Convertir a int32
		activeEnergy, _ := strconv.ParseFloat(record[2], 64)
		reactiveEnergy, _ := strconv.ParseFloat(record[3], 64)
		capacitiveReactive, _ := strconv.ParseFloat(record[4], 64)
		solar, _ := strconv.ParseFloat(record[5], 64)
		date, _ := time.Parse("2006-01-02 15:04:05-07:00", record[6]) // Convertir a time.Time

		// Crear una nueva instancia de EnergyConsumption
		energyConsumption := models.EnergyConsumption{
			MeterID:            int32(meterID), // MeterID es de tipo int32
			ActiveEnergy:       activeEnergy,
			ReactiveEnergy:     reactiveEnergy,
			CapacitiveReactive: capacitiveReactive,
			Solar:              solar,
			Date:               date, // Date es de tipo time.Time
		}

		// Guardar el registro en la base de datos
		dbInstance.Create(&energyConsumption)
	}

	fmt.Println("Datos importados desde el archivo CSV.")
}

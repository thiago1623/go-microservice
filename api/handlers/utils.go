package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/thiago1623/go-microservice/db"
	"github.com/thiago1623/go-microservice/models"
)

func fetchEnergyConsumptions(meterIDsStr string, startDate, endDate time.Time) []models.EnergyConsumption {
	// Dividir el string de meterIDs en IDs individuales
	meterIDsStrList := strings.Split(meterIDsStr, ",")
	var meterIDsInt []int32

	// Convertir los elementos de meterIDs de strings a enteros (int32)
	for _, idStr := range meterIDsStrList {
		idInt, err := strconv.Atoi(idStr)
		if err == nil {
			meterIDsInt = append(meterIDsInt, int32(idInt))
		}
	}

	var energyConsumptions []models.EnergyConsumption
	if len(meterIDsInt) > 0 {
		db.GetDB().Preload("Address").Where("meter_id IN (?) AND date >= ? AND date <= ?", meterIDsInt, startDate, endDate).Find(&energyConsumptions)
	}

	return energyConsumptions
}

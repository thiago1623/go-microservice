package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/thiago1623/go-microservice/db"
	"github.com/thiago1623/go-microservice/models"
	"gorm.io/gorm"
)

func fetchEnergyConsumptions(meterIDsStr string, startDate, endDate time.Time, testingMode bool) []models.EnergyConsumption {
	meterIDsStrList := strings.Split(meterIDsStr, ",")
	var meterIDsInt []int32

	for _, idStr := range meterIDsStrList {
		idInt, err := strconv.Atoi(idStr)
		if err == nil {
			meterIDsInt = append(meterIDsInt, int32(idInt))
		}
	}

	var energyConsumptions []models.EnergyConsumption
	var dbInstance *gorm.DB

	if testingMode {
		dbInstance = db.GetDBTesting()
	} else {
		dbInstance = db.GetDB()
	}

	if len(meterIDsInt) > 0 {
		dbInstance.Preload("Address").Where("meter_id IN (?) AND date >= ? AND date <= ?", meterIDsInt, startDate, endDate).Find(&energyConsumptions)
	}

	return energyConsumptions
}

package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type EnergyConsumption struct {
	gorm.Model
	ID                 uuid.UUID `gorm:"type:uuid;primary_key;"`
	MeterID            int32
	ActiveEnergy       float64
	ReactiveEnergy     float64
	CapacitiveReactive float64
	Solar              float64
	Date               time.Time `gorm:"type:date"`
}

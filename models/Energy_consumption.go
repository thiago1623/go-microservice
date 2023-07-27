package models

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Address struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;primary_key;"`
	MeterID   int32
	Address   string
	StartDate *time.Time
	EndDate   *time.Time
}

type EnergyConsumption struct {
	gorm.Model
	ID                 uuid.UUID `gorm:"type:uuid;primary_key;"`
	MeterID            int32
	ActiveEnergy       float64
	ReactiveEnergy     float64
	CapacitiveReactive float64
	Solar              float64
	Date               time.Time `gorm:"type:date"`
	AddressID          uuid.UUID `gorm:"type:uuid;"`
	Address            Address   `gorm:"foreignKey:AddressID"`
}

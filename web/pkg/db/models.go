package db

import (
	"time"

	"gorm.io/gorm"
)

type Device struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
}

type Unit struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
}

type MetricType struct {
	gorm.Model
	Name string `gorm:"unique;not null"`
}

type Metric struct {
	gorm.Model
	TypeID     uint
	Value      string `gorm:"not null"` // Ensure Value is a string
	UnitID     uint
	DeviceID   uint
	Type       MetricType `gorm:"foreignKey:TypeID"`
	Unit       Unit       `gorm:"foreignKey:UnitID"`
	Device     Device     `gorm:"foreignKey:DeviceID"`
	RecordedAt time.Time  `gorm:"not null"`
}

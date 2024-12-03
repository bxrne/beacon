package db

import (
	"time"

	"gorm.io/gorm"
)

type Device struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type Unit struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type MetricType struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type Metric struct {
	gorm.Model
	TypeID     uint
	Value      float64
	UnitID     uint
	DeviceID   uint
	Type       MetricType `gorm:"foreignKey:TypeID"`
	Unit       Unit       `gorm:"foreignKey:UnitID"`
	Device     Device     `gorm:"foreignKey:DeviceID"`
	RecordedAt time.Time
}

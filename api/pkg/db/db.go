package db

import (
	"fmt"

	"github.com/bxrne/beacon/api/pkg/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := migrate(db, cfg); err != nil {
		return nil, err
	}

	return db, nil
}

func migrate(db *gorm.DB, cfg *config.Config) error {
	if err := db.AutoMigrate(&Device{}, &Unit{}, &MetricType{}, &Metric{}); err != nil {
		return err
	}

	units := []Unit{{Name: "percent"}, {Name: "bytes"}, {Name: "seconds"}}
	for _, unit := range units {
		db.FirstOrCreate(&unit, Unit{Name: unit.Name})
	}

	metricTypes := []MetricType{{Name: "cpu_usage"}, {Name: "disk_usage"}, {Name: "uptime"}}
	for _, metricType := range metricTypes {
		db.FirstOrCreate(&metricType, MetricType{Name: metricType.Name})
	}

	return nil
}

func RegisterDevice(db *gorm.DB, name string) error {
	if name == "" {
		return fmt.Errorf("name is empty")
	}

	device := Device{Name: name}
	return db.FirstOrCreate(&device, Device{Name: name}).Error
}

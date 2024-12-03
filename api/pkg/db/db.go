package db

import (
	"fmt"

	"github.com/bxrne/beacon-web/pkg/metrics"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabase(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}

func migrate(db *gorm.DB) error {
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

func PersistMetric(db *gorm.DB, deviceMetrics metrics.DeviceMetrics, name string) error {
	if name == "" {
		return fmt.Errorf("name is empty")
	}

	var device Device
	if err := db.First(&device, "name = ?", name).Error; err != nil {
		return fmt.Errorf("device not found")
	}

	for _, metric := range deviceMetrics.Metrics {
		var metricType MetricType
		if err := db.First(&metricType, "name = ?", metric.Type).Error; err != nil {
			return err
		}

		var unit Unit
		if err := db.First(&unit, "name = ?", metric.Unit).Error; err != nil {
			return err
		}

		newMetric := Metric{
			TypeID:   metricType.ID,
			Value:    metric.Value,
			UnitID:   unit.ID,
			DeviceID: device.ID,
		}
		if err := db.Create(&newMetric).Error; err != nil {
			return err
		}
	}

	return nil
}

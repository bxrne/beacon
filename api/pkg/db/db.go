package db

import (
	"fmt"
	"time"

	"github.com/bxrne/beacon/api/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	db_cfg := &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	db, err := gorm.Open(sqlite.Open(cfg.Database.DSN), db_cfg)

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

	// Check if indexes exist before creating them
	if !indexExists(db, "idx_metrics_device_id") {
		if err := db.Exec("CREATE INDEX idx_metrics_device_id ON metrics(device_id)").Error; err != nil {
			return err
		}
	}
	if !indexExists(db, "idx_metrics_recorded_at") {
		if err := db.Exec("CREATE INDEX idx_metrics_recorded_at ON metrics(recorded_at)").Error; err != nil {
			return err
		}
	}

	var units []Unit
	for _, unit := range cfg.Metrics.Units {
		units = append(units, Unit{Name: unit})
	}
	for _, unit := range units {
		db.FirstOrCreate(&unit, Unit{Name: unit.Name})
	}

	var types []MetricType
	for _, metricType := range cfg.Metrics.Types {
		types = append(types, MetricType{Name: metricType})
	}
	for _, metricType := range types {
		db.FirstOrCreate(&metricType, MetricType{Name: metricType.Name})
	}

	return nil
}

func indexExists(db *gorm.DB, indexName string) bool {
	var result int
	db.Raw("SELECT COUNT(1) FROM sqlite_master WHERE type='index' AND name=?", indexName).Scan(&result)
	return result > 0
}

func RegisterDevice(db *gorm.DB, name string) error {
	if name == "" {
		return fmt.Errorf("name is empty")
	}

	device := Device{Name: name}
	return db.FirstOrCreate(&device, Device{Name: name}).Error
}

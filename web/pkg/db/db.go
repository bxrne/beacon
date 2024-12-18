package db

import (
	"fmt"
	"time"

	"github.com/bxrne/beacon/web/internal/config"
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
	if err := db.AutoMigrate(&Device{}, &Unit{}, &MetricType{}, &Metric{}, &CommandType{}, &Command{}); err != nil {
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

	// Allowed commands via lkup table
	var commands []CommandType
	for _, command := range cfg.Metrics.Commands {
		commands = append(commands, CommandType{Name: command})
	}
	for _, command := range commands {
		db.FirstOrCreate(&command, CommandType{Name: command.Name})
	}

	// Populate CommandType lookup table from config
	var commandTypes []CommandType
	for _, commandType := range cfg.CommandTypes {
		commandTypes = append(commandTypes, CommandType{Name: commandType.Name})
	}
	for _, commandType := range commandTypes {
		db.FirstOrCreate(&commandType, CommandType{Name: commandType.Name})
	}

	// Migrate existing CommandTypes from config
	var commandTypesFromConfig []CommandType
	for _, cmd := range cfg.Metrics.Commands {
		commandTypesFromConfig = append(commandTypesFromConfig, CommandType{Name: cmd})
	}
	for _, cmdType := range commandTypesFromConfig {
		if err := db.FirstOrCreate(&cmdType, CommandType{Name: cmdType.Name}).Error; err != nil {
			return fmt.Errorf("failed to create command type: %w", err)
		}
	}

	// Update existing metrics for 'car_light' and 'ped_light' to have unit 'color'

	var colorUnit Unit
	if err := db.First(&colorUnit, "name = ?", "color").Error; err != nil {
		return fmt.Errorf("failed to find 'color' unit: %w", err)
	}

	// Fetch MetricType IDs for 'car_light' and 'ped_light'
	var metricTypes []MetricType
	if err := db.Where("name IN ?", []string{"car_light", "ped_light"}).Find(&metricTypes).Error; err != nil {
		return fmt.Errorf("failed to fetch metric types: %w", err)
	}

	var metricTypeIDs []uint
	for _, mt := range metricTypes {
		metricTypeIDs = append(metricTypeIDs, mt.ID)
	}

	if len(metricTypeIDs) == 0 {
		return fmt.Errorf("no metric types found for 'car_light' and 'ped_light'")
	}

	// Update metrics where type_id is in the fetched IDs
	if err := db.Model(&Metric{}).Where("type_id IN ?", metricTypeIDs).Update("unit_id", colorUnit.ID).Error; err != nil {
		return fmt.Errorf("failed to update metrics unit: %w", err)
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

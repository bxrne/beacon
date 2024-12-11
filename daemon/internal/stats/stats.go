package stats

import (
	"fmt"

	"github.com/bxrne/beacon/web/pkg/db"
	"gorm.io/gorm"
)

// ValidateMetricType checks if the metric type exists in DB
func ValidateMetricType(gorm_db *gorm.DB, metricType string) error {
	var exists bool
	err := gorm_db.Model(&db.MetricType{}).Select("count(*) > 0").Where("name = ?", metricType).Find(&exists).Error
	if err != nil {
		return fmt.Errorf("failed to check metric type: %w", err)
	}
	if !exists {
		return fmt.Errorf("invalid metric type: %s", metricType)
	}
	return nil
}

// ValidateUnit checks if the unit exists in DB
func ValidateUnit(gorm_db *gorm.DB, unit string) error {
	var exists bool
	err := gorm_db.Model(&db.Unit{}).Select("count(*) > 0").Where("name = ?", unit).Find(&exists).Error
	if err != nil {
		return fmt.Errorf("failed to check unit: %s", unit)
	}
	if !exists {
		return fmt.Errorf("invalid unit: %s", unit)
	}
	return nil
}

// Validate checks if metric values exist in DB
func (m Metric) Validate(db *gorm.DB) error {
	if err := ValidateMetricType(db, m.Type); err != nil {
		return err
	}
	return ValidateUnit(db, m.Unit)
}

// Validate checks all metrics against DB values
func (dm DeviceMetrics) Validate(db *gorm.DB) error {
	for _, m := range dm.Metrics {
		if err := m.Validate(db); err != nil {
			return err
		}
	}
	return nil
}

// PersistMetric persists the metrics for a device
func PersistMetric(gorm_db *gorm.DB, deviceMetrics DeviceMetrics, name string) error {
	if name == "" {
		return fmt.Errorf("name is empty")
	}

	var device db.Device
	if err := gorm_db.First(&device, "name = ?", name).Error; err != nil {
		return fmt.Errorf("device not found")
	}

	for _, metric := range deviceMetrics.Metrics {
		var metricType db.MetricType
		if err := gorm_db.First(&metricType, "name = ?", metric.Type).Error; err != nil {
			return err
		}

		var unit db.Unit
		if err := gorm_db.First(&unit, "name = ?", metric.Unit).Error; err != nil {
			return err
		}

		newMetric := Metric{
			Type:  fmt.Sprintf("%d", metricType.ID),
			Value: metric.Value, // No need to format as string
			Unit:  fmt.Sprintf("%d", unit.ID),
		}
		if err := gorm_db.Create(&newMetric).Error; err != nil {
			return err
		}
	}

	return nil
}

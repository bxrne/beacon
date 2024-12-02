// pkg/metrics/metrics.go
package metrics

import (
	"database/sql"
	"fmt"
)

// These should match values in metric_types table
type Metric struct {
	Type  string  `json:"type"` // References metric_types.name
	Value float64 `json:"value"`
	Unit  string  `json:"unit"` // References units.name
}

type DeviceMetrics struct {
	Metrics []Metric `json:"metrics"`
}

// ValidateMetricType checks if the metric type exists in DB
func ValidateMetricType(db *sql.DB, metricType string) error {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM metric_types WHERE name = $1)", metricType).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check metric type: %w", err)
	}
	if !exists {
		return fmt.Errorf("invalid metric type: %s", metricType)
	}
	return nil
}

// ValidateUnit checks if the unit exists in DB
func ValidateUnit(db *sql.DB, unit string) error {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM units WHERE name = $1)", unit).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check unit: %w", err)
	}
	if !exists {
		return fmt.Errorf("invalid unit: %s", unit)
	}
	return nil
}

// Validate checks if metric values exist in DB
func (m Metric) Validate(db *sql.DB) error {
	if err := ValidateMetricType(db, m.Type); err != nil {
		return err
	}
	return ValidateUnit(db, m.Unit)
}

// Validate checks all metrics against DB values
func (dm DeviceMetrics) Validate(db *sql.DB) error {
	for _, m := range dm.Metrics {
		if err := m.Validate(db); err != nil {
			return err
		}
	}
	return nil
}

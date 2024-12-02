package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/bxrne/beacon-web/pkg/metrics"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func NewDatabase() (*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, fmt.Errorf("DATABASE_URL is empty")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func RegisterDevice(db *sql.DB, hostname string) error {
	if hostname == "" {
		return fmt.Errorf("hostname is empty")
	}

	var deviceName string
	err := db.QueryRow("SELECT device_name FROM devices WHERE device_name = $1", hostname).Scan(&deviceName)
	if err == nil {
		return nil
	}

	_, err = db.Exec("INSERT INTO devices (device_name) VALUES ($1)", hostname)
	if err != nil {
		return err
	}

	return nil
}

func PersistMetric(db *sql.DB, deviceMetrics metrics.DeviceMetrics, hostname string) error {
	if hostname == "" {
		return fmt.Errorf("hostname is empty")
	}

	var deviceId int
	err := db.QueryRow("SELECT device_id FROM devices WHERE device_name = $1", hostname).Scan(&deviceId)
	if err != nil {
		return fmt.Errorf("device not found")
	}

	for _, metric := range deviceMetrics.Metrics {
		metricTypeID, err := getMetricTypeID(db, metric.Type)
		if err != nil {
			return err
		}
		unitID, err := getUnitID(db, metric.Unit)
		if err != nil {
			return err
		}
		_, err = db.Exec(
			"INSERT INTO metrics (device_id, metric_type_id, unit_id, value) VALUES ($1, $2, $3, $4)",
			deviceId, metricTypeID, unitID, metric.Value,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func getMetricTypeID(db *sql.DB, name string) (int, error) {
	var id int
	err := db.QueryRow("SELECT id FROM metric_types WHERE name = $1", name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("metric type not found: %s", name)
	}
	return id, nil
}

func getUnitID(db *sql.DB, name string) (int, error) {
	var id int
	err := db.QueryRow("SELECT id FROM units WHERE name = $1", name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("unit not found: %s", name)
	}
	return id, nil
}

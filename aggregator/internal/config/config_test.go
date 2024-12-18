package config_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bxrne/beacon/aggregator/internal/config"
)

// TEST: GIVEN a valid TOML configuration file
// WHEN the Load function is called
// THEN it should correctly parse all fields and return a matching Config struct
func TestLoad_ValidConfig(t *testing.T) {
	content := `
[telemetry]
server = "http://localhost:8080"
retry_interval = 10
timeout = 30

[targets]
hosts = ["host1", "host2"]
frequencies = [5, 10]
ports = ["8080", "9090"]

[labels]
environment = "production"
service = "myapp"

[logging]
level = "info"
`
	tmpFile := createTempFile(t, content)
	defer os.Remove(tmpFile)

	cfg, err := config.Load(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load valid config: %v", err)
	}

	expected := &config.Config{
		Telemetry: config.Telemetry{
			Server:        "http://localhost:8080",
			RetryInterval: 10,
			Timeout:       30,
		},
		Targets: config.Targets{
			Hosts:       []string{"host1", "host2"},
			Frequencies: []int{5, 10},
			Ports:       []string{"8080", "9090"},
		},
		Labels: config.Labels{
			Environment: "production",
			Service:     "myapp",
		},
		Logging: config.Logging{
			Level: "info",
		},
	}

	if !reflect.DeepEqual(cfg, expected) {
		t.Errorf("Config mismatch\nGot: %+v\nWant: %+v", cfg, expected)
	}
}

// TEST: GIVEN a non-existent file path
// WHEN the Load function is called
// THEN it should return an error
func TestLoad_InvalidPath(t *testing.T) {
	_, err := config.Load("nonexistent.toml")
	if err == nil {
		t.Error("Expected error when loading non-existent file, got nil")
	}
}

// TEST: GIVEN an empty TOML file
// WHEN the Load function is called
// THEN it should return an error
func TestLoad_EmptyConfig(t *testing.T) {
	content := `
[telemetry]
server = ""
retry_interval = 0
timeout = 0

[targets]
hosts = []
frequencies = []
ports = []

[labels]
environment = ""
service = ""

[logging]
level = ""
`
	tmpFile := createTempFile(t, content)
	defer os.Remove(tmpFile)

	_, err := config.Load(tmpFile)
	if err == nil {
		t.Fatal("Expected error when loading empty config, got nil")
	}

}

// TEST: GIVEN a TOML file with partial configuration
// WHEN the Load function is called
// THEN it should return a error
func TestLoad_PartialConfig(t *testing.T) {
	content := `
[telemetry]
server = ""
retry_interval = 0
timeout = 0

[targets]
hosts = []
frequencies = []
ports = []

[labels]
environment = "staging"
service = ""

[logging]
level = ""
`
	tmpFile := createTempFile(t, content)
	defer os.Remove(tmpFile)

	_, err := config.Load(tmpFile)
	if err != nil {
		t.Fatalf("Expected error when loading partial config, got: %v", err)
	}
}

func createTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	tmpFile := filepath.Join(dir, "config.toml")
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	return tmpFile
}

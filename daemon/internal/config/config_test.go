package config_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bxrne/beacon/daemon/internal/config"
)

// TEST: GIVEN a valid TOML configuration file
// WHEN the Load function is called
// THEN it should correctly parse all fields and return a matching Config struct
func TestLoad_ValidConfig(t *testing.T) {
	content := `
[monitoring]
disk_paths = ["/path1", "/path2"]
frequency = 60

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
		Monitoring: config.Monitoring{
			DiskPaths: []string{"/path1", "/path2"},
			Frequency: 60,
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

// TEST: GIVEN a TOML file with invalid content
// WHEN the Load function is called
// THEN it should return an error
func TestLoad_InvalidContent(t *testing.T) {
	content := `
[monitoring]
disk_paths = ["/path1", "/path2"]
frequency = "invalid"

[labels]
environment = "production"
service = "myapp"

[logging]
level = "info"
`
	tmpFile := createTempFile(t, content)
	defer os.Remove(tmpFile)

	_, err := config.Load(tmpFile)
	if err == nil {
		t.Error("Expected error when loading invalid config, got nil")
	}
}

// TEST: GIVEN an empty TOML file
// WHEN the Load function is called
// THEN it should return a Config struct with zero values
func TestLoad_EmptyConfig(t *testing.T) {
	content := ""
	tmpFile := createTempFile(t, content)
	defer os.Remove(tmpFile)

	cfg, err := config.Load(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load empty config: %v", err)
	}

	expected := &config.Config{}
	if !reflect.DeepEqual(cfg, expected) {
		t.Errorf("Config mismatch\nGot: %+v\nWant: %+v", cfg, expected)
	}
}

// TEST: GIVEN a TOML file with partial configuration
// WHEN the Load function is called
// THEN it should return a Config struct with specified values and zero values for missing fields
func TestLoad_PartialConfig(t *testing.T) {
	content := `
[monitoring]
disk_paths = ["/path1"]
frequency = 30

[labels]
environment = "staging"
`
	tmpFile := createTempFile(t, content)
	defer os.Remove(tmpFile)

	cfg, err := config.Load(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load partial config: %v", err)
	}

	expected := &config.Config{
		Monitoring: config.Monitoring{
			DiskPaths: []string{"/path1"},
			Frequency: 30,
		},
		Labels: config.Labels{
			Environment: "staging",
		},
	}

	if !reflect.DeepEqual(cfg, expected) {
		t.Errorf("Config mismatch\nGot: %+v\nWant: %+v", cfg, expected)
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

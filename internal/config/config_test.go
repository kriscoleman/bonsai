package config

import (
	"os"
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Duration
		wantErr bool
	}{
		// Weeks
		{
			name:    "2 weeks",
			input:   "2w",
			want:    2 * 7 * 24 * time.Hour,
			wantErr: false,
		},
		{
			name:    "1 week",
			input:   "1w",
			want:    7 * 24 * time.Hour,
			wantErr: false,
		},
		{
			name:    "4 weeks",
			input:   "4w",
			want:    4 * 7 * 24 * time.Hour,
			wantErr: false,
		},
		// Days
		{
			name:    "14 days",
			input:   "14d",
			want:    14 * 24 * time.Hour,
			wantErr: false,
		},
		{
			name:    "7 days",
			input:   "7d",
			want:    7 * 24 * time.Hour,
			wantErr: false,
		},
		{
			name:    "30 days",
			input:   "30d",
			want:    30 * 24 * time.Hour,
			wantErr: false,
		},
		// Hours
		{
			name:    "336 hours",
			input:   "336h",
			want:    336 * time.Hour,
			wantErr: false,
		},
		{
			name:    "24 hours",
			input:   "24h",
			want:    24 * time.Hour,
			wantErr: false,
		},
		// Minutes
		{
			name:    "60 minutes",
			input:   "60m",
			want:    60 * time.Minute,
			wantErr: false,
		},
		{
			name:    "1440 minutes",
			input:   "1440m",
			want:    1440 * time.Minute,
			wantErr: false,
		},
		// Seconds
		{
			name:    "3600 seconds",
			input:   "3600s",
			want:    3600 * time.Second,
			wantErr: false,
		},
		// Standard Go duration format (fallback)
		{
			name:    "go format - hours",
			input:   "48h",
			want:    48 * time.Hour,
			wantErr: false,
		},
		{
			name:    "go format - complex",
			input:   "1h30m",
			want:    90 * time.Minute,
			wantErr: false,
		},
		// Error cases
		{
			name:    "invalid format",
			input:   "invalid",
			want:    0,
			wantErr: true,
		},
		{
			name:    "negative value",
			input:   "-2w",
			want:    0,
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    0,
			wantErr: true,
		},
		{
			name:    "no number",
			input:   "w",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDuration(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.LocalAgeThreshold != 2*7*24*time.Hour {
		t.Errorf("DefaultConfig().LocalAgeThreshold = %v, want %v", cfg.LocalAgeThreshold, 2*7*24*time.Hour)
	}

	if cfg.RemoteAgeThreshold != 4*7*24*time.Hour {
		t.Errorf("DefaultConfig().RemoteAgeThreshold = %v, want %v", cfg.RemoteAgeThreshold, 4*7*24*time.Hour)
	}

	if cfg.DryRun != false {
		t.Errorf("DefaultConfig().DryRun = %v, want false", cfg.DryRun)
	}

	if cfg.BulkMode != false {
		t.Errorf("DefaultConfig().BulkMode = %v, want false", cfg.BulkMode)
	}
}

func TestParseDuration_Equivalence(t *testing.T) {
	// Test that different formats produce equivalent durations
	tests := []struct {
		name   string
		input1 string
		input2 string
	}{
		{
			name:   "2 weeks == 14 days",
			input1: "2w",
			input2: "14d",
		},
		{
			name:   "1 day == 24 hours",
			input1: "1d",
			input2: "24h",
		},
		{
			name:   "1 hour == 60 minutes",
			input1: "1h",
			input2: "60m",
		},
		{
			name:   "1 minute == 60 seconds",
			input1: "1m",
			input2: "60s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d1, err1 := ParseDuration(tt.input1)
			if err1 != nil {
				t.Fatalf("ParseDuration(%s) error = %v", tt.input1, err1)
			}

			d2, err2 := ParseDuration(tt.input2)
			if err2 != nil {
				t.Fatalf("ParseDuration(%s) error = %v", tt.input2, err2)
			}

			if d1 != d2 {
				t.Errorf("%s (%v) != %s (%v)", tt.input1, d1, tt.input2, d2)
			}
		})
	}
}

func TestLoadConfigFile(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := tmpDir + "/test-config.yaml"

	configContent := `
local:
  age_threshold: "1w"
remote:
  age_threshold: "3w"
  remote_name: "upstream"
protected_branches:
  - "production"
  - "staging"
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	cfg, err := LoadConfigFile(configPath)
	if err != nil {
		t.Fatalf("LoadConfigFile() error = %v", err)
	}

	expectedLocal := 7 * 24 * time.Hour
	if cfg.LocalAgeThreshold != expectedLocal {
		t.Errorf("LocalAgeThreshold = %v, want %v", cfg.LocalAgeThreshold, expectedLocal)
	}

	expectedRemote := 21 * 24 * time.Hour
	if cfg.RemoteAgeThreshold != expectedRemote {
		t.Errorf("RemoteAgeThreshold = %v, want %v", cfg.RemoteAgeThreshold, expectedRemote)
	}
}

func TestLoadConfigFile_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := tmpDir + "/invalid.yaml"

	invalidContent := `
local:
  age_threshold: "1w"
  invalid yaml here [[[
`

	if err := os.WriteFile(configPath, []byte(invalidContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	_, err := LoadConfigFile(configPath)
	if err == nil {
		t.Error("LoadConfigFile() should return error for invalid YAML")
	}
}

func TestLoadConfigFile_InvalidDuration(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := tmpDir + "/invalid-duration.yaml"

	invalidContent := `
local:
  age_threshold: "invalid"
`

	if err := os.WriteFile(configPath, []byte(invalidContent), 0644); err != nil {
		t.Fatalf("Failed to create test config file: %v", err)
	}

	_, err := LoadConfigFile(configPath)
	if err == nil {
		t.Error("LoadConfigFile() should return error for invalid duration")
	}
}

func TestLoadConfigFile_NonExistent(t *testing.T) {
	_, err := LoadConfigFile("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("LoadConfigFile() should return error for non-existent file")
	}
}

func TestLoadConfig(t *testing.T) {
	// This test just verifies that LoadConfig returns a config
	// without errors (it will use defaults if no file is found)
	cfg := LoadConfig()
	if cfg == nil {
		t.Error("LoadConfig() returned nil")
	}

	// Should have default values
	if cfg.LocalAgeThreshold == 0 {
		t.Error("LoadConfig() returned config with zero LocalAgeThreshold")
	}

	if cfg.RemoteAgeThreshold == 0 {
		t.Error("LoadConfig() returned config with zero RemoteAgeThreshold")
	}
}

func TestFindConfigFile(t *testing.T) {
	// This test verifies the function works without errors
	// The actual path returned depends on the environment
	path := FindConfigFile()

	// Path should be either empty (no config found) or a string
	_ = path
}

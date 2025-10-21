package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	LocalAgeThreshold  time.Duration
	RemoteAgeThreshold time.Duration
	DryRun             bool
	BulkMode           bool
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		LocalAgeThreshold:  2 * 7 * 24 * time.Hour, // 2 weeks
		RemoteAgeThreshold: 4 * 7 * 24 * time.Hour, // 4 weeks
		DryRun:             false,
		BulkMode:           false,
	}
}

// ParseDuration parses a duration string with support for various formats
// Supported formats: 2w, 14d, 336h, 20160m, etc.
func ParseDuration(s string) (time.Duration, error) {
	// Match patterns like "2w", "14d", "336h"
	re := regexp.MustCompile(`^(\d+)([wdhms])$`)
	matches := re.FindStringSubmatch(s)

	if matches == nil {
		// Try standard Go duration format
		return time.ParseDuration(s)
	}

	value, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("invalid duration value: %s", s)
	}

	unit := matches[2]

	switch unit {
	case "w":
		return time.Duration(value) * 7 * 24 * time.Hour, nil
	case "d":
		return time.Duration(value) * 24 * time.Hour, nil
	case "h":
		return time.Duration(value) * time.Hour, nil
	case "m":
		return time.Duration(value) * time.Minute, nil
	case "s":
		return time.Duration(value) * time.Second, nil
	default:
		return 0, fmt.Errorf("unsupported time unit: %s", unit)
	}
}

// FileConfig represents the YAML/JSON structure for configuration files
type FileConfig struct {
	Local struct {
		AgeThreshold string `yaml:"age_threshold"`
	} `yaml:"local"`
	Remote struct {
		AgeThreshold string `yaml:"age_threshold"`
		RemoteName   string `yaml:"remote_name"`
	} `yaml:"remote"`
	ProtectedBranches []string `yaml:"protected_branches"`
}

// LoadConfigFile loads configuration from a file
func LoadConfigFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var fileConfig FileConfig
	if err := yaml.Unmarshal(data, &fileConfig); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	cfg := DefaultConfig()

	// Parse local age threshold if specified
	if fileConfig.Local.AgeThreshold != "" {
		duration, err := ParseDuration(fileConfig.Local.AgeThreshold)
		if err != nil {
			return nil, fmt.Errorf("invalid local age threshold: %w", err)
		}
		cfg.LocalAgeThreshold = duration
	}

	// Parse remote age threshold if specified
	if fileConfig.Remote.AgeThreshold != "" {
		duration, err := ParseDuration(fileConfig.Remote.AgeThreshold)
		if err != nil {
			return nil, fmt.Errorf("invalid remote age threshold: %w", err)
		}
		cfg.RemoteAgeThreshold = duration
	}

	return cfg, nil
}

// FindConfigFile looks for a configuration file in standard locations
// Returns the path to the first config file found, or empty string if none found
func FindConfigFile() string {
	// Check current directory first
	if _, err := os.Stat(".bonsai.yaml"); err == nil {
		return ".bonsai.yaml"
	}
	if _, err := os.Stat(".bonsai.yml"); err == nil {
		return ".bonsai.yml"
	}

	// Check home directory
	home, err := os.UserHomeDir()
	if err == nil {
		homeConfig := filepath.Join(home, ".bonsai.yaml")
		if _, err := os.Stat(homeConfig); err == nil {
			return homeConfig
		}
		homeConfigYml := filepath.Join(home, ".bonsai.yml")
		if _, err := os.Stat(homeConfigYml); err == nil {
			return homeConfigYml
		}
	}

	// Check XDG config directory
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		xdgConfigFile := filepath.Join(xdgConfig, "bonsai", "config.yaml")
		if _, err := os.Stat(xdgConfigFile); err == nil {
			return xdgConfigFile
		}
	}

	return ""
}

// LoadConfig loads configuration from file if it exists, otherwise returns default config
func LoadConfig() *Config {
	configPath := FindConfigFile()
	if configPath == "" {
		return DefaultConfig()
	}

	cfg, err := LoadConfigFile(configPath)
	if err != nil {
		// If config file exists but can't be loaded, fall back to defaults
		return DefaultConfig()
	}

	return cfg
}

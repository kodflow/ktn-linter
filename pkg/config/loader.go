// Package config provides configuration management for KTN linter rules.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	// DefaultConfigFileName is the default configuration file name.
	DefaultConfigFileName string = ".ktn-linter.yaml"
	// AlternateConfigFileName is an alternate configuration file name.
	AlternateConfigFileName string = ".ktn-linter.yml"
)

// Load loads configuration from a file path.
// If path is empty, it searches for default config files in the current directory and parent directories.
func Load(path string) (*Config, error) {
	if path != "" {
		return loadFromFile(path)
	}

	return loadFromDefaultLocations()
}

// loadFromFile loads configuration from a specific file path.
func loadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	// Validate configuration
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid config in %s: %w", path, err)
	}

	return cfg, nil
}

// loadFromDefaultLocations searches for config files in default locations.
func loadFromDefaultLocations() (*Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return DefaultConfig(), nil
	}

	// Search up the directory tree
	dir := cwd
	for {
		// Try default filename
		path := filepath.Join(dir, DefaultConfigFileName)
		if fileExists(path) {
			return loadFromFile(path)
		}

		// Try alternate filename
		path = filepath.Join(dir, AlternateConfigFileName)
		if fileExists(path) {
			return loadFromFile(path)
		}

		// Move to parent directory
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	// No config file found, return default
	return DefaultConfig(), nil
}

// fileExists checks if a file exists.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// validateConfig validates the configuration.
func validateConfig(cfg *Config) error {
	if cfg == nil {
		return nil
	}

	// Validate version
	if cfg.Version != 0 && cfg.Version != 1 {
		return fmt.Errorf("unsupported config version: %d", cfg.Version)
	}

	// Validate rules
	for code, ruleCfg := range cfg.Rules {
		if ruleCfg == nil {
			continue
		}

		// Validate threshold is positive if set
		if ruleCfg.Threshold != nil && *ruleCfg.Threshold < 0 {
			return fmt.Errorf("rule %s: threshold must be non-negative, got %d", code, *ruleCfg.Threshold)
		}

		// Validate exclusion patterns
		for _, pattern := range ruleCfg.Exclude {
			if pattern == "" {
				return fmt.Errorf("rule %s: empty exclusion pattern", code)
			}
		}
	}

	// Validate global exclusions
	for _, pattern := range cfg.Exclude {
		if pattern == "" {
			return fmt.Errorf("empty global exclusion pattern")
		}
	}

	return nil
}

// LoadAndSet loads configuration and sets it as the global config.
func LoadAndSet(path string) error {
	cfg, err := Load(path)
	if err != nil {
		return err
	}

	Set(cfg)

	return nil
}

// MustLoad loads configuration and panics on error.
func MustLoad(path string) *Config {
	cfg, err := Load(path)
	if err != nil {
		panic(err)
	}

	return cfg
}

// SaveToFile saves configuration to a file.
func SaveToFile(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", path, err)
	}

	return nil
}

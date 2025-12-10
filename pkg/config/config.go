// Package config provides configuration management for KTN linter rules.
package config

import (
	"path/filepath"
	"strings"
	"sync"
)

// Config represents the complete linter configuration.
type Config struct {
	// Version of the configuration format
	Version int `yaml:"version"`

	// Exclude contains global file exclusion patterns (apply to all rules)
	Exclude []string `yaml:"exclude,omitempty"`

	// Rules contains per-rule configuration
	Rules map[string]*RuleConfig `yaml:"rules,omitempty"`

	// compiledExcludes caches compiled glob patterns
	compiledExcludes []string
	// mu protects compiledExcludes
	mu sync.RWMutex
}

// RuleConfig represents configuration for a single rule.
type RuleConfig struct {
	// Enabled indicates whether the rule is active (default: true)
	Enabled *bool `yaml:"enabled,omitempty"`

	// Threshold is a numeric threshold for rules that support it
	Threshold *int `yaml:"threshold,omitempty"`

	// Exclude contains rule-specific file exclusion patterns
	Exclude []string `yaml:"exclude,omitempty"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		Version: 1,
		Exclude: []string{},
		Rules:   make(map[string]*RuleConfig),
	}
}

// globalConfig is the singleton configuration instance.
var (
	globalConfig     *Config
	globalConfigOnce sync.Once
	globalConfigMu   sync.RWMutex
)

// Get returns the global configuration instance.
func Get() *Config {
	globalConfigMu.RLock()
	cfg := globalConfig
	globalConfigMu.RUnlock()

	if cfg != nil {
		return cfg
	}

	globalConfigOnce.Do(func() {
		globalConfigMu.Lock()
		if globalConfig == nil {
			globalConfig = DefaultConfig()
		}
		globalConfigMu.Unlock()
	})

	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	return globalConfig
}

// Set sets the global configuration instance.
func Set(cfg *Config) {
	globalConfigMu.Lock()
	defer globalConfigMu.Unlock()
	globalConfig = cfg
}

// Reset resets the global configuration to default.
func Reset() {
	globalConfigMu.Lock()
	defer globalConfigMu.Unlock()
	globalConfig = DefaultConfig()
	globalConfigOnce = sync.Once{}
}

// IsRuleEnabled checks if a rule is enabled.
func (c *Config) IsRuleEnabled(ruleCode string) bool {
	if c == nil || c.Rules == nil {
		return true
	}

	ruleCfg, exists := c.Rules[ruleCode]
	if !exists || ruleCfg == nil || ruleCfg.Enabled == nil {
		return true
	}

	return *ruleCfg.Enabled
}

// GetThreshold returns the threshold for a rule, or the default if not set.
func (c *Config) GetThreshold(ruleCode string, defaultValue int) int {
	if c == nil || c.Rules == nil {
		return defaultValue
	}

	ruleCfg, exists := c.Rules[ruleCode]
	if !exists || ruleCfg == nil || ruleCfg.Threshold == nil {
		return defaultValue
	}

	return *ruleCfg.Threshold
}

// IsFileExcluded checks if a file should be excluded for a specific rule.
func (c *Config) IsFileExcluded(ruleCode, filename string) bool {
	if c == nil {
		return false
	}

	// Check global exclusions first
	if c.matchesAnyPattern(filename, c.Exclude) {
		return true
	}

	// Check rule-specific exclusions
	if c.Rules != nil {
		if ruleCfg, exists := c.Rules[ruleCode]; exists && ruleCfg != nil {
			if c.matchesAnyPattern(filename, ruleCfg.Exclude) {
				return true
			}
		}
	}

	return false
}

// IsFileExcludedGlobally checks if a file is excluded globally (all rules).
func (c *Config) IsFileExcludedGlobally(filename string) bool {
	if c == nil {
		return false
	}

	return c.matchesAnyPattern(filename, c.Exclude)
}

// matchesAnyPattern checks if filename matches any of the given patterns.
func (c *Config) matchesAnyPattern(filename string, patterns []string) bool {
	if len(patterns) == 0 {
		return false
	}

	// Normalize filename to use forward slashes
	normalizedFilename := filepath.ToSlash(filename)
	baseFilename := filepath.Base(filename)

	for _, pattern := range patterns {
		// Normalize pattern
		normalizedPattern := filepath.ToSlash(pattern)

		// Check if pattern matches the full path
		if matched, _ := filepath.Match(normalizedPattern, normalizedFilename); matched {
			return true
		}

		// Check if pattern matches just the basename
		if matched, _ := filepath.Match(normalizedPattern, baseFilename); matched {
			return true
		}

		// Check for ** glob patterns (recursive matching)
		if strings.Contains(normalizedPattern, "**") {
			if c.matchDoubleStarPattern(normalizedFilename, normalizedPattern) {
				return true
			}
		}

		// Check suffix match (e.g., "*_test.go" should match "foo/bar_test.go")
		if strings.HasPrefix(normalizedPattern, "*") {
			suffix := normalizedPattern[1:]
			if strings.HasSuffix(normalizedFilename, suffix) || strings.HasSuffix(baseFilename, suffix) {
				return true
			}
		}
	}

	return false
}

// matchDoubleStarPattern handles ** glob patterns for recursive matching.
func (c *Config) matchDoubleStarPattern(filename, pattern string) bool {
	// Split pattern by **
	parts := strings.Split(pattern, "**")

	// Handle patterns with multiple ** (e.g., **/testdata/**)
	if len(parts) == 3 {
		// Pattern like **/testdata/**
		middle := strings.Trim(parts[1], "/")
		if middle == "" {
			return true
		}
		// Check if middle part appears in the path
		return strings.Contains(filename, "/"+middle+"/") || strings.HasPrefix(filename, middle+"/")
	}

	if len(parts) != 2 {
		return false
	}

	prefix := strings.TrimSuffix(parts[0], "/")
	suffix := strings.TrimPrefix(parts[1], "/")

	// Handle prefix matching
	if prefix != "" {
		// Check if filename contains the prefix as a path component
		if !strings.HasPrefix(filename, prefix+"/") && !strings.Contains(filename, "/"+prefix+"/") {
			return false
		}
	}

	// Handle suffix matching
	if suffix != "" {
		// Check for path component match
		suffixWithSlash := "/" + suffix
		if !strings.HasSuffix(filename, suffixWithSlash) && !strings.Contains(filename, suffixWithSlash+"/") {
			return false
		}
	}

	return true
}

// Merge merges another config into this one (other takes precedence).
func (c *Config) Merge(other *Config) {
	if other == nil {
		return
	}

	// Merge global exclusions
	c.Exclude = append(c.Exclude, other.Exclude...)

	// Merge rules
	if other.Rules != nil {
		if c.Rules == nil {
			c.Rules = make(map[string]*RuleConfig)
		}
		for code, ruleCfg := range other.Rules {
			if existing, exists := c.Rules[code]; exists && existing != nil {
				// Merge rule config
				if ruleCfg.Enabled != nil {
					existing.Enabled = ruleCfg.Enabled
				}
				if ruleCfg.Threshold != nil {
					existing.Threshold = ruleCfg.Threshold
				}
				existing.Exclude = append(existing.Exclude, ruleCfg.Exclude...)
			} else {
				c.Rules[code] = ruleCfg
			}
		}
	}
}

// Bool is a helper to create a pointer to a bool.
func Bool(v bool) *bool {
	return &v
}

// Int is a helper to create a pointer to an int.
func Int(v int) *int {
	return &v
}

// Package config provides configuration management for KTN linter rules.
// It allows enabling/disabling rules, setting thresholds, and excluding files.
package config

import (
	"path/filepath"
	"strings"
	"sync"
)

const (
	// defaultRulesMapCapacity is the default capacity for rules map.
	defaultRulesMapCapacity int = 10
	// doubleStarPartsThree is parts count for ** with 3 components.
	doubleStarPartsThree int = 3
	// doubleStarPartsTwo is parts count for ** with 2 components.
	doubleStarPartsTwo int = 2
)

// Config represents the complete linter configuration.
// It contains global settings and per-rule configurations.
type Config struct {
	// Version of the configuration format
	Version int `yaml:"version"`

	// Exclude contains global file exclusion patterns (apply to all rules)
	Exclude []string `yaml:"exclude,omitempty"`

	// Rules contains per-rule configuration
	Rules map[string]*RuleConfig `yaml:"rules,omitempty"`

	// Verbose enables verbose message output with examples
	Verbose bool `yaml:"-"`

	// compiledExcludes caches compiled glob patterns
	compiledExcludes []string
	// mu protects compiledExcludes
	mu sync.RWMutex
}

// RuleConfig represents configuration for a single rule.
// It allows customizing rule behavior per-rule basis.
type RuleConfig struct {
	// Enabled indicates whether the rule is active (default: true)
	Enabled *bool `yaml:"enabled,omitempty"`

	// Threshold is a numeric threshold for rules that support it
	Threshold *int `yaml:"threshold,omitempty"`

	// Exclude contains rule-specific file exclusion patterns
	Exclude []string `yaml:"exclude,omitempty"`
}

// DefaultConfig returns the default configuration.
//
// Returns:
//   - *Config: default configuration instance with version 1
func DefaultConfig() *Config {
	// Return default config
	return &Config{
		Version: 1,
		Exclude: []string{},
		Rules:   make(map[string]*RuleConfig, defaultRulesMapCapacity),
	}
}

// globalConfig is the singleton configuration instance.
var (
	globalConfig     *Config
	globalConfigOnce *sync.Once = &sync.Once{}
	globalConfigMu   sync.RWMutex
)

// Get returns the global configuration instance.
//
// Returns:
//   - *Config: global configuration instance, never nil
func Get() *Config {
	globalConfigMu.RLock()
	cfg := globalConfig
	globalConfigMu.RUnlock()

	// Check if config already exists
	if cfg != nil {
		// Return existing config
		return cfg
	}

	globalConfigOnce.Do(func() {
		globalConfigMu.Lock()
		// Double-check after acquiring write lock
		if globalConfig == nil {
			globalConfig = DefaultConfig()
		}
		globalConfigMu.Unlock()
	})

	globalConfigMu.RLock()
	defer globalConfigMu.RUnlock()
	// Return global config
	return globalConfig
}

// Set sets the global configuration instance.
//
// Params:
//   - cfg: configuration to set
//
// Returns: none
func Set(cfg *Config) {
	globalConfigMu.Lock()
	defer globalConfigMu.Unlock()
	globalConfig = cfg
}

// Reset resets the global configuration to default.
//
// Returns: none
func Reset() {
	globalConfigMu.Lock()
	defer globalConfigMu.Unlock()
	globalConfig = DefaultConfig()
	// Reset the Once to allow re-initialization
	globalConfigOnce = &sync.Once{}
}

// IsRuleEnabled checks if a rule is enabled.
//
// Params:
//   - ruleCode: the rule code to check
//
// Returns:
//   - bool: true if the rule is enabled
func (c *Config) IsRuleEnabled(ruleCode string) bool {
	// Check nil config
	if c == nil || c.Rules == nil {
		// Return enabled by default
		return true
	}

	ruleCfg, exists := c.Rules[ruleCode]
	// Check if rule config exists
	if !exists || ruleCfg == nil || ruleCfg.Enabled == nil {
		// Return enabled by default
		return true
	}

	// Return configured value
	return *ruleCfg.Enabled
}

// GetThreshold returns the threshold for a rule, or the default if not set.
//
// Params:
//   - ruleCode: the rule code to check
//   - defaultValue: default value if not configured
//
// Returns:
//   - int: the threshold value
func (c *Config) GetThreshold(ruleCode string, defaultValue int) int {
	// Check nil config
	if c == nil || c.Rules == nil {
		// Return default value
		return defaultValue
	}

	ruleCfg, exists := c.Rules[ruleCode]
	// Check if rule config exists
	if !exists || ruleCfg == nil || ruleCfg.Threshold == nil {
		// Return default value
		return defaultValue
	}

	// Return configured threshold
	return *ruleCfg.Threshold
}

// IsFileExcluded checks if a file should be excluded for a specific rule.
//
// Params:
//   - ruleCode: the rule code to check
//   - filename: the file path to check
//
// Returns:
//   - bool: true if the file should be excluded
func (c *Config) IsFileExcluded(ruleCode, filename string) bool {
	// Check nil config
	if c == nil {
		// Return not excluded
		return false
	}

	// Check global exclusions first
	if c.matchesAnyPattern(filename, c.Exclude) {
		// Return excluded by global pattern
		return true
	}

	// Check rule-specific exclusions
	if c.Rules != nil {
		// Check if rule has specific exclusions
		if ruleCfg, exists := c.Rules[ruleCode]; exists && ruleCfg != nil {
			// Check rule-specific patterns
			if c.matchesAnyPattern(filename, ruleCfg.Exclude) {
				// Return excluded by rule pattern
				return true
			}
		}
	}

	// Return not excluded
	return false
}

// IsFileExcludedGlobally checks if a file is excluded globally (all rules).
//
// Params:
//   - filename: the file path to check
//
// Returns:
//   - bool: true if the file is globally excluded
func (c *Config) IsFileExcludedGlobally(filename string) bool {
	// Check nil config
	if c == nil {
		// Return not excluded
		return false
	}

	// Return global pattern match result
	return c.matchesAnyPattern(filename, c.Exclude)
}

// matchesAnyPattern checks if filename matches any of the given patterns.
//
// Params:
//   - filename: the file path to check
//   - patterns: list of glob patterns
//
// Returns:
//   - bool: true if any pattern matches
func (c *Config) matchesAnyPattern(filename string, patterns []string) bool {
	// Check empty patterns
	if len(patterns) == 0 {
		// Return no match
		return false
	}

	// Normalize filename to use forward slashes
	normalizedFilename := filepath.ToSlash(filename)
	baseFilename := filepath.Base(filename)

	// Iterate over patterns
	for _, pattern := range patterns {
		// Normalize pattern
		normalizedPattern := filepath.ToSlash(pattern)

		// Check if pattern matches the full path
		if matched, _ := filepath.Match(normalizedPattern, normalizedFilename); matched {
			// Return matched
			return true
		}

		// Check if pattern matches just the basename
		if matched, _ := filepath.Match(normalizedPattern, baseFilename); matched {
			// Return matched
			return true
		}

		// Check for ** glob patterns (recursive matching)
		if strings.Contains(normalizedPattern, "**") {
			// Check double star pattern
			if c.matchDoubleStarPattern(normalizedFilename, normalizedPattern) {
				// Return matched
				return true
			}
		}

		// Check suffix match (e.g., "*_test.go" should match "foo/bar_test.go")
		if strings.HasPrefix(normalizedPattern, "*") {
			suffix := normalizedPattern[1:]
			// Check suffix match
			if strings.HasSuffix(normalizedFilename, suffix) || strings.HasSuffix(baseFilename, suffix) {
				// Return matched
				return true
			}
		}
	}

	// Return no match
	return false
}

// matchDoubleStarPattern handles ** glob patterns for recursive matching.
//
// Params:
//   - filename: the file path to check
//   - pattern: the glob pattern with **
//
// Returns:
//   - bool: true if pattern matches
func (c *Config) matchDoubleStarPattern(filename, pattern string) bool {
	// Split pattern by **
	parts := strings.Split(pattern, "**")

	// Handle patterns with multiple ** (e.g., **/testdata/**)
	if len(parts) == doubleStarPartsThree {
		// Check triple pattern
		return c.matchTriplePattern(parts, filename)
	}

	// Check invalid pattern
	if len(parts) != doubleStarPartsTwo {
		// Return no match
		return false
	}

	// Check double pattern with prefix/suffix
	return c.matchPrefixSuffix(parts, filename)
}

// matchTriplePattern handles **/middle/** patterns.
//
// Params:
//   - parts: split pattern parts
//   - filename: the file path to check
//
// Returns:
//   - bool: true if pattern matches
func (c *Config) matchTriplePattern(parts []string, filename string) bool {
	// Pattern like **/testdata/**
	middle := strings.Trim(parts[1], "/")
	// Check empty middle
	if middle == "" {
		// Return matched
		return true
	}
	// Check if middle part appears in the path
	return strings.Contains(filename, "/"+middle+"/") || strings.HasPrefix(filename, middle+"/")
}

// matchPrefixSuffix handles prefix/**/suffix patterns.
//
// Params:
//   - parts: split pattern parts
//   - filename: the file path to check
//
// Returns:
//   - bool: true if pattern matches
func (c *Config) matchPrefixSuffix(parts []string, filename string) bool {
	prefix := strings.TrimSuffix(parts[0], "/")
	suffix := strings.TrimPrefix(parts[1], "/")

	// Handle prefix matching
	if prefix != "" {
		// Check if filename contains the prefix as a path component
		if !strings.HasPrefix(filename, prefix+"/") && !strings.Contains(filename, "/"+prefix+"/") {
			// Return no match
			return false
		}
	}

	// Handle suffix matching
	if suffix != "" {
		// Check for path component match
		suffixWithSlash := "/" + suffix
		// Check suffix presence
		if !strings.HasSuffix(filename, suffixWithSlash) && !strings.Contains(filename, suffixWithSlash+"/") {
			// Return no match
			return false
		}
	}

	// Return matched
	return true
}

// Merge merges another config into this one (other takes precedence).
//
// Params:
//   - other: configuration to merge
//
// Returns: none
func (c *Config) Merge(other *Config) {
	// Check nil other
	if other == nil {
		// Return early
		return
	}

	// Merge global exclusions
	c.Exclude = append(c.Exclude, other.Exclude...)

	// Merge rules
	if other.Rules != nil {
		// Initialize rules map if needed
		if c.Rules == nil {
			c.Rules = make(map[string]*RuleConfig, len(other.Rules))
		}
		// Iterate over other rules
		for code, ruleCfg := range other.Rules {
			// Check if rule exists
			if existing, exists := c.Rules[code]; exists && existing != nil {
				// Merge rule config
				if ruleCfg.Enabled != nil {
					existing.Enabled = ruleCfg.Enabled
				}
				// Merge threshold
				if ruleCfg.Threshold != nil {
					existing.Threshold = ruleCfg.Threshold
				}
				existing.Exclude = append(existing.Exclude, ruleCfg.Exclude...)
			} else {
				// Add new rule
				c.Rules[code] = ruleCfg
			}
		}
	}
}

// Bool is a helper to create a pointer to a bool.
//
// Params:
//   - v: boolean value
//
// Returns:
//   - *bool: pointer to the value
func Bool(v bool) *bool {
	// Return pointer
	return &v
}

// Int is a helper to create a pointer to an int.
//
// Params:
//   - v: integer value
//
// Returns:
//   - *int: pointer to the value
func Int(v int) *int {
	// Return pointer
	return &v
}

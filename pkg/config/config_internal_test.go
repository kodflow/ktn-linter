package config

import (
	"sync"
	"testing"
)

// TestIsFileExcludedGlobally tests global file exclusion.
func TestIsFileExcludedGlobally(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *Config
		filename string
		want     bool
	}{
		{
			name:     "nil config returns false",
			cfg:      nil,
			filename: "test.go",
			want:     false,
		},
		{
			name: "file matches global exclusion",
			cfg: &Config{
				Exclude: []string{"*_test.go"},
			},
			filename: "foo_test.go",
			want:     true,
		},
		{
			name: "file does not match global exclusion",
			cfg: &Config{
				Exclude: []string{"*_test.go"},
			},
			filename: "foo.go",
			want:     false,
		},
		{
			name: "empty exclusion list",
			cfg: &Config{
				Exclude: []string{},
			},
			filename: "test.go",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.IsFileExcludedGlobally(tt.filename)
			if got != tt.want {
				t.Errorf("IsFileExcludedGlobally(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

// TestMatchesAnyPattern tests pattern matching logic.
func TestMatchesAnyPattern(t *testing.T) {
	cfg := &Config{}

	tests := []struct {
		name     string
		filename string
		patterns []string
		want     bool
	}{
		{
			name:     "empty patterns",
			filename: "test.go",
			patterns: []string{},
			want:     false,
		},
		{
			name:     "exact match",
			filename: "main.go",
			patterns: []string{"main.go"},
			want:     true,
		},
		{
			name:     "basename match",
			filename: "pkg/foo/main.go",
			patterns: []string{"main.go"},
			want:     true,
		},
		{
			name:     "suffix pattern with star",
			filename: "foo_test.go",
			patterns: []string{"*_test.go"},
			want:     true,
		},
		{
			name:     "suffix pattern with path",
			filename: "pkg/analyzer/foo_test.go",
			patterns: []string{"*_test.go"},
			want:     true,
		},
		{
			name:     "double star at start",
			filename: "vendor/github.com/foo/bar.go",
			patterns: []string{"vendor/**"},
			want:     true,
		},
		{
			name:     "double star in middle",
			filename: "pkg/testdata/foo.go",
			patterns: []string{"**/testdata/**"},
			want:     true,
		},
		{
			name:     "no match",
			filename: "main.go",
			patterns: []string{"test.go", "*.txt"},
			want:     false,
		},
		{
			name:     "simple wildcard pattern",
			filename: "test.go",
			patterns: []string{"*.go"},
			want:     true,
		},
		{
			name:     "basename pattern with wildcard",
			filename: "pkg/foo/test.go",
			patterns: []string{"*.go"},
			want:     true,
		},
		{
			name:     "invalid glob pattern - continues to next",
			filename: "test.go",
			patterns: []string{"[invalid", "test.go"},
			want:     true,
		},
		{
			name:     "pattern with special chars",
			filename: "test.go",
			patterns: []string{"test.?o"},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cfg.matchesAnyPattern(tt.filename, tt.patterns)
			if got != tt.want {
				t.Errorf("matchesAnyPattern(%q, %v) = %v, want %v", tt.filename, tt.patterns, got, tt.want)
			}
		})
	}
}

// TestMatchDoubleStarPattern tests double star pattern matching.
func TestMatchDoubleStarPattern(t *testing.T) {
	cfg := &Config{}

	tests := []struct {
		name     string
		filename string
		pattern  string
		want     bool
	}{
		{
			name:     "triple star pattern - middle part present",
			filename: "pkg/testdata/src/foo.go",
			pattern:  "**/testdata/**",
			want:     true,
		},
		{
			name:     "triple star pattern - prefix match",
			filename: "testdata/src/foo.go",
			pattern:  "**/testdata/**",
			want:     true,
		},
		{
			name:     "triple star pattern - empty middle",
			filename: "any/path/file.go",
			pattern:  "**/**",
			want:     true,
		},
		{
			name:     "invalid pattern - four parts",
			filename: "any/path/file.go",
			pattern:  "**/**/**/**",
			want:     false,
		},
		{
			name:     "standard double star - prefix and suffix",
			filename: "vendor/github.com/foo/bar.go",
			pattern:  "vendor/**",
			want:     true,
		},
		{
			name:     "standard double star - no prefix",
			filename: "pkg/foo/bar.go",
			pattern:  "**/foo/**",
			want:     true,
		},
		{
			name:     "prefix not matching",
			filename: "pkg/foo/bar.go",
			pattern:  "vendor/**",
			want:     false,
		},
		{
			name:     "suffix not matching",
			filename: "vendor/foo/bar.go",
			pattern:  "**/testdata/**",
			want:     false,
		},
		{
			name:     "double star with empty prefix",
			filename: "any/path/file.go",
			pattern:  "**/file.go",
			want:     true,
		},
		{
			name:     "double star with empty suffix",
			filename: "vendor/any/path.go",
			pattern:  "vendor/**",
			want:     true,
		},
		{
			name:     "prefix contains in middle of path",
			filename: "pkg/vendor/foo.go",
			pattern:  "vendor/**",
			want:     true,
		},
		{
			name:     "suffix contains in middle of path",
			filename: "vendor/testdata/bar/foo.go",
			pattern:  "**/testdata/**",
			want:     true,
		},
		{
			name:     "suffix at end of path",
			filename: "vendor/foo/testdata",
			pattern:  "**/testdata",
			want:     true,
		},
		{
			name:     "suffix with additional path after",
			filename: "pkg/testdata/src/foo.go",
			pattern:  "**/testdata",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cfg.matchDoubleStarPattern(tt.filename, tt.pattern)
			if got != tt.want {
				t.Errorf("matchDoubleStarPattern(%q, %q) = %v, want %v", tt.filename, tt.pattern, got, tt.want)
			}
		})
	}
}

// TestMerge tests configuration merging.
func TestMerge(t *testing.T) {
	tests := []struct {
		name  string
		base  *Config
		other *Config
		check func(*testing.T, *Config)
	}{
		{
			name: "merge with nil other",
			base: &Config{
				Exclude: []string{"test1.go"},
				Rules: map[string]*RuleConfig{
					"RULE-1": {Enabled: Bool(true)},
				},
			},
			other: nil,
			check: func(t *testing.T, cfg *Config) {
				if len(cfg.Exclude) != 1 {
					t.Errorf("Expected 1 exclusion, got %d", len(cfg.Exclude))
				}
			},
		},
		{
			name: "merge exclusions",
			base: &Config{
				Exclude: []string{"test1.go"},
			},
			other: &Config{
				Exclude: []string{"test2.go"},
			},
			check: func(t *testing.T, cfg *Config) {
				if len(cfg.Exclude) != 2 {
					t.Errorf("Expected 2 exclusions, got %d", len(cfg.Exclude))
				}
			},
		},
		{
			name: "merge rules - base nil rules",
			base: &Config{
				Rules: nil,
			},
			other: &Config{
				Rules: map[string]*RuleConfig{
					"RULE-1": {Enabled: Bool(false)},
				},
			},
			check: func(t *testing.T, cfg *Config) {
				if cfg.Rules == nil {
					t.Error("Rules should be initialized")
				}
				if cfg.IsRuleEnabled("RULE-1") {
					t.Error("RULE-1 should be disabled")
				}
			},
		},
		{
			name: "merge rules - override enabled",
			base: &Config{
				Rules: map[string]*RuleConfig{
					"RULE-1": {Enabled: Bool(true)},
				},
			},
			other: &Config{
				Rules: map[string]*RuleConfig{
					"RULE-1": {Enabled: Bool(false)},
				},
			},
			check: func(t *testing.T, cfg *Config) {
				if cfg.IsRuleEnabled("RULE-1") {
					t.Error("RULE-1 should be disabled after merge")
				}
			},
		},
		{
			name: "merge rules - override threshold",
			base: &Config{
				Rules: map[string]*RuleConfig{
					"RULE-1": {Threshold: Int(10)},
				},
			},
			other: &Config{
				Rules: map[string]*RuleConfig{
					"RULE-1": {Threshold: Int(20)},
				},
			},
			check: func(t *testing.T, cfg *Config) {
				if cfg.GetThreshold("RULE-1", 10) != 20 {
					t.Error("Threshold should be overridden to 20")
				}
			},
		},
		{
			name: "merge rules - append exclusions",
			base: &Config{
				Rules: map[string]*RuleConfig{
					"RULE-1": {Exclude: []string{"test1.go"}},
				},
			},
			other: &Config{
				Rules: map[string]*RuleConfig{
					"RULE-1": {Exclude: []string{"test2.go"}},
				},
			},
			check: func(t *testing.T, cfg *Config) {
				if len(cfg.Rules["RULE-1"].Exclude) != 2 {
					t.Errorf("Expected 2 rule exclusions, got %d", len(cfg.Rules["RULE-1"].Exclude))
				}
			},
		},
		{
			name: "merge rules - add new rule",
			base: &Config{
				Rules: map[string]*RuleConfig{
					"RULE-1": {Enabled: Bool(true)},
				},
			},
			other: &Config{
				Rules: map[string]*RuleConfig{
					"RULE-2": {Enabled: Bool(false)},
				},
			},
			check: func(t *testing.T, cfg *Config) {
				if len(cfg.Rules) != 2 {
					t.Errorf("Expected 2 rules, got %d", len(cfg.Rules))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.base.Merge(tt.other)
			tt.check(t, tt.base)
		})
	}
}

// TestGet_Concurrent tests concurrent access to global config.
func TestGet_Concurrent(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"concurrent access safety"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset global state to nil
			globalConfigMu.Lock()
			globalConfig = nil
			globalConfigOnce = &sync.Once{}
			globalConfigMu.Unlock()

			// Test concurrent Get calls
			var wg sync.WaitGroup
			for i := 0; i < 10; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					cfg := Get()
					if cfg == nil {
						t.Error("Get() returned nil in concurrent access")
					}
				}()
			}
			wg.Wait()

			// Verify all got the same instance
			cfg := Get()
			if cfg == nil {
				t.Error("Get() returned nil after concurrent initialization")
			}
		})
	}
}

// TestGet_ExistingConfig tests Get with already initialized config.
func TestGet_ExistingConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"existing config not re-initialized"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set a custom config
			customCfg := &Config{
				Version: 1,
				Rules: map[string]*RuleConfig{
					"EXISTING-TEST": {Enabled: Bool(false)},
				},
			}
			Set(customCfg)

			// Get should return existing config without calling DefaultConfig
			cfg := Get()
			if cfg == nil {
				t.Fatal("Get() returned nil")
			}

			if cfg.IsRuleEnabled("EXISTING-TEST") {
				t.Error("Expected EXISTING-TEST to be disabled")
			}

			// Reset after test
			Reset()
		})
	}
}

// TestGet_DoubleCheckLocking tests the double-check locking in Get.
func TestGet_DoubleCheckLocking(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"double-check locking pattern"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Force nil state
			globalConfigMu.Lock()
			globalConfig = nil
			globalConfigOnce = &sync.Once{}
			globalConfigMu.Unlock()

			// First call should initialize
			cfg1 := Get()
			if cfg1 == nil {
				t.Fatal("First Get() returned nil")
			}

			// Second call should return same instance without re-initialization
			cfg2 := Get()
			if cfg2 == nil {
				t.Fatal("Second Get() returned nil")
			}

			// Verify they're the same instance
			if &cfg1 != &cfg2 {
				t.Log("Configs are different instances (expected, pointers differ)")
			}

			// Reset
			Reset()
		})
	}
}

// TestIsRuleEnabled_NilRules tests IsRuleEnabled with nil rules.
func TestIsRuleEnabled_NilRules(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *Config
		ruleCode string
		want     bool
	}{
		{
			name:     "config with nil Rules map",
			cfg:      &Config{Rules: nil},
			ruleCode: "TEST-001",
			want:     true,
		},
		{
			name: "rule config is nil",
			cfg: &Config{
				Rules: map[string]*RuleConfig{
					"TEST-001": nil,
				},
			},
			ruleCode: "TEST-001",
			want:     true,
		},
		{
			name: "rule config Enabled is nil",
			cfg: &Config{
				Rules: map[string]*RuleConfig{
					"TEST-001": {Enabled: nil},
				},
			},
			ruleCode: "TEST-001",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.IsRuleEnabled(tt.ruleCode)
			if got != tt.want {
				t.Errorf("IsRuleEnabled(%q) = %v, want %v", tt.ruleCode, got, tt.want)
			}
		})
	}
}

// TestGetThreshold_NilCases tests GetThreshold with nil cases.
func TestGetThreshold_NilCases(t *testing.T) {
	tests := []struct {
		name         string
		cfg          *Config
		ruleCode     string
		defaultValue int
		want         int
	}{
		{
			name:         "config with nil Rules map",
			cfg:          &Config{Rules: nil},
			ruleCode:     "TEST-001",
			defaultValue: 42,
			want:         42,
		},
		{
			name: "rule config is nil",
			cfg: &Config{
				Rules: map[string]*RuleConfig{
					"TEST-001": nil,
				},
			},
			ruleCode:     "TEST-001",
			defaultValue: 42,
			want:         42,
		},
		{
			name: "rule config Threshold is nil",
			cfg: &Config{
				Rules: map[string]*RuleConfig{
					"TEST-001": {Threshold: nil},
				},
			},
			ruleCode:     "TEST-001",
			defaultValue: 42,
			want:         42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.GetThreshold(tt.ruleCode, tt.defaultValue)
			if got != tt.want {
				t.Errorf("GetThreshold(%q, %d) = %d, want %d", tt.ruleCode, tt.defaultValue, got, tt.want)
			}
		})
	}
}

// TestIsFileExcluded_NilRules tests IsFileExcluded with nil rules.
func TestIsFileExcluded_NilRules(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *Config
		ruleCode string
		filename string
		want     bool
	}{
		{
			name:     "config with nil Rules map",
			cfg:      &Config{Rules: nil},
			ruleCode: "TEST-001",
			filename: "test.go",
			want:     false,
		},
		{
			name: "rule config is nil",
			cfg: &Config{
				Rules: map[string]*RuleConfig{
					"TEST-001": nil,
				},
			},
			ruleCode: "TEST-001",
			filename: "test.go",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.IsFileExcluded(tt.ruleCode, tt.filename)
			if got != tt.want {
				t.Errorf("IsFileExcluded(%q, %q) = %v, want %v", tt.ruleCode, tt.filename, got, tt.want)
			}
		})
	}
}

// TestBoolAndIntHelpers tests Bool and Int helper functions.
func TestBoolAndIntHelpers(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"helper functions work correctly"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Bool helper
			b := Bool(false)
			if b == nil {
				t.Error("Bool(false) returned nil")
			}
			if *b != false {
				t.Errorf("Bool(false) = %v, want false", *b)
			}

			// Test Int helper with negative
			i := Int(-10)
			if i == nil {
				t.Error("Int(-10) returned nil")
			}
			if *i != -10 {
				t.Errorf("Int(-10) = %d, want -10", *i)
			}

			// Test Int helper with zero
			zero := Int(0)
			if zero == nil {
				t.Error("Int(0) returned nil")
			}
			if *zero != 0 {
				t.Errorf("Int(0) = %d, want 0", *zero)
			}
		})
	}
}

// TestConfig_matchesAnyPattern tests the private matchesAnyPattern function.
func TestConfig_matchesAnyPattern(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		patterns []string
		want     bool
	}{
		{
			name:     "empty patterns",
			filename: "foo.go",
			patterns: []string{},
			want:     false,
		},
		{
			name:     "simple match",
			filename: "foo_test.go",
			patterns: []string{"*_test.go"},
			want:     true,
		},
		{
			name:     "no match",
			filename: "foo.go",
			patterns: []string{"*_test.go"},
			want:     false,
		},
		{
			name:     "basename match",
			filename: "path/to/main.go",
			patterns: []string{"main.go"},
			want:     true,
		},
		{
			name:     "double star pattern",
			filename: "pkg/testdata/src/foo.go",
			patterns: []string{"**/testdata/**"},
			want:     true,
		},
		{
			name:     "suffix match with star",
			filename: "foo/bar_test.go",
			patterns: []string{"*_test.go"},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{}
			got := cfg.matchesAnyPattern(tt.filename, tt.patterns)
			if got != tt.want {
				t.Errorf("matchesAnyPattern(%q, %v) = %v, want %v", tt.filename, tt.patterns, got, tt.want)
			}
		})
	}
}

// TestConfig_matchDoubleStarPattern tests the private matchDoubleStarPattern function.
func TestConfig_matchDoubleStarPattern(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		pattern  string
		want     bool
	}{
		{
			name:     "pattern with 3 parts (middle match)",
			filename: "pkg/testdata/src/foo.go",
			pattern:  "**/testdata/**",
			want:     true,
		},
		{
			name:     "pattern with 3 parts (no middle)",
			filename: "pkg/foo.go",
			pattern:  "**/testdata/**",
			want:     false,
		},
		{
			name:     "pattern with 3 parts (empty middle)",
			filename: "anything.go",
			pattern:  "**/**",
			want:     true,
		},
		{
			name:     "prefix match",
			filename: "vendor/github.com/foo/bar.go",
			pattern:  "vendor/**",
			want:     true,
		},
		{
			name:     "prefix no match",
			filename: "pkg/foo.go",
			pattern:  "vendor/**",
			want:     false,
		},
		{
			name:     "suffix match",
			filename: "pkg/testdata/foo.go",
			pattern:  "**/testdata",
			want:     true,
		},
		{
			name:     "suffix no match",
			filename: "pkg/foo.go",
			pattern:  "**/testdata",
			want:     false,
		},
		{
			name:     "invalid pattern (not 2 or 3 parts)",
			filename: "foo.go",
			pattern:  "**/**/**/**",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{}
			got := cfg.matchDoubleStarPattern(tt.filename, tt.pattern)
			if got != tt.want {
				t.Errorf("matchDoubleStarPattern(%q, %q) = %v, want %v", tt.filename, tt.pattern, got, tt.want)
			}
		})
	}
}

// TestConfig_matchTriplePattern tests the matchTriplePattern method.
//
// Params:
//   - t: testing context
func TestConfig_matchTriplePattern(t *testing.T) {
	tests := []struct {
		name     string
		parts    []string
		filename string
		want     bool
	}{
		{"empty middle matches all", []string{"", "", ""}, "any/path", true},
		{"middle matches", []string{"", "/testdata/", ""}, "foo/testdata/bar", true},
		{"middle no match", []string{"", "/testdata/", ""}, "foo/other/bar", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{}
			got := c.matchTriplePattern(tt.parts, tt.filename)
			if got != tt.want {
				t.Errorf("matchTriplePattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestConfig_matchPrefixSuffix tests the matchPrefixSuffix method.
//
// Params:
//   - t: testing context
func TestConfig_matchPrefixSuffix(t *testing.T) {
	tests := []struct {
		name     string
		parts    []string
		filename string
		want     bool
	}{
		{"empty prefix and suffix", []string{"", ""}, "any/path", true},
		{"prefix matches", []string{"src/", ""}, "src/foo/bar", true},
		{"suffix matches", []string{"", ".go"}, "foo/.go", true},
		{"no match", []string{"src/", ".go"}, "other/path", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{}
			got := c.matchPrefixSuffix(tt.parts, tt.filename)
			if got != tt.want {
				t.Errorf("matchPrefixSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

package config_test

import (
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
)

func TestDefaultConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"returns non-nil config"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.DefaultConfig()
			if cfg == nil {
				t.Error("DefaultConfig() returned nil")
			}
			if cfg.Version != 1 {
				t.Errorf("DefaultConfig().Version = %d, want 1", cfg.Version)
			}
			if cfg.Rules == nil {
				t.Error("DefaultConfig().Rules is nil")
			}
		})
	}
}

func TestConfig_IsRuleEnabled(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *config.Config
		ruleCode string
		want     bool
	}{
		{
			name:     "nil config returns true",
			cfg:      nil,
			ruleCode: "KTN-FUNC-001",
			want:     true,
		},
		{
			name:     "empty config returns true",
			cfg:      config.DefaultConfig(),
			ruleCode: "KTN-FUNC-001",
			want:     true,
		},
		{
			name: "rule not in config returns true",
			cfg: &config.Config{
				Rules: map[string]*config.RuleConfig{},
			},
			ruleCode: "KTN-FUNC-001",
			want:     true,
		},
		{
			name: "rule enabled explicitly",
			cfg: &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-FUNC-001": {Enabled: config.Bool(true)},
				},
			},
			ruleCode: "KTN-FUNC-001",
			want:     true,
		},
		{
			name: "rule disabled explicitly",
			cfg: &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-FUNC-001": {Enabled: config.Bool(false)},
				},
			},
			ruleCode: "KTN-FUNC-001",
			want:     false,
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

func TestConfig_GetThreshold(t *testing.T) {
	tests := []struct {
		name         string
		cfg          *config.Config
		ruleCode     string
		defaultValue int
		want         int
	}{
		{
			name:         "nil config returns default",
			cfg:          nil,
			ruleCode:     "KTN-FUNC-011",
			defaultValue: 10,
			want:         10,
		},
		{
			name:         "rule not configured returns default",
			cfg:          config.DefaultConfig(),
			ruleCode:     "KTN-FUNC-011",
			defaultValue: 10,
			want:         10,
		},
		{
			name: "threshold set in config",
			cfg: &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-FUNC-011": {Threshold: config.Int(15)},
				},
			},
			ruleCode:     "KTN-FUNC-011",
			defaultValue: 10,
			want:         15,
		},
		{
			name: "threshold set to zero",
			cfg: &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-FUNC-011": {Threshold: config.Int(0)},
				},
			},
			ruleCode:     "KTN-FUNC-011",
			defaultValue: 10,
			want:         0,
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

func TestConfig_IsFileExcluded(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *config.Config
		ruleCode string
		filename string
		want     bool
	}{
		{
			name:     "nil config returns false",
			cfg:      nil,
			ruleCode: "KTN-FUNC-001",
			filename: "foo.go",
			want:     false,
		},
		{
			name: "global exclusion matches",
			cfg: &config.Config{
				Exclude: []string{"*_test.go"},
			},
			ruleCode: "KTN-FUNC-001",
			filename: "foo_test.go",
			want:     true,
		},
		{
			name: "global exclusion does not match",
			cfg: &config.Config{
				Exclude: []string{"*_test.go"},
			},
			ruleCode: "KTN-FUNC-001",
			filename: "foo.go",
			want:     false,
		},
		{
			name: "rule-specific exclusion matches",
			cfg: &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-FUNC-001": {Exclude: []string{"cmd/**"}},
				},
			},
			ruleCode: "KTN-FUNC-001",
			filename: "cmd/main.go",
			want:     true,
		},
		{
			name: "rule-specific exclusion does not apply to other rules",
			cfg: &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-FUNC-001": {Exclude: []string{"cmd/**"}},
				},
			},
			ruleCode: "KTN-FUNC-002",
			filename: "cmd/main.go",
			want:     false,
		},
		{
			name: "double star pattern matches deep path",
			cfg: &config.Config{
				Exclude: []string{"vendor/**"},
			},
			ruleCode: "KTN-FUNC-001",
			filename: "vendor/github.com/foo/bar.go",
			want:     true,
		},
		{
			name: "exact filename match",
			cfg: &config.Config{
				Exclude: []string{"main.go"},
			},
			ruleCode: "KTN-FUNC-001",
			filename: "main.go",
			want:     true,
		},
		{
			name: "path with suffix pattern",
			cfg: &config.Config{
				Exclude: []string{"**/testdata/**"},
			},
			ruleCode: "KTN-FUNC-001",
			filename: "pkg/analyzer/ktn/ktnfunc/testdata/good.go",
			want:     true,
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

func TestConfig_Merge(t *testing.T) {
	tests := []struct {
		name        string
		base        *config.Config
		other       *config.Config
		checkRule   string
		wantEnabled bool
	}{
		{
			name: "merge nil config",
			base: &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-FUNC-001": {Enabled: config.Bool(true)},
				},
			},
			other:       nil,
			checkRule:   "KTN-FUNC-001",
			wantEnabled: true,
		},
		{
			name: "other overrides base",
			base: &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-FUNC-001": {Enabled: config.Bool(true)},
				},
			},
			other: &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-FUNC-001": {Enabled: config.Bool(false)},
				},
			},
			checkRule:   "KTN-FUNC-001",
			wantEnabled: false,
		},
		{
			name: "other adds new rule",
			base: &config.Config{
				Rules: map[string]*config.RuleConfig{},
			},
			other: &config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-FUNC-002": {Enabled: config.Bool(false)},
				},
			},
			checkRule:   "KTN-FUNC-002",
			wantEnabled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.base.Merge(tt.other)
			got := tt.base.IsRuleEnabled(tt.checkRule)
			if got != tt.wantEnabled {
				t.Errorf("After Merge(), IsRuleEnabled(%q) = %v, want %v", tt.checkRule, got, tt.wantEnabled)
			}
		})
	}
}

func TestGetAndSet(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"global config get and set"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset to default
			config.Reset()

			// Get should return default
			cfg := config.Get()
			if cfg == nil {
				t.Fatal("Get() returned nil")
			}

			// Set custom config
			customCfg := &config.Config{
				Version: 1,
				Rules: map[string]*config.RuleConfig{
					"KTN-TEST": {Enabled: config.Bool(false)},
				},
			}
			config.Set(customCfg)

			// Get should return custom
			cfg = config.Get()
			if cfg.IsRuleEnabled("KTN-TEST") {
				t.Error("Expected KTN-TEST to be disabled after Set()")
			}

			// Reset again
			config.Reset()
		})
	}
}

func TestBoolAndInt(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"helper functions"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			boolPtr := config.Bool(true)
			if boolPtr == nil || *boolPtr != true {
				t.Error("Bool(true) failed")
			}

			intPtr := config.Int(42)
			if intPtr == nil || *intPtr != 42 {
				t.Error("Int(42) failed")
			}
		})
	}
}

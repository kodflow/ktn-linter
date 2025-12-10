package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadFromFile_ReadError tests loadFromFile with read errors.
func TestLoadFromFile_ReadError(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "nonexistent file",
			path:    "/nonexistent/file/path.yaml",
			wantErr: true,
		},
		{
			name:    "directory instead of file",
			path:    os.TempDir(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := loadFromFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadFromFile(%q) error = %v, wantErr %v", tt.path, err, tt.wantErr)
			}
		})
	}
}

// TestLoadFromFile_InvalidYAML tests loadFromFile with invalid YAML.
func TestLoadFromFile_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "invalid.yaml")

	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "malformed yaml",
			content: "invalid: yaml: [[[",
			wantErr: true,
		},
		{
			name:    "tab in yaml",
			content: "version:\t1",
			wantErr: false, // YAML actually allows tabs in some contexts
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := os.WriteFile(path, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			_, err := loadFromFile(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadFromFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateConfig_InvalidVersion tests validateConfig with invalid versions.
func TestValidateConfig_InvalidVersion(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name:    "nil config",
			cfg:     nil,
			wantErr: false,
		},
		{
			name: "version 0 is valid",
			cfg: &Config{
				Version: 0,
			},
			wantErr: false,
		},
		{
			name: "version 1 is valid",
			cfg: &Config{
				Version: 1,
			},
			wantErr: false,
		},
		{
			name: "version 2 is invalid",
			cfg: &Config{
				Version: 2,
			},
			wantErr: true,
		},
		{
			name: "version 99 is invalid",
			cfg: &Config{
				Version: 99,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateConfig_NegativeThreshold tests validateConfig with negative thresholds.
func TestValidateConfig_NegativeThreshold(t *testing.T) {
	tests := []struct {
		name      string
		threshold *int
		wantErr   bool
	}{
		{
			name:      "nil threshold",
			threshold: nil,
			wantErr:   false,
		},
		{
			name:      "zero threshold",
			threshold: Int(0),
			wantErr:   false,
		},
		{
			name:      "positive threshold",
			threshold: Int(10),
			wantErr:   false,
		},
		{
			name:      "negative threshold -1",
			threshold: Int(-1),
			wantErr:   true,
		},
		{
			name:      "negative threshold -100",
			threshold: Int(-100),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Version: 1,
				Rules: map[string]*RuleConfig{
					"TEST-RULE": {Threshold: tt.threshold},
				},
			}
			err := validateConfig(cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateConfig_EmptyPattern tests validateConfig with empty patterns.
func TestValidateConfig_EmptyPattern(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "empty global exclusion pattern",
			cfg: &Config{
				Version: 1,
				Exclude: []string{""},
			},
			wantErr: true,
		},
		{
			name: "valid global exclusion pattern",
			cfg: &Config{
				Version: 1,
				Exclude: []string{"*_test.go"},
			},
			wantErr: false,
		},
		{
			name: "empty rule exclusion pattern",
			cfg: &Config{
				Version: 1,
				Rules: map[string]*RuleConfig{
					"TEST-RULE": {Exclude: []string{""}},
				},
			},
			wantErr: true,
		},
		{
			name: "valid rule exclusion pattern",
			cfg: &Config{
				Version: 1,
				Rules: map[string]*RuleConfig{
					"TEST-RULE": {Exclude: []string{"vendor/**"}},
				},
			},
			wantErr: false,
		},
		{
			name: "nil rule config",
			cfg: &Config{
				Version: 1,
				Rules: map[string]*RuleConfig{
					"TEST-RULE": nil,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestLoadFromDefaultLocations_NoConfigFile tests loading with no config file.
func TestLoadFromDefaultLocations_NoConfigFile(t *testing.T) {
	// Create a temp directory with no config files
	tmpDir := t.TempDir()

	// Change to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Load should return default config
	cfg, err := loadFromDefaultLocations()
	if err != nil {
		t.Errorf("loadFromDefaultLocations() unexpected error: %v", err)
	}
	if cfg == nil {
		t.Error("loadFromDefaultLocations() returned nil config")
	}
	if cfg.Version != 1 {
		t.Errorf("Expected default config with version 1, got %d", cfg.Version)
	}
}

// TestLoadFromDefaultLocations_FindsInCurrentDir tests finding config in current directory.
func TestLoadFromDefaultLocations_FindsInCurrentDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create config in current directory
	configPath := filepath.Join(tmpDir, DefaultConfigFileName)
	content := `version: 1
rules:
  CURRENT-DIR-TEST:
    enabled: false
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Change to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Load should find config in current directory
	cfg, err := loadFromDefaultLocations()
	if err != nil {
		t.Fatalf("loadFromDefaultLocations() error = %v", err)
	}

	if cfg.IsRuleEnabled("CURRENT-DIR-TEST") {
		t.Error("Expected CURRENT-DIR-TEST to be disabled")
	}
}

// TestLoadFromDefaultLocations_ErrorInFile tests error handling when loading file.
func TestLoadFromDefaultLocations_ErrorInFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create invalid config file
	configPath := filepath.Join(tmpDir, DefaultConfigFileName)
	content := `version: 99
invalid yaml content
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Change to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Load should return error from invalid file
	_, err = loadFromDefaultLocations()
	if err == nil {
		t.Error("loadFromDefaultLocations() expected error for invalid file")
	}
}

// TestLoadFromDefaultLocations_AlternateFileName tests alternate config filename.
func TestLoadFromDefaultLocations_AlternateFileName(t *testing.T) {
	tmpDir := t.TempDir()

	// Create config with alternate name
	configPath := filepath.Join(tmpDir, AlternateConfigFileName)
	content := `version: 1
rules:
  ALT-TEST:
    enabled: false
`
	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	// Change to temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Load should find alternate file
	cfg, err := loadFromDefaultLocations()
	if err != nil {
		t.Fatalf("loadFromDefaultLocations() error = %v", err)
	}

	if cfg.IsRuleEnabled("ALT-TEST") {
		t.Error("Expected ALT-TEST to be disabled")
	}
}

// TestFileExists tests fileExists function.
func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name  string
		setup func() string
		want  bool
	}{
		{
			name: "file exists",
			setup: func() string {
				path := filepath.Join(tmpDir, "exists.txt")
				os.WriteFile(path, []byte("test"), 0644)
				return path
			},
			want: true,
		},
		{
			name: "file does not exist",
			setup: func() string {
				return filepath.Join(tmpDir, "notexists.txt")
			},
			want: false,
		},
		{
			name: "path is directory",
			setup: func() string {
				dirPath := filepath.Join(tmpDir, "testdir")
				os.MkdirAll(dirPath, 0755)
				return dirPath
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			got := fileExists(path)
			if got != tt.want {
				t.Errorf("fileExists(%q) = %v, want %v", path, got, tt.want)
			}
		})
	}
}

// TestLoadAndSet_Error tests LoadAndSet with error.
func TestLoadAndSet_Error(t *testing.T) {
	err := LoadAndSet("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("LoadAndSet() expected error for nonexistent file")
	}
}

// TestMustLoad_Success tests MustLoad with valid file.
func TestMustLoad_Success(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "config.yaml")
	content := `version: 1
rules:
  MUST-TEST:
    enabled: false
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	cfg := MustLoad(path)
	if cfg == nil {
		t.Error("MustLoad() returned nil")
	}
	if cfg.IsRuleEnabled("MUST-TEST") {
		t.Error("Expected MUST-TEST to be disabled")
	}
}

// TestSaveToFile_MarshalError tests SaveToFile edge cases.
func TestSaveToFile_WriteError(t *testing.T) {
	cfg := DefaultConfig()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "invalid directory",
			path:    "/nonexistent/directory/config.yaml",
			wantErr: true,
		},
		{
			name:    "valid path",
			path:    filepath.Join(t.TempDir(), "valid.yaml"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SaveToFile(cfg, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestLoad_EmptyPathReturnsDefault tests Load with empty path.
func TestLoad_EmptyPathReturnsDefault(t *testing.T) {
	// Change to directory with no config
	tmpDir := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	cfg, err := Load("")
	if err != nil {
		t.Errorf("Load('') unexpected error: %v", err)
	}
	if cfg == nil {
		t.Error("Load('') returned nil")
	}
	if cfg.Version != 1 {
		t.Errorf("Expected default version 1, got %d", cfg.Version)
	}
}

// TestValidateConfig_MultipleErrors tests multiple validation errors.
func TestValidateConfig_MultipleErrors(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "invalid version and negative threshold",
			cfg: &Config{
				Version: 99,
				Rules: map[string]*RuleConfig{
					"TEST-RULE": {Threshold: Int(-5)},
				},
			},
			wantErr: true,
		},
		{
			name: "empty global and rule exclusion patterns",
			cfg: &Config{
				Version: 1,
				Exclude: []string{"valid.go", ""},
				Rules: map[string]*RuleConfig{
					"TEST-RULE": {Exclude: []string{""}},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

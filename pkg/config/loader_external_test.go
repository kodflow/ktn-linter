package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(t *testing.T) string
		cleanupFunc func(path string)
		wantErr     bool
	}{
		{
			name: "load from explicit path",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				path := filepath.Join(tmpDir, "config.yaml")
				content := `version: 1
exclude:
  - "*_test.go"
rules:
  KTN-FUNC-011:
    threshold: 15
`
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return path
			},
			cleanupFunc: func(path string) {},
			wantErr:     false,
		},
		{
			name: "file not found",
			setupFunc: func(t *testing.T) string {
				return "/nonexistent/path/config.yaml"
			},
			cleanupFunc: func(path string) {},
			wantErr:     true,
		},
		{
			name: "invalid yaml",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				path := filepath.Join(tmpDir, "config.yaml")
				content := `invalid: yaml: content:`
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return path
			},
			cleanupFunc: func(path string) {},
			wantErr:     true,
		},
		{
			name: "invalid version",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				path := filepath.Join(tmpDir, "config.yaml")
				content := `version: 99`
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return path
			},
			cleanupFunc: func(path string) {},
			wantErr:     true,
		},
		{
			name: "negative threshold",
			setupFunc: func(t *testing.T) string {
				tmpDir := t.TempDir()
				path := filepath.Join(tmpDir, "config.yaml")
				content := `version: 1
rules:
  KTN-FUNC-011:
    threshold: -5
`
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				return path
			},
			cleanupFunc: func(path string) {},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setupFunc(t)
			defer tt.cleanupFunc(path)

			cfg, err := config.Load(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && cfg == nil {
				t.Error("Load() returned nil config without error")
			}
		})
	}
}

func TestLoad_EmptyPath(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"returns default config when no file found"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Change to temp dir with no config file
			tmpDir := t.TempDir()
			oldDir, _ := os.Getwd()
			os.Chdir(tmpDir)
			defer os.Chdir(oldDir)

			cfg, err := config.Load("")
			if err != nil {
				t.Errorf("Load('') error = %v", err)
				return
			}

			if cfg == nil {
				t.Error("Load('') returned nil")
			}
		})
	}
}

func TestLoad_FindsConfigInParentDir(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"finds config file in parent directory"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create a temp directory structure
			tmpDir := t.TempDir()
			subDir := filepath.Join(tmpDir, "subdir")
			os.MkdirAll(subDir, 0755)

			// Create config in parent dir
			configPath := filepath.Join(tmpDir, ".ktn-linter.yaml")
			content := `version: 1
rules:
  KTN-TEST-RULE:
    enabled: false
`
			os.WriteFile(configPath, []byte(content), 0644)

			// Change to subdir
			oldDir, _ := os.Getwd()
			os.Chdir(subDir)
			defer os.Chdir(oldDir)

			cfg, err := config.Load("")
			if err != nil {
				t.Errorf("Load('') error = %v", err)
				return
			}

			if cfg.IsRuleEnabled("KTN-TEST-RULE") {
				t.Error("Expected KTN-TEST-RULE to be disabled")
			}
		})
	}
}

func TestLoadAndSet(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"load and set successfully", false},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create temp config
			tmpDir := t.TempDir()
			path := filepath.Join(tmpDir, "config.yaml")
			content := `version: 1
rules:
  KTN-LOAD-TEST:
    enabled: false
`
			os.WriteFile(path, []byte(content), 0644)

			// Reset before test
			config.Reset()

			err := config.LoadAndSet(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAndSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check global config was updated
			if config.Get().IsRuleEnabled("KTN-LOAD-TEST") {
				t.Error("Expected KTN-LOAD-TEST to be disabled after LoadAndSet")
			}

			// Reset after test
			config.Reset()
		})
	}
}

func TestSaveToFile(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
	}{
		{
			name: "save valid config",
			cfg: &config.Config{
				Version: 1,
				Exclude: []string{"*_test.go"},
				Rules: map[string]*config.RuleConfig{
					"KTN-FUNC-011": {Threshold: config.Int(15)},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			path := filepath.Join(tmpDir, "saved-config.yaml")

			err := config.SaveToFile(tt.cfg, path)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveToFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify file was created
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Error("SaveToFile() did not create file")
			}

			// Load and verify
			loaded, err := config.Load(path)
			if err != nil {
				t.Errorf("Failed to load saved config: %v", err)
				return
			}

			if loaded.GetThreshold("KTN-FUNC-011", 10) != 15 {
				t.Error("Loaded config has wrong threshold")
			}
		})
	}
}

func TestMustLoad(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		wantPanic bool
	}{
		{
			name:      "panics on invalid path",
			path:      "/nonexistent/config.yaml",
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if (r != nil) != tt.wantPanic {
					t.Errorf("MustLoad() panic = %v, wantPanic %v", r != nil, tt.wantPanic)
				}
			}()

			config.MustLoad(tt.path)
		})
	}
}

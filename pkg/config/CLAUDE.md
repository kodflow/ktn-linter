# pkg/config/ - Configuration Management

## Purpose
Load and manage `.ktn-linter.yaml` configuration with singleton pattern.

## Structure
```
config/
├── config.go         # Config struct + Get() singleton
├── loader.go         # YAML loading + search paths
└── *_test.go         # Tests
```

## Config File Locations (priority order)
1. `--config` flag path
2. `.ktn-linter.yaml` in current directory
3. `.ktn-linter.yaml` in parent directories (up to root)

## Config Structure
```yaml
# .ktn-linter.yaml
thresholds:
  max_function_lines: 35
  max_params: 5
  max_struct_fields: 10

rules:
  KTN-FUNC-001:
    enabled: true
    severity: error
  KTN-VAR-001:
    enabled: false  # Disable specific rule

exclude:
  paths:
    - "vendor/"
    - "generated/"
```

## Usage in Analyzers
```go
cfg := config.Get()
if cfg.IsRuleEnabled("KTN-FUNC-001") {
    // Run rule
}
threshold := cfg.GetThreshold("max_function_lines", 35)
```

## Singleton Pattern
- `Get()` returns cached config (loaded once)
- `Load(path)` explicitly loads from path
- Thread-safe access

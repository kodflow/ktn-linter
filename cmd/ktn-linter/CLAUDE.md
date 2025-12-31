# cmd/ktn-linter/ - Main CLI Application

## Purpose
Entry point for the KTN-Linter command-line tool. Contains the main package
and Cobra command implementations.

## Structure
```
cmd/ktn-linter/
├── main.go                  # Entry point, version injection
├── main_internal_test.go    # Main package tests
└── cmd/
    ├── root.go              # Root command, global flags
    ├── root_external_test.go
    ├── root_internal_test.go
    ├── lint.go              # Main lint orchestration
    ├── lint_internal_test.go
    └── diag.go              # Diagnostic type definitions
```

## main.go
- Package `main` with minimal setup
- Version injected via `-ldflags`:
  ```bash
  go build -ldflags "-X main.version=1.0.0" ./cmd/ktn-linter
  ```

## Build
```bash
# Standard build
make build

# With version
go build -ldflags "-X main.version=$(git describe --tags)" \
    -o builds/ktn-linter ./cmd/ktn-linter
```

## Dependencies
- `github.com/spf13/cobra` for CLI framework
- `pkg/orchestrator` for linting pipeline
- `pkg/formatter` for output formatting
- `pkg/config` for YAML configuration

## Test Conventions
- `*_external_test.go` for black-box tests (package `cmd_test`)
- `*_internal_test.go` for white-box tests (package `cmd`)

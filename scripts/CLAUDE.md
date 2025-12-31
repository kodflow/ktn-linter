# scripts/ - Build & Validation Scripts

## Purpose
Contains shell scripts for validation, coverage reporting, and CI tasks.

## Files
| Script | Purpose |
|--------|---------|
| `validate-testdata.sh` | Validates all testdata files directly |
| `generate-coverage.sh` | Generates COVERAGE.MD report |

## validate-testdata.sh
- Runs ktn-linter directly on each testdata file
- Uses `--only-rule` to isolate each rule
- Verifies: good.go = 0 errors, bad.go = expected errors only
- Checks for redeclarations between good.go and bad.go
- Required because `go list ./...` excludes testdata directories

## generate-coverage.sh
- Runs `go test -coverprofile` on all packages
- Generates COVERAGE.MD with per-package breakdown
- Uses icons: 100%, >90%, 80-90%, <80%, 0%

## Usage
```bash
# Validate all testdata
./scripts/validate-testdata.sh

# Generate coverage report
./scripts/generate-coverage.sh
```

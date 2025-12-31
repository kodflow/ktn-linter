# bin/ - Utility Scripts

## Purpose
Contains wrapper scripts for integrating KTN-Linter with external tools.

## Files
| Script | Purpose |
|--------|---------|
| `golangci-lint-wrapper` | Combines golangci-lint + ktn-linter output |

## golangci-lint-wrapper
- Auto-builds ktn-linter if missing (`/workspace/builds/ktn-linter`)
- Runs real golangci-lint only if linters are enabled in `.golangci.yml`
- Skips testdata directories (use `make validate` instead)
- Combines output from both linters into single stream

## Usage
```bash
# Called automatically by CI or IDE integrations
./bin/golangci-lint-wrapper run ./...
```

## Dependencies
- Real golangci-lint at `/usr/local/bin/golangci-lint`
- KTN-Linter source in `/workspace/cmd/ktn-linter`

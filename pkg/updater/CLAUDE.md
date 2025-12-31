# pkg/updater/ - Self-Update System

## Purpose
Handles automatic updates for ktn-linter binary via GitHub releases.

## Structure
```
updater/
├── updater.go         # Core update logic (check, download, replace)
├── types.go           # Update result types
└── *_test.go          # Tests (external + internal)
```

## Key Functions
| Function | Description |
|----------|-------------|
| `CheckForUpdate()` | Queries GitHub API for latest release |
| `PerformUpdate()` | Downloads and replaces current binary |
| `GetCurrentVersion()` | Returns running binary version |

## Update Flow
1. Fetch latest release from GitHub API
2. Compare with current version (semver)
3. Download platform-specific binary
4. Replace current executable atomically

## Dependencies
- GitHub Releases API
- Platform detection (GOOS/GOARCH)

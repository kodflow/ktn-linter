# cmd/ - CLI Entry Points

## Structure
```
cmd/
└── ktn-linter/
    ├── main.go          # Version injection, entry point
    └── cmd/
        ├── root.go      # Cobra root command, global flags
        ├── lint.go      # Main lint orchestration (~600 lines)
        └── diag.go      # Diagnostic types
```

## Conventions
- Package `main` in `main.go` (version via ldflags)
- Cobra commands in `cmd/` subdirectory
- Tests: `*_internal_test.go` (white-box), `*_external_test.go` (black-box)

## Key Files
- `lint.go`: Core orchestration (loadPackages → selectAnalyzers → runAnalyzers)
- `root.go`: Flags parsing (--category, --only-rule, --verbose, --config)

## Dependencies
- `github.com/spf13/cobra` for CLI
- `golang.org/x/tools/go/packages` for AST loading
- Internal: `pkg/analyzer/ktn`, `pkg/formatter`, `pkg/config`

# pkg/orchestrator/ - Linting Orchestration

## Purpose
Coordinates the complete linting pipeline: loading packages, selecting analyzers,
running analysis, and processing diagnostics.

## Architecture
```
┌─────────────────────────────────────────────────┐
│              Orchestrator                        │
│  ┌──────────┐ ┌──────────┐ ┌──────────────────┐ │
│  │  Loader  │ │ Selector │ │     Runner       │ │
│  └──────────┘ └──────────┘ └──────────────────┘ │
│  ┌──────────────────────────────────────────┐   │
│  │          DiagnosticsProcessor            │   │
│  └──────────────────────────────────────────┘   │
└─────────────────────────────────────────────────┘
```

## Files
| File | Responsibility |
|------|----------------|
| `orchestrator.go` | Main coordinator, `NewOrchestrator()`, `Run()` |
| `loader.go` | `PackageLoader`: loads Go packages |
| `selector.go` | `AnalyzerSelector`: chooses analyzers |
| `runner.go` | `AnalysisRunner`: executes analyzers |
| `diagnostics.go` | `DiagnosticsProcessor`: filters, deduplicates |
| `options.go` | `Options` struct |
| `types.go` | `DiagnosticResult` type |

## Usage
```go
import "github.com/kodflow/ktn-linter/pkg/orchestrator"

orch := orchestrator.NewOrchestrator(os.Stderr, verbose)

// Full pipeline
diags, err := orch.Run(patterns, opts)

// Or step-by-step
pkgs, _ := orch.LoadPackages(patterns)
analyzers, _ := orch.SelectAnalyzers(opts)
rawDiags := orch.RunAnalyzers(pkgs, analyzers)
filtered := orch.FilterDiagnostics(rawDiags)
diags := orch.ExtractDiagnostics(filtered)
```

## Options
```go
type Options struct {
    Verbose    bool   // Enable verbose logging
    Category   string // Filter by category (e.g., "func")
    OnlyRule   string // Run single rule (e.g., "KTN-FUNC-001")
    ConfigPath string // Path to config file
}
```

## Dependencies
- `golang.org/x/tools/go/packages` for package loading
- `golang.org/x/tools/go/analysis` for analyzer framework
- `pkg/analyzer/ktn` for rule registry
- `pkg/config` for configuration

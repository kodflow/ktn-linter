# cmd/ktn-linter/cmd/ - Cobra Commands

## Purpose
Implements the CLI commands using the Cobra framework. Handles argument parsing,
flag processing, and orchestrates the linting pipeline.

## Files
| File | Responsibility |
|------|----------------|
| `root.go` | Root command, global flags, subcommand registration |
| `lint.go` | Main lint command, delegates to orchestrator |
| `rules.go` | Display rules with descriptions and examples |
| `diag.go` | Diagnostic type definitions |

## root.go - Root Command
Global flags available to all subcommands:
- `--verbose`, `-v`: Enable verbose output
- `--config`, `-c`: Path to config file (default: `.ktn-linter.yaml`)
- `--category`: Filter by rule category (func, var, const, etc.)
- `--only-rule`: Run single rule by code (e.g., `KTN-FUNC-001`)

## lint.go - Lint Command
Pipeline execution:
1. Parse target patterns (default: `./...`)
2. Create orchestrator with options
3. Run analysis pipeline
4. Format and display diagnostics

```go
orch := orchestrator.NewOrchestrator(os.Stderr, verbose)
diags, err := orch.Run(patterns, opts)
```

## rules.go - Rules Command
Display all KTN rules for AI-assisted development:
- `--format`: Output format (markdown, json, text)
- `--no-examples`: Skip loading good.go examples

```bash
ktn-linter rules --format=markdown
ktn-linter rules --category=func --format=json
```

## diag.go - Diagnostics
Defines `Diag` struct for lint findings:
- Position (file, line, column)
- Message and rule code
- Severity level

## Test Organization
| Test File | Type | Purpose |
|-----------|------|---------|
| `root_external_test.go` | Black-box | Test CLI interface |
| `root_internal_test.go` | White-box | Test internal helpers |
| `lint_internal_test.go` | White-box | Test lint command logic |
| `rules_external_test.go` | Black-box | Test rules output |
| `rules_internal_test.go` | White-box | Test rules formatting |

## Dependencies
- `pkg/orchestrator`: Linting pipeline coordination
- `pkg/analyzer/ktn`: Rule registry access
- `pkg/formatter`: Output formatting
- `pkg/config`: Configuration loading
- `pkg/messages`: Rule messages
- `pkg/rules`: Rule info extraction

## Adding New Commands
1. Create `newcmd.go` in this directory
2. Define command with `&cobra.Command{}`
3. Add to root command in `root.go`
4. Create tests: `newcmd_external_test.go` and/or `newcmd_internal_test.go`

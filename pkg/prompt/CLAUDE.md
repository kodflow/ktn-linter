# pkg/prompt/ - AI Prompt Generation

## Purpose
Generates structured prompts for AI-assisted code fixes based on linting violations.

## Structure
```
prompt/
├── generator.go       # Main prompt builder
├── formatter.go       # Output formatting for prompts
├── phases.go          # Multi-phase prompt structure
├── phase_group.go     # Phase grouping utilities
├── types.go           # Core data types
├── violation.go       # Violation representation
├── rule_violations.go # Rule-grouped violations
├── prompt_output.go   # Final output structure
└── *_test.go          # Comprehensive tests
```

## Key Types
| Type | Description |
|------|-------------|
| `Generator` | Builds prompts from violations |
| `Phase` | Represents a fix phase |
| `Violation` | Single linting error |
| `RuleViolations` | Violations grouped by rule |

## Prompt Structure
- Phase 1: Understanding (code context)
- Phase 2: Analysis (violation details)
- Phase 3: Fix suggestions (remediation)

## Usage
Called by `ktn-linter prompt` command to generate AI-ready output.

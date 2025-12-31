# pkg/ - Core Library Packages

## Structure
```
pkg/
├── analyzer/         # Static analysis (rules + utilities)
│   ├── ktn/          # KTN rules registry + categories
│   ├── modernize/    # golang.org/x/tools wrapper
│   ├── utils/        # AST helpers (100% coverage)
│   └── shared/       # Classification utilities
├── config/           # YAML configuration (singleton)
├── formatter/        # Output formatting (100% coverage)
├── messages/         # Rule message templates
├── orchestrator/     # Linting pipeline coordination
├── rules/            # Rule info extraction for docs/AI
└── severity/         # Severity levels + colors
```

## Package Roles
| Package | Responsibility |
|---------|----------------|
| `analyzer/ktn` | Rule registry, category aggregation |
| `config` | Load .ktn-linter.yaml, per-rule settings |
| `formatter` | Terminal output (colors, grouping) |
| `messages` | Short + verbose messages per rule |
| `orchestrator` | Load packages, select/run analyzers, collect diagnostics |
| `rules` | Extract rule metadata for `ktn-linter rules` command |
| `severity` | INFO/WARN/ERROR levels |

## No Circular Dependencies
Dependency flow: cmd → orchestrator → analyzer → (utils, shared, config, messages, severity) → formatter

# docs/ - Documentation

## Purpose
Contains rule documentation for all KTN linting rules.

## Structure
```
docs/
└── rules/
    ├── KTN-API-001.md
    ├── KTN-COMMENT-001.md .. KTN-COMMENT-007.md
    ├── KTN-CONST-001.md .. KTN-CONST-003.md
    ├── KTN-FUNC-001.md .. KTN-FUNC-013.md
    ├── KTN-INTERFACE-001.md
    ├── KTN-STRUCT-001.md .. KTN-STRUCT-006.md
    ├── KTN-TEST-001.md .. KTN-TEST-011.md (after refactor)
    └── KTN-VAR-001.md .. KTN-VAR-018.md
```

## Rule Documentation Format
Each `.md` file contains:
- Rule code and name
- Description of what it detects
- Code examples (good vs bad)
- Severity level

## Naming Convention
`KTN-<CATEGORY>-<NNN>.md` where NNN is zero-padded (001, 002, etc.)

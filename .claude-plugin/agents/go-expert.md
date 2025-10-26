# Go Expert Agent

Expert Go 1.25+ strict avec **KTN-Linter**. Code production uniquement.

## Règle Absolue

**0 issues = 0 issues**. Même pas INFO. Corriger TOUT avant de continuer.

## Workflow KTN-Linter

Chaque modification de .go :
```bash
# 1. Build le linter si nécessaire
make build

# 2. Lint avec KTN-Linter (94 règles strictes)
./builds/ktn-linter lint ./...

# 3. Vérifier les tests
make test

# Si > 0 issues → STOP et corriger
# Si tests FAIL → STOP et corriger
```

Ne JAMAIS passer à la suite tant que le linter affiche autre chose que "✅ No issues found".

## Installation Auto KTN-Linter

Premier usage sur un projet :
```bash
# Clone et build automatique
git clone https://github.com/kodflow/ktn-linter
cd ktn-linter
make build
# Binaire disponible: ./builds/ktn-linter
```

Ou utiliser le Makefile du projet :
```bash
make build  # Compile ktn-linter
make lint   # Exécute ktn-linter lint ./...
make test   # Lance les tests avec coverage
```

## Patterns (Quand Utiliser)

**Trop de paramètres (>3)** → Functional Options
```go
type Option func(*Config)
func NewServer(opts ...Option) *Server
```

**Construction complexe** → Builder
```go
NewBuilder().Field1(x).Field2(y).Build()
```

**Algo interchangeable** → Strategy
```go
type Compressor interface { Compress([]byte) []byte }
```

**I/O parallèle** → Worker Pool
```go
tasks := make(chan Task)
for i := 0; i < 10; i++ { go worker(tasks) }
```

**Traitement pipeline** → Pipeline
```go
stage1 → stage2 → stage3 (channels)
```

## Conventions Go 1.25

```go
// Constantes
const MAX_RETRY int = 3  // ALL_CAPS + type explicite

// Variables
var userName string      // camelCase

// Erreurs
return fmt.Errorf("...: %w", err)  // wrap avec %w

// Context
func Do(ctx context.Context, ...) error  // toujours 1er param

// Range integers (Go 1.22+)
for i := range 10 { }
```

## Anti-Patterns

❌ Variables globales mutables → DI
❌ Error ignoré (`_`) → Toujours check
❌ Context dans struct → Param fonction
❌ Naked returns → Return explicite
❌ `fmt.Println` → Logger

## Checklist

- [ ] `make lint` → 0 issues (pas même INFO)
- [ ] `make test` → 100% PASS
- [ ] Coverage ≥ 90%
- [ ] Godoc complet
- [ ] Erreurs wrappées

## Réflexe

Après modification → `make lint && make test`
Si échec → Corriger IMMÉDIATEMENT

Pas d'exception. Pas de "je corrige plus tard". Jamais.

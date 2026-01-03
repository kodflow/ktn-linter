# KTN-Linter

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.25+-blue.svg)](https://golang.org)

Linter Go strict pour l'application des bonnes pratiques et règles de style.

**Règle stricte** : 0 issues = 0 issues (même INFO). STOP et corriger immédiatement.

## Installation

### Installation Universelle (Recommandée)

Pour installer ktn-linter sur **n'importe quel projet Go** :

```bash
curl -sSL https://raw.githubusercontent.com/kodflow/ktn-linter/main/install.sh | bash
```

Ou téléchargez et exécutez le script :

```bash
wget https://raw.githubusercontent.com/kodflow/ktn-linter/main/install.sh
chmod +x install.sh
./install.sh
```

Le script :
- ✅ Télécharge le binaire depuis GitHub releases (linux/darwin, amd64/arm64)
- ✅ Installe dans `/usr/local/bin` ou `~/.local/bin`
- ✅ Configure optionnellement golangci-lint
- ✅ Crée un Makefile avec targets ktn-linter

### Installation depuis les sources

```bash
git clone https://github.com/kodflow/ktn-linter
cd ktn-linter
make build      # Compile le binaire dans builds/
```

## Utilisation sur n'importe quel projet

Une fois installé (via `install.sh`), utilisez ktn-linter sur n'importe quel projet Go :

```bash
# Dans votre projet Go
ktn-linter lint ./...                # Lint tout le projet
ktn-linter lint --help               # Affiche l'aide
ktn-linter lint --simple ./pkg/...   # Format simplifié sur pkg/
ktn-linter lint --fix ./...          # Applique automatiquement les fixes modernize
ktn-linter lint --config .ktn-linter.yaml ./...  # Utilise un fichier de config
```

## Configuration (v1.4.0+)

KTN-Linter peut être configuré via un fichier `.ktn-linter.yaml` :

```yaml
version: 1

# Exclusions globales (toutes les règles)
exclude:
  - "**/testdata/**"
  - "**/*_generated.go"
  - "vendor/**"

# Configuration par règle
rules:
  KTN-FUNC-005:
    enabled: true
    threshold: 50          # Lignes max (défaut: 35)
    exclude:
      - "cmd/**"           # Exclure pour cette règle

  KTN-FUNC-011:
    threshold: 15          # Complexité cyclomatique max (défaut: 10)

  KTN-COMMENT-001:
    enabled: false         # Désactiver la règle

  KTN-VAR-009:
    threshold: 100         # Taille struct pour pointeur (défaut: 64)
```

**Règles avec seuils configurables** :
| Règle | Paramètre | Défaut |
|-------|-----------|--------|
| KTN-COMMENT-001 | maxCommentLength | 80 |
| KTN-COMMENT-002 | minPackageCommentLength | 3 |
| KTN-COMMENT-005 | minStructDocLines | 2 |
| KTN-FUNC-005 | maxFunctionLength | 35 |
| KTN-FUNC-006 | maxParameters | 5 |
| KTN-FUNC-010 | maxReturnValues | 3 |
| KTN-FUNC-011 | maxCyclomaticComplexity | 10 |
| KTN-FUNC-012 | maxNestedDepth | 4 |
| KTN-VAR-009 | maxScopeLines | 50 |
| KTN-VAR-012 | maxLineLength | 120 |
| KTN-VAR-016 | maxDeclarations | 10 |

**Recherche du fichier config** :
1. Chemin spécifié avec `--config`
2. `.ktn-linter.yaml` dans le répertoire courant
3. `.ktn-linter.yml` dans le répertoire courant
4. Remonte récursivement dans les répertoires parents

**Flag --fix (v1.3.0+)** :

Applique automatiquement les fixes suggérés par les analyseurs modernize SÛRS :
- ✅ `interface{}` → `any` (Go 1.18+) - Seul analyseur sûr actuellement
- ⚠️ Fixes complexes (slices.Contains, CutSuffix, etc.) : utiliser `go install golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest && modernize -fix ./...`

Le flag `--fix` n'applique que les transformations simples qui ne nécessitent pas d'ajout d'imports, pour éviter de corrompre le code.

**Intégration avec golangci-lint** (optionnel) :

Le script `install.sh` propose de configurer automatiquement `.golangci.yml` pour intégrer ktn-linter comme linter custom.

```bash
# Après installation
golangci-lint run ./...   # Exécute golangci-lint + ktn-linter
```

## Utilisation (développement du linter)

```bash
make test      # Tests + couverture (génère COVERAGE.MD)
make coverage  # Génère uniquement le rapport COVERAGE.MD
make lint      # Lance le linter KTN sur le code de production
make validate  # Valide que tous les testdata good.go/bad.go sont corrects
make build     # Compile le binaire ktn-linter dans builds/
make install   # Compile et installe ktn-linter dans /usr/local/bin
make fmt       # Formate le code Go avec go fmt sur tout le projet
make help      # Aide
```

**Validation testdata** : `make validate` vérifie automatiquement que :
- ✅ Tous les **good.go** : 0 erreur (100% conformes)
- ✅ Tous les **bad.go** : UNIQUEMENT les erreurs de leur règle spécifique
  - Ex: `func001/bad.go` → **seulement** KTN-FUNC-001 (pas de KTN-CONST-001, etc.)
- ✅ Aucune redeclaration entre good.go et bad.go

Voir [COVERAGE.MD](COVERAGE.MD) pour le rapport détaillé de couverture.

### Intégration VSCode

**Linting automatique** : L'extension Go lance automatiquement le linter à la sauvegarde (`Ctrl+S`).

**Voir les erreurs dans les fichiers testdata** :
1. Ouvrir un fichier testdata (ex: `pkg/analyzer/ktn/const/testdata/src/const001/const001.go`)
2. Sauvegarder (`Ctrl+S`) → Les erreurs apparaissent immédiatement
3. Ouvrir l'onglet Problèmes (`Ctrl+Shift+M`) → 50 erreurs détectées

**Fonctionnalités** :
- ✅ Linting automatique (production + testdata)
- ✅ Format simple pour VSCode (`file:line:col: message (CODE)`)
- ✅ Erreurs visibles dans l'éditeur et l'onglet Problèmes
- ✅ Build automatique du binaire à chaque sauvegarde

**Commandes** :
```bash
make lint           # Lint production seulement (exclut testdata)
make lint-testdata  # Vérifie détection sur testdata (784 erreurs)
```

**Configuration** : `.vscode/settings.json`, `.vscode/tasks.json`, `.vscode/keybindings.json`
**Wrapper** : `bin/golangci-lint-wrapper` (format simple, inclut testdata)

## Règles Implémentées (ordonnées par criticité)

### Commentaires et Documentation (7 règles) - INFO/WARNING
| Code | Sévérité | Description |
|------|----------|-------------|
| [KTN-COMMENT-001](docs/rules/KTN-COMMENT-001.md) | INFO | Commentaires inline trop longs (>80 caractères) |
| [KTN-COMMENT-002](docs/rules/KTN-COMMENT-002.md) | WARNING | Commentaire descriptif avant `package` |
| [KTN-COMMENT-003](docs/rules/KTN-COMMENT-003.md) | WARNING | Commentaire obligatoire pour constantes |
| [KTN-COMMENT-004](docs/rules/KTN-COMMENT-004.md) | WARNING | Commentaire obligatoire pour var package |
| [KTN-COMMENT-005](docs/rules/KTN-COMMENT-005.md) | WARNING | Documentation obligatoire pour struct (≥2 lignes) |
| [KTN-COMMENT-006](docs/rules/KTN-COMMENT-006.md) | WARNING | Documentation fonction (Params/Returns) |
| [KTN-COMMENT-007](docs/rules/KTN-COMMENT-007.md) | WARNING | Commentaires sur branches/returns/logique |

### Constantes (6 règles) - ERROR/WARNING/INFO
| Code | Sévérité | Description |
|------|----------|-------------|
| [KTN-CONST-001](docs/rules/KTN-CONST-001.md) | ERROR | Type explicite obligatoire |
| [KTN-CONST-002](docs/rules/KTN-CONST-002.md) | INFO | Groupement et placement avant var |
| [KTN-CONST-003](docs/rules/KTN-CONST-003.md) | INFO | Nommage CamelCase (pas d'underscores) |
| KTN-CONST-004 | WARNING | Constantes non utilisées |
| KTN-CONST-005 | INFO | Constantes dupliquées |
| KTN-CONST-006 | INFO | Constantes magiques (préférer nommées) |

### Variables (36 règles) - ERROR/WARNING/INFO
| Code | Sévérité | Description |
|------|----------|-------------|
| [KTN-VAR-001](docs/rules/KTN-VAR-001.md) | WARNING | Type explicite obligatoire pour var package |
| [KTN-VAR-002](docs/rules/KTN-VAR-002.md) | WARNING | Déclarations ordonnées (const avant var) |
| [KTN-VAR-003](docs/rules/KTN-VAR-003.md) | ERROR | Nommage camelCase obligatoire |
| [KTN-VAR-004](docs/rules/KTN-VAR-004.md) | WARNING | Longueur min variable (scope-aware) |
| [KTN-VAR-005](docs/rules/KTN-VAR-005.md) | WARNING | Longueur max 30 caractères |
| [KTN-VAR-006](docs/rules/KTN-VAR-006.md) | ERROR | Détection shadowing variables |
| [KTN-VAR-007](docs/rules/KTN-VAR-007.md) | INFO | := vs var (zero-value aware) |
| [KTN-VAR-008](docs/rules/KTN-VAR-008.md) | INFO | Préallocation slices avec capacité connue |
| [KTN-VAR-009](docs/rules/KTN-VAR-009.md) | INFO | Éviter make([]T, length) avec append |
| [KTN-VAR-010](docs/rules/KTN-VAR-010.md) | INFO | Préallocation bytes.Buffer avec Grow |
| [KTN-VAR-011](docs/rules/KTN-VAR-011.md) | INFO | Utiliser strings.Builder pour concaténations |
| [KTN-VAR-012](docs/rules/KTN-VAR-012.md) | WARNING | Éviter allocations dans boucles chaudes |
| [KTN-VAR-013](docs/rules/KTN-VAR-013.md) | INFO | Pointeurs pour structs >64 bytes en paramètre |
| [KTN-VAR-014](docs/rules/KTN-VAR-014.md) | INFO | sync.Pool pour buffers répétés |
| [KTN-VAR-015](docs/rules/KTN-VAR-015.md) | INFO | Éviter string() conversions répétées |
| [KTN-VAR-016](docs/rules/KTN-VAR-016.md) | INFO | Groupement dans un seul bloc var() |
| [KTN-VAR-017](docs/rules/KTN-VAR-017.md) | INFO | Préallocation maps avec capacité connue |
| [KTN-VAR-018](docs/rules/KTN-VAR-018.md) | INFO | Utiliser [N]T au lieu de make([]T, N) ≤64 bytes |
| KTN-VAR-019 | ERROR | Copies de mutex (sync.Mutex, sync.RWMutex) |
| KTN-VAR-020 | INFO | Préférer nil slice à empty slice |
| KTN-VAR-021 | WARNING | Consistance receiver (pointer vs value) |
| KTN-VAR-022 | WARNING | Éviter pointeur vers interface |
| KTN-VAR-023 | WARNING | crypto/rand pour données sensibles |
| KTN-VAR-024 | INFO | any vs interface{} (Go 1.18+) |
| KTN-VAR-025 | INFO | Utiliser clear() built-in (Go 1.21+) |
| KTN-VAR-026 | INFO | Utiliser min()/max() built-in (Go 1.21+) |
| KTN-VAR-027 | INFO | range over integer (Go 1.22+) |
| KTN-VAR-028 | INFO | Loop var copy obsolète (Go 1.22+) |
| KTN-VAR-029 | INFO | slices.Grow au lieu de make+copy (Go 1.21+) |
| KTN-VAR-030 | INFO | slices.Clone au lieu de make+copy (Go 1.21+) |
| KTN-VAR-031 | INFO | maps.Clone au lieu de boucle manuelle (Go 1.21+) |
| KTN-VAR-033 | INFO | cmp.Or pour valeurs par défaut (Go 1.22+) |
| KTN-VAR-034 | INFO | WaitGroup.Go (Go 1.25+) |
| KTN-VAR-035 | INFO | slices.Contains au lieu de boucle (Go 1.21+) |
| KTN-VAR-036 | INFO | slices.Index au lieu de boucle (Go 1.21+) |
| KTN-VAR-037 | INFO | maps.Keys/Values iterateurs (Go 1.23+) |

### Fonctions (13 règles) - ERROR/WARNING/INFO
| Code | Sévérité | Description |
|------|----------|-------------|
| [KTN-FUNC-001](docs/rules/KTN-FUNC-001.md) | ERROR | Erreur toujours en dernière position retour |
| [KTN-FUNC-002](docs/rules/KTN-FUNC-002.md) | ERROR | Context toujours en premier paramètre |
| [KTN-FUNC-003](docs/rules/KTN-FUNC-003.md) | ERROR | Éviter else après return/continue/break |
| [KTN-FUNC-004](docs/rules/KTN-FUNC-004.md) | ERROR | Fonctions privées non utilisées (code mort) |
| [KTN-FUNC-005](docs/rules/KTN-FUNC-005.md) | WARNING | Longueur max 35 lignes de code pur |
| [KTN-FUNC-006](docs/rules/KTN-FUNC-006.md) | WARNING | Max 5 paramètres par fonction |
| [KTN-FUNC-007](docs/rules/KTN-FUNC-007.md) | WARNING | Pas de side effects dans les getters |
| [KTN-FUNC-008](docs/rules/KTN-FUNC-008.md) | WARNING | Paramètres non utilisés préfixés par _ |
| [KTN-FUNC-009](docs/rules/KTN-FUNC-009.md) | INFO | Pas de magic numbers (constantes nommées) |
| [KTN-FUNC-010](docs/rules/KTN-FUNC-010.md) | INFO | Pas de naked returns (sauf <5 lignes) |
| [KTN-FUNC-011](docs/rules/KTN-FUNC-011.md) | INFO | Complexité cyclomatique max 10 |
| [KTN-FUNC-012](docs/rules/KTN-FUNC-012.md) | INFO | Named returns pour >3 valeurs de retour |
| KTN-FUNC-013 | WARNING | Préférer slice/map vide à nil |

### Structures (6 règles) - WARNING/INFO
| Code | Sévérité | Description |
|------|----------|-------------|
| [KTN-STRUCT-001](docs/rules/KTN-STRUCT-001.md) | INFO | Convention getters/setters: `Field()` et `SetField()` |
| [KTN-STRUCT-002](docs/rules/KTN-STRUCT-002.md) | WARNING | Constructeur NewX() requis (suffixes autorisés: NewXxxWithOption) |
| [KTN-STRUCT-003](docs/rules/KTN-STRUCT-003.md) | WARNING | Pas de préfixe Get pour getters |
| [KTN-STRUCT-004](docs/rules/KTN-STRUCT-004.md) | INFO | Un fichier Go par struct (DTOs peuvent être groupés) |
| [KTN-STRUCT-005](docs/rules/KTN-STRUCT-005.md) | INFO | Ordre des champs (exportés avant privés) |
| [KTN-STRUCT-006](docs/rules/KTN-STRUCT-006.md) | INFO | Pas de tags de sérialisation sur champs privés |

**Convention Getters/Setters (STRUCT-001)**:
- Getters/setters sont **OPTIONNELS**
- Si présents: `x.Value()` pour get, `x.SetValue(v)` pour set
- Si getter existe mais nom ≠ champ (ex: `Value()` retourne `foo`), suggérer renommage vers `Foo()`

### Tests (13 règles) - ERROR/WARNING/INFO
| Code | Sévérité | Description |
|------|----------|-------------|
| [KTN-TEST-001](docs/rules/KTN-TEST-001.md) | ERROR | Fichiers test doivent finir par _internal/_external_test.go |
| [KTN-TEST-002](docs/rules/KTN-TEST-002.md) | WARNING | Package xxx_test obligatoire (désactivée) |
| [KTN-TEST-003](docs/rules/KTN-TEST-003.md) | WARNING | Fichier test sans fichier source correspondant |
| [KTN-TEST-004](docs/rules/KTN-TEST-004.md) | WARNING | Fonctions publiques sans tests |
| [KTN-TEST-005](docs/rules/KTN-TEST-005.md) | WARNING | Tests sans table-driven pattern |
| [KTN-TEST-006](docs/rules/KTN-TEST-006.md) | WARNING | Pattern 1:1 fichiers test/source |
| [KTN-TEST-007](docs/rules/KTN-TEST-007.md) | WARNING | Interdiction t.Skip() |
| [KTN-TEST-008](docs/rules/KTN-TEST-008.md) | WARNING | Règle 1:2 (_internal_test.go ET _external_test.go) |
| [KTN-TEST-009](docs/rules/KTN-TEST-009.md) | WARNING | Tests publics dans _external_test.go uniquement |
| [KTN-TEST-010](docs/rules/KTN-TEST-010.md) | WARNING | Tests privés dans _internal_test.go uniquement |
| [KTN-TEST-011](docs/rules/KTN-TEST-011.md) | WARNING | Convention package (white-box/black-box) |
| [KTN-TEST-012](docs/rules/KTN-TEST-012.md) | WARNING | Tests doivent contenir des assertions |
| [KTN-TEST-013](docs/rules/KTN-TEST-013.md) | INFO | Coverage cas d'erreur |

### Interfaces (1 règle) - WARNING
| Code | Sévérité | Description |
|------|----------|-------------|
| [KTN-INTERFACE-001](docs/rules/KTN-INTERFACE-001.md) | WARNING | Interface non utilisée |

### API (1 règle) - WARNING
| Code | Sévérité | Description |
|------|----------|-------------|
| [KTN-API-001](docs/rules/KTN-API-001.md) | WARNING | Interfaces minimales côté consumer pour dépendances externes |

### Génériques (5 règles) - ERROR/WARNING/INFO (Go 1.18+)
| Code | Sévérité | Description |
|------|----------|-------------|
| KTN-GENERIC-001 | ERROR | Contrainte comparable requise pour == et != |
| KTN-GENERIC-002 | WARNING | Génériques inutiles sur types interface |
| KTN-GENERIC-003 | WARNING | golang.org/x/exp/constraints déprécié → cmp |
| KTN-GENERIC-005 | WARNING | Type params ne doivent pas shadower identifiants prédéclarés |
| KTN-GENERIC-006 | ERROR | Contrainte cmp.Ordered requise pour <, >, +, -, *, /, % |

### Modernize (17 règles actives / 18 totales) ✅ golang.org/x/tools

Suite officielle d'analyseurs Go pour moderniser le code avec les dernières fonctionnalités du langage et de la stdlib:

**Go 1.18+**
- **any**: `interface{}` → `any`

**Go 1.21+**
- **minmax**: `if a > b { return a }` → `max(a, b)`
- **slicescontains**: Loop manuel → `slices.Contains()`
- **slicessort**: `sort.Slice()` → `slices.Sort()`
- **slicesdelete**: `append(a[:i], a[i+1:]...)` → `slices.Delete()`

**Go 1.22+**
- **rangeint**: `for i := 0; i < n; i++` → `for range n`
- **forvar**: Supprime `x := x` inutiles dans loops
- **reflecttypefor**: `reflect.TypeOf(T{})` → `reflect.TypeFor[T]()`

**Go 1.23+**
- **mapsloop**: Loop manuel → `maps.Keys/Values()`
- **stditerators**: Modernise vers iterateurs stdlib
- **stringsseq**: Modernise manipulation strings

**Go 1.24+**
- **bloop**: `for b.N` → `b.Loop()`
- **testingcontext**: Context manuel → `t.Context()`

**Optimisations générales**
- **fmtappendf**: `append(x, fmt.Sprintf(...))` → `fmt.Appendf()`
- **stringsbuilder**: Concaténation `+=` → `strings.Builder`
- **stringscutprefix**: `HasPrefix+TrimPrefix` → `CutPrefix()`
- **omitzero**: Supprime valeurs zéro redondantes
- **waitgroup**: Pattern manuel → `wg.Go()`

**Analyseurs désactivés** (bugs connus ou instabilité):
- ~~**newexpr**~~: `&T{}` → `new(T)` (désactivé: panic dans certains cas)

**Mise à jour**: `go get -u golang.org/x/tools/go/analysis/passes/modernize@latest && go mod tidy`

## Statistiques

- **Couverture globale**: 96.3%
- **Packages 100%**: utils, formatter, ktn, ktnconst, ktngeneric, ktninterface, modernize, severity
- **Go version**: 1.25+
- **Total règles KTN**: 80 (7 comment + 6 const + 36 var + 13 func + 6 struct + 5 generic + 11 test + 1 interface + 1 api)
- **Total modernize**: 17 analyseurs actifs / 18 totaux
- **Rapport détaillé**: Voir [COVERAGE.MD](COVERAGE.MD)

## Structure

```
/workspace/
├── cmd/ktn-linter/     # Binaire
├── pkg/analyzer/       # Règles d'analyse
└── pkg/formatter/      # Formatage sortie
```

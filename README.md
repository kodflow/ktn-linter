# KTN-Linter

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
```

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
| KTN-COMMENT-001 | INFO | Commentaires inline trop longs (>80 caractères) |
| KTN-COMMENT-002 | WARNING | Commentaire descriptif avant `package` |
| KTN-COMMENT-003 | WARNING | Commentaire obligatoire pour constantes |
| KTN-COMMENT-004 | WARNING | Commentaire obligatoire pour var package |
| KTN-COMMENT-005 | WARNING | Documentation obligatoire pour struct (≥2 lignes) |
| KTN-COMMENT-006 | WARNING | Documentation fonction (Params/Returns) |
| KTN-COMMENT-007 | WARNING | Commentaires sur branches/returns/logique |

### Constantes (3 règles) - WARNING/INFO
| Code | Sévérité | Description |
|------|----------|-------------|
| KTN-CONST-001 | WARNING | Type explicite obligatoire |
| KTN-CONST-002 | INFO | Groupement et placement avant var |
| KTN-CONST-003 | INFO | Nommage SCREAMING_SNAKE_CASE |

### Variables (17 règles) - ERROR/WARNING/INFO
| Code | Sévérité | Description |
|------|----------|-------------|
| KTN-VAR-001 | ERROR | Variables package en camelCase (pas SCREAMING_SNAKE) |
| KTN-VAR-002 | WARNING | Type explicite obligatoire |
| KTN-VAR-003 | WARNING | Utiliser := pour variables locales |
| KTN-VAR-004 | WARNING | Préallocation slices avec capacité connue |
| KTN-VAR-005 | WARNING | Éviter make([]T, length) avec append |
| KTN-VAR-006 | WARNING | Préallocation bytes.Buffer/strings.Builder avec Grow |
| KTN-VAR-007 | WARNING | Utiliser strings.Builder pour >2 concaténations |
| KTN-VAR-008 | WARNING | Éviter allocations dans boucles chaudes |
| KTN-VAR-009 | WARNING | Pointeurs pour structs >64 bytes |
| KTN-VAR-010 | WARNING | sync.Pool pour buffers répétés |
| KTN-VAR-011 | WARNING | Shadowing de variables |
| KTN-VAR-012 | WARNING | Conversions string() répétées |
| KTN-VAR-013 | INFO | Groupement dans un seul bloc var() |
| KTN-VAR-014 | INFO | Variables après constantes (ordre déclarations) |
| KTN-VAR-015 | INFO | Préallocation maps avec capacité connue |
| KTN-VAR-016 | INFO | Utiliser [N]T au lieu de make([]T, N) |
| KTN-VAR-017 | INFO | Copies de mutex (sync.Mutex, sync.RWMutex) |

### Fonctions (12 règles) - ERROR/WARNING/INFO
| Code | Sévérité | Description |
|------|----------|-------------|
| KTN-FUNC-001 | ERROR | Erreur toujours en dernière position retour |
| KTN-FUNC-002 | ERROR | Context toujours en premier paramètre |
| KTN-FUNC-003 | ERROR | Éviter else après return/continue/break |
| KTN-FUNC-004 | ERROR | Fonctions privées non utilisées (code mort) |
| KTN-FUNC-005 | WARNING | Longueur max 35 lignes de code pur |
| KTN-FUNC-006 | WARNING | Max 5 paramètres par fonction |
| KTN-FUNC-007 | WARNING | Pas de side effects dans les getters |
| KTN-FUNC-008 | WARNING | Paramètres non utilisés préfixés par _ |
| KTN-FUNC-009 | INFO | Pas de magic numbers (constantes nommées) |
| KTN-FUNC-010 | INFO | Pas de naked returns (sauf <5 lignes) |
| KTN-FUNC-011 | INFO | Complexité cyclomatique max 10 |
| KTN-FUNC-012 | INFO | Named returns pour >3 valeurs de retour |

### Structures (5 règles) - WARNING/INFO
| Code | Sévérité | Description |
|------|----------|-------------|
| KTN-STRUCT-001 | WARNING | Interface obligatoire (100% méthodes publiques) |
| KTN-STRUCT-002 | WARNING | Constructeur NewX() requis |
| KTN-STRUCT-003 | WARNING | Pas de préfixe Get pour getters |
| KTN-STRUCT-004 | INFO | Un fichier Go par struct |
| KTN-STRUCT-005 | INFO | Ordre des champs (exportés avant privés) |

### Tests (13 règles) - ERROR/WARNING/INFO
| Code | Sévérité | Description |
|------|----------|-------------|
| KTN-TEST-001 | ERROR | Fichiers test doivent finir par _internal/_external_test.go |
| KTN-TEST-002 | WARNING | Package xxx_test obligatoire (désactivée) |
| KTN-TEST-003 | WARNING | Fichier test sans fichier source correspondant |
| KTN-TEST-004 | WARNING | Fonctions publiques sans tests |
| KTN-TEST-005 | WARNING | Tests sans table-driven pattern |
| KTN-TEST-006 | WARNING | Pattern 1:1 fichiers test/source |
| KTN-TEST-007 | WARNING | Interdiction t.Skip() |
| KTN-TEST-008 | WARNING | Règle 1:2 (_internal_test.go ET _external_test.go) |
| KTN-TEST-009 | WARNING | Tests publics dans _external_test.go uniquement |
| KTN-TEST-010 | WARNING | Tests privés dans _internal_test.go uniquement |
| KTN-TEST-011 | WARNING | Convention package (white-box/black-box) |
| KTN-TEST-012 | WARNING | Tests doivent contenir des assertions |
| KTN-TEST-013 | INFO | Coverage cas d'erreur |

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

- **Couverture globale**: 91.6% 🟡
- **Packages 100%**: utils, formatter 🟢
- **Go version**: 1.25
- **Total règles KTN**: 57 (7 comment + 3 const + 17 var + 12 func + 5 struct + 13 test)
- **Total modernize**: 17 analyseurs actifs / 18 totaux
- **Rapport détaillé**: Voir [COVERAGE.MD](COVERAGE.MD)

## Structure

```
/workspace/
├── cmd/ktn-linter/     # Binaire
├── pkg/analyzer/       # Règles d'analyse
└── pkg/formatter/      # Formatage sortie
```

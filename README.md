# KTN-Linter

Linter Go pour l'application des bonnes pratiques.

## Installation

```bash
go mod download
```

## Utilisation

```bash
make test      # Tests + couverture (génère COVERAGE.MD)
make coverage  # Génère uniquement le rapport COVERAGE.MD
make lint      # Lance le linter KTN sur le code de production
make validate  # Valide que tous les testdata good.go/bad.go sont corrects
make build     # Compile le binaire ktn-linter dans builds/
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

## Règles Implémentées

### Constantes (4 règles) ✅ 100%

- **KTN-CONST-001**: Type explicite obligatoire
- **KTN-CONST-002**: Groupement et placement avant var
- **KTN-CONST-003**: Nommage SCREAMING_SNAKE_CASE
- **KTN-CONST-004**: Commentaire obligatoire

### Variables (6 règles) ✅ 100%

- **KTN-VAR-001**: Type explicite obligatoire
- **KTN-VAR-002**: Groupement dans un seul bloc var ()
- **KTN-VAR-003**: Nommage camelCase/PascalCase (pas SCREAMING_SNAKE_CASE)
- **KTN-VAR-004**: Commentaire obligatoire
- **KTN-VAR-005**: Pas d'initialisation multiple sur une ligne
- **KTN-VAR-006**: Variables déclarées après les constantes (ordre imports → const → var → types → fonctions)

### Fonctions (12 règles) ✅ 100%

- **KTN-FUNC-001**: Longueur max 35 lignes de code pur
- **KTN-FUNC-002**: Max 5 paramètres par fonction
- **KTN-FUNC-003**: Pas de magic numbers (constantes nommées)
- **KTN-FUNC-004**: Pas de naked returns (sauf <5 lignes)
- **KTN-FUNC-005**: Complexité cyclomatique max 10
- **KTN-FUNC-006**: Erreur toujours en dernière position
- **KTN-FUNC-007**: Documentation stricte (Params/Returns)
- **KTN-FUNC-008**: Context toujours en premier paramètre
- **KTN-FUNC-009**: Pas de side effects dans les getters
- **KTN-FUNC-010**: Named returns pour >3 valeurs de retour
- **KTN-FUNC-011**: Commentaires sur branches/returns/logique
- **KTN-FUNC-012**: Éviter else après return/continue/break

### Structures (6 règles) ✅ 100%

- **KTN-STRUCT-001**: Un fichier Go par struct (évite fichiers de 10000 lignes)
- **KTN-STRUCT-002**: Interface obligatoire reprenant 100% des méthodes publiques de chaque struct
- **KTN-STRUCT-003**: Ordre des champs (exportés avant privés)
- **KTN-STRUCT-004**: Documentation obligatoire pour structs exportées (≥2 lignes)
- **KTN-STRUCT-005**: Constructeur NewX() requis pour structs avec méthodes
- **KTN-STRUCT-006**: Champs privés + getters pour structs avec méthodes (>3 champs)

## Statistiques

- **Couverture globale**: 95.6% 🟡
- **Packages 100%**: utils, formatter, testhelper 🟢
- **Package const**: 96.6% 🟡
- **Package func**: 94.7% 🟡
- **Package var**: 89.8% 🔴
- **Go version**: 1.25
- **Total règles**: 28 (4 const + 6 var + 12 func + 6 struct)
- **Rapport détaillé**: Voir [COVERAGE.MD](COVERAGE.MD) pour le détail des fonctions < 100%

## Structure

```
/workspace/
├── cmd/ktn-linter/     # Binaire
├── pkg/analyzer/       # Règles d'analyse
└── pkg/formatter/      # Formatage sortie
```

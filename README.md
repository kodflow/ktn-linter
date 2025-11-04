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
```

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

### Retours (1 règle) ✅ 100%

- **KTN-RETURN-002**: Préférer slice/map vide à nil pour éviter nil pointer dereference

### Interfaces (1 règle) ✅ 100%

- **KTN-INTERFACE-001**: Interface déclarée mais jamais utilisée (code mort)

### Commentaires (1 règle) ✅ 100%

- **KTN-COMMENT-002**: Commentaires inline trop verbeux (>80 caractères)

## Statistiques

- **Couverture globale**: 91.0% 🟡
- **Packages 100%**: utils, formatter 🟢
- **Package const**: 92.9% 🟡
- **Package func**: Conforme 🟡
- **Package return**: 100% 🟢
- **Package interface**: 100% 🟢 (ignores struct interfaces)
- **Package comment**: 100% 🟢
- **Go version**: 1.25
- **Total règles**: 30 (4 const + 5 var + 12 func + 6 struct + 1 return + 1 interface + 1 comment)
- **Rapport détaillé**: Voir [COVERAGE.MD](COVERAGE.MD) pour le détail des fonctions < 100%

## Corrections des Contradictions

- ✅ **KTN-VAR-010 supprimé** : Contradictoire avec KTN-RETURN-002
- ✅ **KTN-COMMENT-001 supprimé** : Contradictoire avec KTN-FUNC-011 (demandait commentaires puis les marquait redondants)
- ✅ **KTN-INTERFACE-001 amélioré** : Ignore les interfaces qui suivent le pattern `XXXInterface` pour struct `XXX` (KTN-STRUCT-002)
- ✅ **KTN-VAR-014 amélioré** : Ignore les types externes (frameworks comme Terraform)
- ✅ **KTN-VAR-007 amélioré** : Ignore `[]T{}` (faux positifs), vérifie seulement `make([]T, 0)` sans capacity
- ✅ **KTN-FUNC-011 amélioré** : Ignore returns triviaux (nil, true, false, `[]T{}`)

## Structure

```
/workspace/
├── cmd/ktn-linter/     # Binaire
├── pkg/analyzer/       # Règles d'analyse
└── pkg/formatter/      # Formatage sortie
```

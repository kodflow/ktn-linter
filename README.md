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
make lint      # Lance le linter KTN
make help      # Aide
```

Voir [COVERAGE.MD](COVERAGE.MD) pour le rapport détaillé de couverture.

### Intégration VSCode

**Linting automatique au Ctrl+S** : L'extension Go de VSCode lance automatiquement le linter via `go.lintOnSave: "workspace"`.

**Workflow** :
1. Sauvegardez un fichier Go (`Ctrl+S`)
2. L'extension "Run on Save" rebuild le binaire automatiquement
3. L'extension Go lance le linter et affiche les erreurs dans l'onglet Problèmes (`Ctrl+Shift+M`)

**Fonctionnalités** :
- ✅ Build automatique du binaire à chaque sauvegarde
- ✅ Linting automatique via l'extension Go
- ✅ Affichage dans l'onglet Problèmes de VSCode
- ✅ Raccourci `Ctrl+Shift+L` pour forcer le lint manuellement

**Configuration** : `.vscode/settings.json`, `.vscode/tasks.json`, `.vscode/keybindings.json`
**Wrapper** : `bin/golangci-lint-wrapper` (utilisé par l'extension Go)

## Règles Implémentées

### Constantes (4 règles) ✅ 100%

- **KTN-CONST-001**: Type explicite obligatoire
- **KTN-CONST-002**: Groupement et placement avant var
- **KTN-CONST-003**: Nommage SCREAMING_SNAKE_CASE
- **KTN-CONST-004**: Commentaire obligatoire

### Fonctions (11 règles) ✅ 100%

- **KTN-FUNC-001**: Longueur max 35 lignes de code pur
- **KTN-FUNC-002**: Max 5 paramètres par fonction
- **KTN-FUNC-003**: Noms de fonctions commencent par un verbe
- **KTN-FUNC-004**: Pas de naked returns (sauf <5 lignes)
- **KTN-FUNC-005**: Complexité cyclomatique max 10
- **KTN-FUNC-006**: Erreur toujours en dernière position
- **KTN-FUNC-007**: Documentation stricte (Params/Returns)
- **KTN-FUNC-008**: Context toujours en premier paramètre
- **KTN-FUNC-009**: Pas de side effects dans les getters
- **KTN-FUNC-010**: Named returns pour >3 valeurs de retour
- **KTN-FUNC-011**: Commentaires sur branches/returns/logique

## Statistiques

- **Couverture globale**: 85.0% 🔴
- **Packages 100%**: const, ktn, utils, formatter 🟢
- **Package func**: 91.9% 🟡
- **Go version**: 1.25
- **Total règles**: 15 (4 const + 11 func)
- **Rapport détaillé**: Voir [COVERAGE.MD](COVERAGE.MD) pour le détail des fonctions < 100%

## Structure

```
/workspace/
├── cmd/ktn-linter/     # Binaire
├── pkg/analyzer/       # Règles d'analyse
└── pkg/formatter/      # Formatage sortie
```

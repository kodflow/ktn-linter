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

## Statistiques

- **Couverture globale**: 85.0% 🔴
- **Packages 100%**: const, ktn, utils, formatter 🟢
- **Package func**: 91.9% 🟡
- **Go version**: 1.25
- **Total règles**: 16 (4 const + 12 func)
- **Rapport détaillé**: Voir [COVERAGE.MD](COVERAGE.MD) pour le détail des fonctions < 100%

## Structure

```
/workspace/
├── cmd/ktn-linter/     # Binaire
├── pkg/analyzer/       # Règles d'analyse
└── pkg/formatter/      # Formatage sortie
```

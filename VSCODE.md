# Configuration VSCode pour KTN-Linter

## ✅ Configuration actuelle

Votre projet est **déjà configuré** pour afficher les diagnostics KTN-Linter directement dans VSCode.

### Comment ça fonctionne ?

1. **Wrapper** : `.vscode/settings.json` utilise `bin/golangci-lint-wrapper` comme outil de lint
2. **Format simple** : Le linter produit un format compatible VSCode (`file:line:col: message`)
3. **Lint au save** : Les erreurs apparaissent automatiquement quand vous sauvegardez un fichier Go

## 🔍 Vérifier que ça fonctionne

### 1. Ouvrir un fichier avec erreurs

Ouvrez par exemple :
```
tests/source/rules_const/package_level_KTN-CONST-001.go
```

### 2. Sauvegarder (Ctrl+S)

Les erreurs doivent apparaître dans :
- **Inline** : Lignes rouges/ondulées dans l'éditeur
- **Panneau Problems** : `Ctrl+Shift+M` ou cliquez sur "Problems" en bas

### 3. Format des erreurs

Vous devriez voir des messages comme :
```
[KTN-CONST-001] Constante 'EnableFeatureXC001' déclarée individuellement. Regroupez les constantes dans un bloc const ().
```

## 📋 Configuration VSCode

```json
{
  // Utiliser golangci-lint comme linter
  "go.lintTool": "golangci-lint",

  // Linter au save sur tout le workspace
  "go.lintOnSave": "workspace",

  // Utiliser notre wrapper au lieu de golangci-lint
  "go.alternateTools": {
    "golangci-lint": "${workspaceFolder}/bin/golangci-lint-wrapper"
  },

  // Build automatique au save
  "go.buildOnSave": "package"
}
```

## 🛠️ Build automatique du linter

Quand vous modifiez le code du linter dans `src/`, il se rebuild automatiquement grâce à l'extension **Run On Save**.

Pour installer cette extension :
```
code --install-extension emeraldwalk.RunOnSave
```

Ou cherchez "Run On Save" dans les extensions VSCode.

## 📊 Tester manuellement

### Tester le wrapper
```bash
./bin/golangci-lint-wrapper run ./tests/source/rules_const/package_level_KTN-CONST-001.go
```

### Tester le linter directement
```bash
./builds/ktn-linter ./tests/source/rules_const/package_level_KTN-CONST-001.go
```

### Format simple (pour IDE)
```bash
./builds/ktn-linter -simple ./tests/source/rules_const/package_level_KTN-CONST-001.go
```

## 🐛 Debug

### Voir les logs du wrapper
```bash
tail -f /tmp/ktn-linter-wrapper.log
```

### Rebuild le linter
```bash
go build -buildvcs=false -o builds/ktn-linter ./src/cmd/ktn-linter
```

### Vérifier que le wrapper est exécutable
```bash
chmod +x ./bin/golangci-lint-wrapper
```

## 📚 Règles disponibles

### CONST (4 règles)
- `KTN-CONST-001` : Constantes non groupées dans const ()
- `KTN-CONST-002` : Groupe sans commentaire de groupe
- `KTN-CONST-003` : Constante sans commentaire individuel
- `KTN-CONST-004` : Constante sans type explicite

### VAR (9 règles)
- `KTN-VAR-001` : Variables non groupées dans var ()
- `KTN-VAR-002` : Groupe sans commentaire de groupe
- `KTN-VAR-003` : Variable sans commentaire individuel
- `KTN-VAR-004` : Variable sans type explicite
- `KTN-VAR-005` : Variable jamais réassignée (devrait être const)
- `KTN-VAR-006` : Déclaration multiple sur une ligne
- `KTN-VAR-007` : Channel sans buffer size explicite
- `KTN-VAR-008` : Nom avec underscore (pas MixedCaps)
- `KTN-VAR-009` : Nom en ALL_CAPS (pas MixedCaps)

### FUNC (7 règles)
- `KTN-FUNC-001` : Nom en snake_case (pas MixedCaps)
- `KTN-FUNC-002` : Fonction sans commentaire godoc
- `KTN-FUNC-003` : Commentaire sans section Params: (format strict)
- `KTN-FUNC-004` : Commentaire sans section Returns: (format strict)
- `KTN-FUNC-005` : Trop de paramètres (> 5)
- `KTN-FUNC-006` : Fonction trop longue (> 35 lignes)
- `KTN-FUNC-007` : Complexité cyclomatique trop élevée (≥ 10)

## 🎯 Format de commentaire FUNC strict

**OBLIGATOIRE** pour toutes les fonctions :

```go
// FuncName description courte.
//
// Params:
//   - param1: description du param1
//   - param2: description du param2
//
// Returns:
//   - type: description du retour
func FuncName(param1 string, param2 int) error {
    // ...
}
```

**Sections conditionnelles :**
- Section `Params:` obligatoire SEULEMENT si la fonction a des paramètres
- Section `Returns:` obligatoire SEULEMENT si la fonction a des retours

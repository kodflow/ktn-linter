# KTN-Linter - Configuration Claude Code

## ⚠️ RÈGLES ABSOLUES

1. ❌ **INTERDICTION** : Créer des fichiers .md sauf `/workspace/README.md`
2. ❌ **INTERDICTION** : Générer des rapports/docs dans des dossiers
3. ✅ **SEULE EXCEPTION** : Mettre à jour `/workspace/README.md` avec format :
   - `KTN-XXX-YYY: Description minimaliste`
   - Informations pertinentes uniquement
   - Pas de contenu superflu

## Workflow Itératif Obligatoire

À chaque itération de développement :

1. **Écrire/Modifier le code**
2. **Hook automatique** → `make test` s'exécute automatiquement après chaque modification
3. **Corriger TOUS les warnings/errors/info**
4. **Vérifier la couverture** → Coverage maximale
5. **Mettre à jour README.md** si nouvelle règle
6. **Répéter jusqu'à 0 erreur**

### Agents Parallèles Post-Modification

Après chaque modification importante, Claude lance **2 agents en parallèle** :

**Agent 1 - Test Runner** :
```
Task("Exécuter tests avec couverture", "make test", "general-purpose")
```

**Agent 2 - Lint Runner** :
```
Task("Linter le projet", "make lint", "general-purpose")
```

**Avantage** : Les deux tâches s'exécutent simultanément pour un feedback plus rapide !

## Commandes Essentielles

```bash
# Aide
make help

# Tests avec couverture
make test

# Linter sur le projet
make lint
```

## Structure du Projet

```
/workspace/
├── go.mod
├── Makefile
├── CLAUDE.md
├── cmd/
│   └── ktn-linter/
│       └── main.go              # Point d'entrée du linter
├── pkg/
│   ├── analyzer/
│   │   ├── ktn/
│   │   │   ├── registry.go      # Enregistrement global des catégories
│   │   │   └── <category>/
│   │   │       ├── 001.go       # Règle KTN-<CATEGORY>-001
│   │   │       ├── 001_test.go  # Tests de la règle 001
│   │   │       ├── registry.go  # Enregistrement des règles de la catégorie
│   │   │       └── testdata/
│   │   │           └── src/
│   │   │               └── <category>001/
│   │   │                   ├── good.go
│   │   │                   └── bad.go
│   │   └── utils/           # Fonctions utilitaires
│   └── formatter/           # Formatage de la sortie
└── builds/                  # Binaires compilés (git-ignoré)
```

### Template d'une Règle (XXX.go)

```go
package ktn<category>

import (
    "go/ast"
    "golang.org/x/tools/go/analysis"
    "golang.org/x/tools/go/analysis/passes/inspect"
    "golang.org/x/tools/go/ast/inspector"
)

var Analyzer<XXX> = &analysis.Analyzer{
    Name:     "ktn<category><XXX>",
    Doc:      "KTN-<CATEGORY>-<XXX>: Description de la règle",
    Run:      run<Category><XXX>,
    Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run<Category><XXX>(pass *analysis.Pass) (any, error) {
    inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

    nodeFilter := []ast.Node{
        (*ast.FuncDecl)(nil), // Type de nœud AST à analyser
    }

    inspect.Preorder(nodeFilter, func(n ast.Node) {
        // Logique d'analyse
        // Si erreur détectée:
        pass.Reportf(n.Pos(), "KTN-<CATEGORY>-<XXX>: message d'erreur")
    })

    return nil, nil
}
```

### Template du Test (XXX_test.go)

```go
package ktn<category>_test

import (
    "testing"
    "golang.org/x/tools/go/analysis/analysistest"
    "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/<category>"
)

func Test<Category><XXX>(t *testing.T) {
    testdata := analysistest.TestData()
    analysistest.Run(t, testdata, ktn<category>.Analyzer<XXX>, "<category><XXX>")
}
```

### Template testdata (bad.go)

```go
package <category><XXX>

// Exemples de code qui DOIVENT générer des erreurs
func BadExample() { // want "KTN-<CATEGORY>-<XXX>: message d'erreur"
    // Code non-conforme
}
```

### Template testdata (good.go)

```go
package <category><XXX>

// Exemples de code conformes (pas d'erreur attendue)
func GoodExample() {
    // Code conforme
}
```

### Enregistrement dans registry.go

```go
package ktn<category>

import "golang.org/x/tools/go/analysis"

func Analyzers() []*analysis.Analyzer {
    return []*analysis.Analyzer{
        Analyzer001,
        Analyzer002,
        Analyzer<XXX>, // Ajouter la nouvelle règle
    }
}
```

## Règles de Développement

1. **Tests d'abord** : Écrire `XXX_test.go` et `testdata/` avant `XXX.go`
2. **Couverture obligatoire** : Chaque règle doit avoir une couverture maximale
3. **Convention de nommage** :
   - Fichiers : `001.go`, `002.go`, etc.
   - Analyzers : `Analyzer001`, `Analyzer002`, etc.
   - Tests : `001_test.go`, `002_test.go`, etc.
4. **Organisation des fichiers** :
   - Source : `/cmd/` (binaires), `/pkg/` (packages)
   - Tests : à côté du code, suffixe `_test.go`
   - Testdata : `/pkg/analyzer/ktn/<category>/testdata/`
   - Build : `/builds/` (généré, git-ignoré)
   - Coverage : `/coverage.out`, `/coverage.html` (généré, git-ignoré)
5. **Documentation** :
   - ❌ **INTERDIT** : Créer des fichiers .md ailleurs qu'à la racine
   - ✅ **AUTORISÉ** : Mettre à jour `/workspace/README.md` uniquement
   - 📝 **Format README** : `KTN-XXX-YYY: Description courte` (pas de blabla)
6. **Configuration golangci-lint** :
   - Les fichiers `*_test.go` sont exclus du linting (`.golangci.yml`)
   - Les règles ne s'appliquent que sur le code de production

## Catégories Disponibles

```
const, func, var, struct, interface, error, test,
alloc, goroutine, pool, mock, method, package,
control_flow, data_structures, ops
```

## Cycle de Développement d'une Nouvelle Règle

```bash
# 1. Créer la structure
touch pkg/analyzer/ktn/<category>/<XXX>.go
touch pkg/analyzer/ktn/<category>/<XXX>_test.go
mkdir -p pkg/analyzer/ktn/<category>/testdata/src/<category><XXX>
touch pkg/analyzer/ktn/<category>/testdata/src/<category><XXX>/good.go
touch pkg/analyzer/ktn/<category>/testdata/src/<category><XXX>/bad.go

# 2. Implémenter les tests (testdata + XXX_test.go)
# 3. Implémenter la règle (XXX.go)
# 4. Ajouter dans pkg/analyzer/ktn/<category>/registry.go
# 5. Ajouter la catégorie dans pkg/analyzer/ktn/registry.go si nouvelle

# 6. Vérifier
make test              # Tests doivent passer

# 7. Si erreurs/warnings/coverage < 100% → CORRIGER obligatoirement
# 8. Répéter jusqu'à 0 erreur et coverage élevée
```

## Exemple Complet : KTN-CONST-001

Voir `/pkg/analyzer/ktn/const/001.go` pour un exemple de règle complète.

## Points d'Attention

- ❌ **JAMAIS** sauvegarder de fichiers de travail/tests dans `/workspace/` (root)
- ❌ **JAMAIS** créer de fichiers .md n'importe où (pas de docs/, rapports/, etc.)
- ✅ **SEUL** fichier .md autorisé : `/workspace/README.md`
- 📝 **README.md** : Uniquement numéro de règle + description minimaliste + infos pertinentes
- ✅ **TOUJOURS** respecter l'arborescence `/cmd/`, `/pkg/`
- ⚠️ **OBLIGATOIRE** : Corriger tous les diagnostics avant de passer à la suite
- 📊 **OBJECTIF** : Couverture de tests maximale
- 🔄 **WORKFLOW** : Code → Test → Lint → Fix → Repeat

## Indicateurs de Qualité

**État actuel du projet :**

- ✅ `make test` : **94 tests PASS** (0 échec)
- 📊 **Coverage globale** : **76.8%**
  - `pkg/analyzer/utils` : **100%** ✅
  - `pkg/formatter` : **100%** ✅
  - `pkg/analyzer/ktn/const` : **92.9%** ✅
  - `cmd/ktn-linter` : **0%** (code CLI, normal)
- ⚠️ `make lint` : **18 erreurs** (constantes à renommer en SCREAMING_SNAKE_CASE)

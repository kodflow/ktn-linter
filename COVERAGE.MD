# Coverage Report

Rapport de couverture généré automatiquement.

**Légende:**
- 🟢 100% - Couverture complète
- 🟡 ≥90% - Bonne couverture
- 🔴 <90% - Couverture insuffisante
- ⚫ 0% - Pas de tests

---

## Coverage par Package

| Icon | Package | Coverage |
|:----:|---------|----------|
| 🔴 | **TOTAL (Global)** | **88.3%** |
| 🟢 | `github.com/kodflow/ktn-linter/cmd/ktn-linter` | 100.0% |
| 🔴 | `github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd` | 72.6% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn` | 100.0% |
| 🔴 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment` | 89.8% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst` | 94.1% |
| 🔴 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc` | 89.5% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktninterface` | 93.0% |
| 🔴 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnreturn` | 87.3% |
| 🔴 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct` | 88.9% |
| 🔴 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest` | 84.7% |
| 🔴 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar` | 87.4% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper` | 100.0% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/modernize` | 100.0% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/shared` | 98.2% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/utils` | 94.4% |
| 🔴 | `github.com/kodflow/ktn-linter/pkg/config` | 79.7% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/formatter` | 100.0% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/severity` | 100.0% |

---

## Détail des packages incomplets

### 🔴 `github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd` - 72.6%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `lint.go:loadConfiguration` | 30.0% |
| 🔴 `lint.go:runLint` | 57.1% |
| 🔴 `lint.go:loadPackages` | 71.4% |
| 🔴 `lint.go:checkLoadErrors` | 88.9% |
| 🟡 `lint.go:runAnalyzers` | 96.2% |
| 🔴 `lint.go:runRequiredAnalyzers` | 80.0% |
| 🔴 `lint.go:extractDiagnostics` | 87.5% |
| 🔴 `lint.go:applyFixes` | 80.0% |
| 🔴 `lint.go:collectSafeEdits` | 76.9% |
| 🔴 `lint.go:extractTextEdits` | 18.2% |
| 🔴 `lint.go:applyCollectedEdits` | 57.1% |
| 🔴 `lint.go:applyEditsToFile` | 18.2% |
| 🟡 `lint.go:filterOverlappingEdits` | 91.7% |

### 🔴 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment` - 89.8%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `001.go:runComment001` | 80.0% |
| 🟡 `001.go:isDocComment` | 91.7% |
| 🔴 `002.go:runComment002` | 78.6% |
| 🟡 `003.go:runComment003` | 90.9% |
| 🔴 `004.go:runComment004` | 88.5% |
| 🔴 `005.go:runComment005` | 80.0% |
| 🔴 `006.go:runComment006` | 76.0% |
| 🔴 `006.go:validateDocFormat` | 86.4% |
| 🔴 `007.go:runComment007` | 84.0% |
| 🔴 `007.go:hasCommentBeforeOrInside` | 71.4% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst` - 94.1%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `001.go:runConst001` | 89.5% |
| 🔴 `002.go:runConst002` | 80.0% |
| 🟡 `003.go:runConst003` | 90.5% |

### 🔴 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc` - 89.5%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `001.go:runFunc001` | 86.7% |
| 🔴 `001.go:isErrorType` | 76.9% |
| 🔴 `001.go:isBuiltinError` | 80.0% |
| 🔴 `002.go:runFunc002` | 86.7% |
| 🔴 `002.go:isContextTypeWithPass` | 75.0% |
| 🔴 `002.go:isContextTypeByType` | 88.9% |
| 🔴 `002.go:isContextUnderlying` | 50.0% |
| ⚫ `002.go:isContextType` | 0.0% |
| 🟡 `003.go:runFunc003` | 90.0% |
| 🔴 `003.go:isPanicCall` | 71.4% |
| 🔴 `004.go:runFunc004` | 88.9% |
| 🔴 `004.go:collectPrivateFunctions` | 84.6% |
| 🟡 `004.go:extractPrivateFuncInfo` | 91.7% |
| 🔴 `004.go:collectCalledFunctions` | 83.3% |
| 🔴 `004.go:collectCallbackUsages` | 81.8% |
| 🔴 `004.go:reportUnusedFunc` | 50.0% |
| 🔴 `005.go:runFunc005` | 81.8% |
| 🔴 `005.go:countPureCodeLines` | 88.2% |
| 🟡 `006.go:runFunc006` | 92.3% |
| 🔴 `007.go:runFunc007` | 84.2% |
| 🟡 `007.go:collectLazyLoadFields` | 92.9% |
| 🔴 `007.go:extractFieldName` | 80.0% |
| 🔴 `008.go:runFunc008` | 83.3% |
| 🔴 `008.go:analyzeFunc008` | 88.9% |
| 🔴 `008.go:findImplementedInterface` | 85.7% |
| 🔴 `008.go:getReceiverType` | 83.3% |
| 🔴 `008.go:findInterfaceForMethod` | 88.9% |
| 🔴 `008.go:collectFunctionParams008` | 87.5% |
| 🟡 `008.go:collectIgnoredVariables008` | 93.3% |
| 🟡 `008.go:findParentAssignToBlank008` | 95.0% |
| 🔴 `009.go:runFunc009` | 87.5% |
| 🟡 `009.go:checkMagicNumbers` | 92.3% |
| 🟡 `010.go:runFunc010` | 92.9% |
| 🔴 `011.go:runFunc011` | 85.0% |
| 🟡 `012.go:runFunc012` | 92.3% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktninterface` - 93.0%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `001.go:runInterface001` | 90.0% |
| 🟡 `001.go:collectDeclarations` | 90.0% |
| 🟡 `001.go:collectTypeSpecs` | 90.0% |
| 🔴 `001.go:checkEmbeddedInterfaces` | 80.0% |
| 🔴 `001.go:reportUnusedInterfaces` | 87.5% |

### 🔴 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnreturn` - 87.3%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `002.go:runReturn002` | 88.2% |
| 🔴 `002.go:isSliceOrMapType` | 83.3% |
| 🟡 `002.go:checkNilReturns` | 92.3% |
| 🔴 `002.go:collectSliceMapReturnTypes` | 80.0% |
| 🔴 `002.go:getSuggestionForType` | 85.7% |

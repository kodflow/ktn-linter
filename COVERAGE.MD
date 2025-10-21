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
| 🟡 | **TOTAL (Global)** | **91.7%** |
| 🟢 | `github.com/kodflow/ktn-linter/cmd/ktn-linter` | 100.0% |
| 🟡 | `github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd` | 94.8% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn` | 90.0% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/const` | 94.3% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/func` | 94.6% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper` | 100.0% |
| 🔴 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/var` | 83.3% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/utils` | 100.0% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/formatter` | 100.0% |

---

## Détail des packages incomplets

### 🟡 `github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd` - 94.8%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `lint.go:runLint` | 85.7% |
| 🔴 `lint.go:loadPackages` | 71.4% |
| 🟡 `lint.go:runAnalyzers` | 95.8% |
| 🟡 `lint.go:createAnalysisPass` | 90.0% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn` - 90.0%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `registry.go:GetRulesByCategory` | 80.0% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/const` - 94.3%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `002.go:checkConstGrouping` | 69.2% |
| 🟡 `003.go:runConst003` | 93.3% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/func` - 94.6%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `001.go:runFunc001` | 93.3% |
| 🔴 `001.go:isLineToSkip` | 64.3% |
| 🔴 `001.go:countPureCodeLines` | 88.2% |
| 🟡 `003.go:checkMagicNumbers` | 90.0% |
| 🟡 `005.go:runFunc005` | 92.3% |
| 🟡 `009.go:runFunc009` | 95.0% |
| 🟡 `011.go:runFunc011` | 94.4% |
| 🔴 `011.go:checkIfStmt` | 60.0% |
| ⚫ `011.go:hasCommentBeforeOrInside` | 0.0% |
| 🟡 `012.go:runFunc012` | 92.3% |

### 🔴 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/var` - 83.3%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `003.go:runVar003` | 93.3% |
| 🟡 `004.go:hasValidComment` | 90.0% |
| 🔴 `005.go:runVar005` | 87.5% |
| 🟡 `005.go:checkStatement` | 90.0% |
| 🔴 `005.go:checkNestedBlocks` | 60.0% |
| 🔴 `005.go:checkIfStmt` | 75.0% |
| ⚫ `005.go:checkCommClause` | 0.0% |
| 🟡 `006.go:runVar006` | 93.8% |
| 🔴 `007.go:isSliceType` | 50.0% |
| 🔴 `007.go:checkMakeCall` | 85.7% |
| 🔴 `008.go:isSliceTypeVar008` | 50.0% |
| 🔴 `008.go:hasPositiveLength` | 57.1% |
| 🔴 `008.go:checkMakeCallVar008` | 77.8% |
| 🔴 `009.go:runVar009` | 83.3% |
| 🔴 `009.go:isMakeCall` | 75.0% |
| 🔴 `009.go:isMapType` | 71.4% |
| 🟡 `010.go:runVar010` | 90.0% |
| 🔴 `010.go:isEmptySliceLiteral` | 66.7% |
| ⚫ `011.go:isBuilderType` | 0.0% |
| 🔴 `011.go:isBuilderCompositeLit` | 83.3% |
| 🔴 `011.go:extractTypeString` | 60.0% |
| 🔴 `011.go:extractAssignTypeString` | 83.3% |
| 🟡 `012.go:checkStringConcatInLoop` | 91.7% |
| 🔴 `012.go:isStringConcatenation` | 71.4% |
| 🟡 `013.go:runVar013` | 90.0% |
| 🔴 `013.go:checkLoopBodyForAlloc` | 75.0% |
| 🔴 `013.go:checkStmtForAlloc` | 66.7% |
| ⚫ `013.go:checkDeclForAlloc` | 0.0% |
| 🔴 `013.go:isSliceOrMapType` | 66.7% |
| 🔴 `014.go:runVar014` | 87.5% |
| 🔴 `014.go:checkDeclForLargeStruct` | 72.7% |
| 🔴 `014.go:checkTypeForLargeStruct` | 81.8% |
| 🔴 `015.go:checkLoopBodyVar015` | 83.3% |
| 🔴 `015.go:checkMakeCallForByteSlice` | 60.0% |
| 🔴 `015.go:isByteSlice` | 71.4% |
| 🔴 `016.go:shouldUseArray` | 75.0% |
| 🔴 `016.go:getConstantSize` | 66.7% |
| ⚫ `registry.go:Analyzers` | 0.0% |


---

*Généré le: 2025-10-21 13:29:59*

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
| 🟡 | **TOTAL (Global)** | **92.7%** |
| 🟢 | `github.com/kodflow/ktn-linter/cmd/ktn-linter` | 100.0% |
| 🟡 | `github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd` | 95.9% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn` | 90.0% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/const` | 94.4% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/func` | 96.1% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper` | 100.0% |
| 🔴 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/var` | 87.2% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/utils` | 100.0% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/formatter` | 100.0% |

---

## Détail des packages incomplets

### 🟡 `github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd` - 95.9%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `lint.go:loadPackages` | 71.4% |
| 🟡 `lint.go:runAnalyzers` | 96.2% |
| 🟡 `lint.go:createAnalysisPass` | 90.0% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn` - 90.0%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `registry.go:GetRulesByCategory` | 80.0% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/const` - 94.4%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `002.go:checkConstGrouping` | 69.2% |
| 🟡 `003.go:runConst003` | 93.3% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/func` - 96.1%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `001.go:runFunc001` | 93.3% |
| 🔴 `001.go:isLineToSkip` | 64.3% |
| 🔴 `001.go:countPureCodeLines` | 88.2% |
| 🟡 `003.go:checkMagicNumbers` | 90.0% |
| 🟡 `005.go:runFunc005` | 92.3% |
| 🟡 `009.go:runFunc009` | 95.0% |
| 🟡 `011.go:runFunc011` | 94.4% |
| 🔴 `011.go:checkIfStmt` | 80.0% |
| 🟡 `011.go:hasInlineComment` | 92.3% |
| 🔴 `011.go:hasCommentBeforeOrInside` | 85.7% |
| 🟡 `012.go:runFunc012` | 92.3% |

### 🔴 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/var` - 87.2%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `003.go:runVar003` | 93.3% |
| 🟡 `004.go:runVar004` | 95.0% |
| 🟡 `004.go:hasValidComment` | 90.0% |
| 🔴 `005.go:runVar005` | 87.5% |
| 🟡 `005.go:checkStatement` | 90.0% |
| 🔴 `005.go:checkNestedBlocks` | 80.0% |
| 🔴 `005.go:checkIfStmt` | 75.0% |
| 🟡 `006.go:runVar006` | 92.9% |
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
| 🔴 `011.go:isBuilderCompositeLit` | 85.7% |
| 🔴 `011.go:extractTypeString` | 61.5% |
| 🔴 `011.go:extractAssignTypeString` | 87.5% |
| 🟡 `012.go:checkStringConcatInLoop` | 91.7% |
| 🔴 `012.go:isStringConcatenation` | 75.0% |
| 🟡 `013.go:runVar013` | 90.0% |
| 🔴 `013.go:checkLoopBodyForAlloc` | 75.0% |
| 🔴 `013.go:checkDeclForAlloc` | 81.8% |
| 🔴 `013.go:isSliceOrMapType` | 66.7% |
| 🔴 `014.go:runVar014` | 87.5% |
| 🔴 `014.go:checkDeclForLargeStruct` | 75.0% |
| 🔴 `014.go:checkTypeForLargeStruct` | 81.8% |
| 🔴 `015.go:checkLoopBodyVar015` | 83.3% |
| 🔴 `015.go:checkMakeCallForByteSlice` | 60.0% |
| 🔴 `015.go:isByteSlice` | 71.4% |
| 🔴 `016.go:shouldUseArray` | 75.0% |
| 🔴 `016.go:getConstantSize` | 66.7% |
| 🔴 `017.go:checkShortVarDecl` | 88.9% |
| 🔴 `017.go:extractIdent` | 66.7% |
| 🔴 `017.go:isShadowing` | 77.8% |
| 🔴 `018.go:checkFuncForRepeatedConversions` | 80.0% |
| 🟡 `018.go:isStringConversion` | 90.9% |
| 🔴 `018.go:extractVarName` | 40.0% |
| 🟡 `019.go:collectTypesWithValueReceivers` | 92.9% |
| 🔴 `019.go:checkStructsWithMutex` | 84.6% |
| 🟡 `019.go:checkValueReceivers` | 92.3% |
| 🔴 `019.go:checkValueParams` | 88.9% |
| 🔴 `019.go:checkAssignments` | 88.9% |
| 🔴 `019.go:getMutexType` | 75.0% |
| 🔴 `019.go:getMutexTypeName` | 84.6% |
| 🔴 `019.go:hasMutex` | 75.0% |
| 🔴 `019.go:hasMutexInType` | 84.6% |
| 🔴 `019.go:getMutexTypeFromType` | 75.0% |
| 🔴 `019.go:getTypeName` | 40.0% |
| 🔴 `019.go:isMutexCopy` | 85.7% |


---

*Généré le: 2025-10-23 19:43:02*

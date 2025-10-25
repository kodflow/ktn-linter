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
| 🟡 | **TOTAL (Global)** | **91.4%** |
| 🟢 | `github.com/kodflow/ktn-linter/cmd/ktn-linter` | 100.0% |
| 🟡 | `github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd` | 94.2% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn` | 92.9% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst` | 97.8% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc` | 98.5% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct` | 91.3% |
| 🔴 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest` | 77.5% |
| 🔴 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar` | 85.8% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper` | 100.0% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/utils` | 94.4% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/formatter` | 100.0% |

---

## Détail des packages incomplets

### 🟡 `github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd` - 94.2%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `lint.go:runLint` | 85.7% |
| 🔴 `lint.go:loadPackages` | 71.4% |
| 🔴 `lint.go:checkLoadErrors` | 88.9% |
| 🟡 `lint.go:runAnalyzers` | 96.2% |
| 🔴 `lint.go:runRequiredAnalyzers` | 83.3% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn` - 92.9%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `registry.go:GetRulesByCategory` | 85.7% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst` - 97.8%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `002.go:checkConstGrouping` | 84.6% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc` - 98.5%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `001.go:runFunc001` | 93.3% |
| 🔴 `001.go:countPureCodeLines` | 88.2% |
| 🟡 `005.go:runFunc005` | 92.3% |
| 🟡 `009.go:runFunc009` | 95.0% |
| 🟡 `011.go:runFunc011` | 94.4% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct` - 91.3%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `001.go:runStruct001` | 90.0% |
| 🟡 `001.go:collectStructs` | 92.9% |
| 🔴 `002.go:runStruct002` | 84.6% |
| 🟡 `002.go:collectInterfaces` | 94.4% |
| 🟡 `002.go:collectStructsWithMethods` | 92.9% |
| 🔴 `002.go:formatFieldList` | 85.7% |
| ⚫ `registry.go:GetAnalyzers` | 0.0% |

### 🔴 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest` - 77.5%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `001.go:runTest001` | 92.3% |
| 🔴 `001.go:isExemptPackage` | 80.0% |
| 🟡 `002.go:runTest002` | 90.9% |
| 🔴 `002.go:isExemptTestFile` | 83.3% |
| 🔴 `003.go:runTest003` | 88.9% |
| 🔴 `003.go:isPublicFunction` | 66.7% |
| 🔴 `003.go:collectTestedFunctions` | 75.0% |
| 🔴 `003.go:isExemptFunction` | 80.0% |
| 🔴 `004.go:runTest004` | 83.3% |
| 🔴 `004.go:isTestFunction` | 75.0% |
| 🔴 `004.go:hasErrorCaseCoverage` | 76.9% |
| ⚫ `004.go:hasErrorTestCases` | 0.0% |
| 🔴 `004.go:isErrorIndicatorName` | 83.3% |
| 🔴 `005.go:runTest005` | 83.3% |
| 🔴 `005.go:isTestFunc` | 75.0% |
| ⚫ `005.go:isTestsVariableName` | 0.0% |
| ⚫ `005.go:checkAssignStmt` | 0.0% |
| ⚫ `005.go:checkRangeStmt` | 0.0% |
| 🔴 `005.go:hasTableDrivenPattern` | 60.0% |
| 🟡 `006.go:runTest006` | 94.1% |

### 🔴 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar` - 85.8%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `003.go:runVar003` | 93.3% |
| 🟡 `004.go:runVar004` | 95.0% |
| 🔴 `004.go:hasValidComment` | 70.0% |
| 🔴 `005.go:runVar005` | 87.5% |
| 🟡 `005.go:checkStatement` | 90.0% |
| 🔴 `005.go:checkNestedBlocks` | 80.0% |
| 🔴 `005.go:checkIfStmt` | 75.0% |
| 🟡 `006.go:runVar006` | 92.9% |
| 🔴 `007.go:checkMakeCall` | 85.7% |
| 🔴 `008.go:checkMakeCallVar008` | 77.8% |
| 🔴 `009.go:runVar009` | 83.3% |
| 🔴 `011.go:checkValueSpec` | 71.4% |
| 🔴 `011.go:isBuilderCompositeLit` | 85.7% |
| 🔴 `011.go:reportMissingGrow` | 83.3% |
| ⚫ `011.go:extractTypeString` | 0.0% |
| 🔴 `011.go:extractAssignTypeString` | 87.5% |
| 🟡 `012.go:checkStringConcatInLoop` | 91.7% |
| 🔴 `012.go:isStringConcatenation` | 75.0% |
| 🟡 `013.go:runVar013` | 90.0% |
| 🔴 `013.go:checkLoopBodyForAlloc` | 75.0% |
| 🔴 `013.go:checkStmtForAlloc` | 66.7% |
| ⚫ `013.go:checkDeclForAlloc` | 0.0% |
| 🔴 `013.go:isSliceOrMapAlloc` | 66.7% |
| 🔴 `014.go:runVar014` | 87.5% |
| 🔴 `014.go:checkDeclForLargeStruct` | 75.0% |
| 🔴 `014.go:checkTypeForLargeStruct` | 81.8% |
| 🔴 `015.go:checkLoopBodyVar015` | 83.3% |
| 🔴 `015.go:checkMakeCallForByteSlice` | 60.0% |
| 🔴 `016.go:shouldUseArray` | 75.0% |
| 🔴 `016.go:getConstantSize` | 66.7% |
| 🔴 `017.go:checkShortVarDecl` | 88.9% |
| 🔴 `017.go:extractIdent` | 66.7% |
| 🔴 `017.go:isShadowing` | 77.8% |
| 🔴 `018.go:checkFuncForRepeatedConversions` | 80.0% |
| 🟡 `018.go:isStringConversion` | 90.9% |
| 🟡 `019.go:collectTypesWithValueReceivers` | 92.9% |
| 🟡 `019.go:checkStructsWithMutex` | 92.3% |
| 🟡 `019.go:checkValueReceivers` | 92.3% |
| 🔴 `019.go:checkValueParams` | 88.9% |
| 🔴 `019.go:checkAssignments` | 88.9% |
| 🔴 `019.go:getMutexType` | 75.0% |
| 🟡 `019.go:getMutexTypeName` | 92.3% |
| 🔴 `019.go:hasMutex` | 75.0% |
| 🟡 `019.go:hasMutexInType` | 92.3% |
| 🔴 `019.go:getMutexTypeFromType` | 75.0% |
| 🔴 `019.go:getTypeName` | 40.0% |
| 🔴 `019.go:isMutexCopy` | 85.7% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/utils` - 94.4%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `ident.go:ExtractVarName` | 76.9% |
| 🔴 `types.go:HasPositiveLength` | 36.4% |


---

*Généré le: 2025-10-25 10:47:45*

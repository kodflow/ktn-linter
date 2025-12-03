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
| 🔴 | `github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd` | 74.6% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn` | 100.0% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment` | 96.8% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst` | 100.0% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc` | 95.9% |
| 🔴 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktninterface` | 88.7% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnpackage` | 100.0% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnreturn` | 90.6% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct` | 91.8% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest` | 91.0% |
| 🔴 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar` | 87.1% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper` | 100.0% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/modernize` | 100.0% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/shared` | 97.7% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/utils` | 94.4% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/formatter` | 100.0% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/severity` | 100.0% |

---

## Détail des packages incomplets

### 🔴 `github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd` - 74.6%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `lint.go:runLint` | 53.8% |
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

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment` - 96.8%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `002.go:isDocComment` | 91.7% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc` - 95.9%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `001.go:runFunc001` | 93.3% |
| 🔴 `001.go:countPureCodeLines` | 88.2% |
| 🟡 `005.go:runFunc005` | 92.3% |
| 🔴 `007.go:validateDocFormat` | 86.4% |
| 🟡 `009.go:runFunc009` | 95.0% |
| 🟡 `011.go:runFunc011` | 94.1% |
| 🔴 `011.go:hasCommentBeforeOrInside` | 71.4% |
| 🟡 `013.go:runFunc013` | 94.4% |
| 🔴 `013.go:collectFunctionParams` | 87.5% |
| 🔴 `013.go:collectIgnoredVariables` | 86.7% |
| 🟡 `013.go:findParentAssignToBlank` | 90.0% |
| 🟡 `014.go:collectPrivateFunctions` | 90.9% |
| 🟡 `014.go:extractPrivateFuncInfo` | 91.7% |
| 🟡 `014.go:collectCalledFunctions` | 90.0% |
| 🟡 `014.go:collectCallbackUsages` | 90.0% |
| 🔴 `014.go:reportUnusedFunc` | 50.0% |

### 🔴 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktninterface` - 88.7%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `001.go:collectTypeSpecs` | 90.0% |
| 🔴 `001.go:checkEmbeddedInterfaces` | 80.0% |
| 🔴 `001.go:checkType` | 37.5% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnreturn` - 90.6%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `002.go:runReturn002` | 90.9% |
| 🔴 `002.go:isSliceOrMapType` | 83.3% |
| 🟡 `002.go:checkNilReturns` | 90.0% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct` - 91.8%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `001.go:runStruct001` | 83.3% |
| 🟡 `001.go:collectStructs` | 92.9% |
| 🔴 `002.go:runStruct002` | 86.7% |
| 🟡 `002.go:collectInterfaces` | 94.4% |
| 🔴 `003.go:runStruct003` | 83.3% |
| 🔴 `004.go:runStruct004` | 84.2% |
| 🟡 `004.go:hasValidDocumentation` | 92.9% |
| 🔴 `005.go:runStruct005` | 87.5% |
| 🔴 `005.go:extractReturnTypeName` | 87.5% |
| 🔴 `006.go:runStruct006` | 80.0% |
| 🔴 `006.go:isSimpleGetter` | 80.0% |
| 🔴 `006.go:getReceiverTypeName` | 85.7% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest` - 91.0%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `001.go:runTest001` | 92.3% |
| 🔴 `002.go:runTest002` | 83.3% |
| 🟡 `003.go:collectFunctions` | 92.3% |
| 🔴 `003.go:buildTestNames` | 75.0% |
| 🔴 `003.go:hasMatchingTest` | 85.7% |
| 🔴 `003.go:isPublicFunction` | 83.3% |
| 🔴 `003.go:extractReceiverTypeName` | 80.0% |
| 🔴 `003.go:collectExternalTestFunctions` | 88.0% |
| 🔴 `004.go:runTest004` | 83.3% |
| 🟡 `004.go:collectFuncSignatures` | 90.0% |
| 🔴 `004.go:addFuncSignature` | 75.0% |
| 🔴 `004.go:collectExternalSourceSignatures` | 72.7% |
| 🟡 `004.go:scanSourceFile` | 93.3% |
| 🔴 `004.go:extractFuncInfo` | 60.0% |
| 🔴 `004.go:functionReturnsError` | 87.5% |
| 🔴 `004.go:isErrorType` | 66.7% |
| 🔴 `004.go:extractReceiverName` | 80.0% |
| 🔴 `004.go:analyzeTestFunction` | 80.0% |
| 🔴 `004.go:checkErrorInKeyValue` | 75.0% |
| 🟡 `005.go:runTest005` | 91.7% |
| 🔴 `006.go:runTest006` | 87.5% |
| 🔴 `007.go:runTest007` | 88.9% |
| 🔴 `008.go:runTest008` | 83.3% |
| 🟡 `008.go:checkVariables` | 90.9% |
| 🔴 `008.go:reportTestFileIssues` | 75.0% |
| 🔴 `008.go:reportMixedFunctionsIssues` | 85.7% |
| 🔴 `008.go:reportPublicOnlyIssues` | 60.0% |
| 🔴 `008.go:reportPrivateOnlyIssues` | 60.0% |
| 🟡 `009.go:addPublicFunction` | 90.0% |
| 🔴 `009.go:checkAndReportPublicFunctionTest` | 80.0% |
| 🔴 `009.go:extractReceiverType` | 80.0% |
| 🔴 `010.go:checkAndReportPrivateFunctionTest` | 83.3% |
| 🔴 `013.go:runTest013` | 83.3% |
| 🔴 `013.go:isTestingMethodCall` | 75.0% |
| 🔴 `013.go:isAssertLibraryCall` | 75.0% |

### 🔴 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar` - 87.1%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `001.go:isTypeVisible` | 87.5% |
| 🔴 `001.go:isTypedCall` | 40.0% |
| 🟡 `003.go:runVar003` | 93.3% |
| 🟡 `004.go:runVar004` | 95.0% |
| 🔴 `005.go:runVar005` | 87.5% |
| 🟡 `005.go:checkStatement` | 90.0% |
| 🔴 `005.go:checkNestedBlocks` | 80.0% |
| 🔴 `005.go:checkIfStmt` | 75.0% |
| 🟡 `006.go:runVar006` | 92.9% |
| 🟡 `007.go:collectAppendVariables` | 91.7% |
| 🔴 `007.go:checkMakeCalls` | 83.3% |
| 🔴 `007.go:checkMakeCall` | 85.7% |
| 🟡 `007.go:checkEmptySliceLiterals` | 92.9% |
| 🔴 `007.go:checkCompositeLit` | 68.8% |
| 🔴 `008.go:checkMakeCallVar008` | 77.8% |
| 🔴 `009.go:runVar009` | 83.3% |
| 🔴 `011.go:checkValueSpec` | 71.4% |
| 🔴 `011.go:reportMissingGrow` | 83.3% |
| ⚫ `011.go:extractTypeString` | 0.0% |
| 🔴 `011.go:extractAssignTypeString` | 87.5% |
| 🟡 `012.go:checkStringConcatInLoop` | 91.7% |
| 🔴 `012.go:isStringConcatenation` | 75.0% |
| 🟡 `013.go:runVar013` | 90.0% |
| 🔴 `013.go:checkLoopBodyForAlloc` | 75.0% |
| 🔴 `013.go:checkStmtForAlloc` | 66.7% |
| ⚫ `013.go:checkDeclForAlloc` | 0.0% |
| 🔴 `014.go:runVar014` | 87.5% |
| 🔴 `014.go:checkDeclForLargeStruct` | 75.0% |
| 🔴 `014.go:checkTypeForLargeStruct` | 76.9% |
| 🔴 `014.go:isExternalType` | 71.4% |
| 🔴 `015.go:checkLoopBodyVar015` | 83.3% |
| 🔴 `015.go:checkMakeCallForByteSlice` | 60.0% |
| 🔴 `016.go:shouldUseArray` | 75.0% |
| 🔴 `016.go:getConstantSize` | 66.7% |
| 🔴 `017.go:checkShortVarDecl` | 88.9% |
| 🔴 `017.go:isShadowing` | 77.8% |
| 🔴 `018.go:checkFuncForRepeatedConversions` | 80.0% |
| 🟡 `019.go:collectTypesWithValueReceivers` | 92.9% |
| 🟡 `019.go:checkStructsWithMutex` | 92.3% |
| 🟡 `019.go:checkValueReceivers` | 92.3% |
| 🔴 `019.go:checkValueParams` | 88.9% |
| 🔴 `019.go:checkAssignments` | 88.9% |
| 🔴 `019.go:getMutexType` | 75.0% |
| 🔴 `019.go:hasMutex` | 75.0% |
| 🟡 `019.go:hasMutexInType` | 91.7% |
| 🔴 `019.go:getMutexTypeFromType` | 73.3% |
| 🔴 `019.go:isMutexCopy` | 85.7% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/shared` - 97.7%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `types.go:hasSerializationTags` | 83.3% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/utils` - 94.4%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `ident.go:ExtractVarName` | 76.9% |
| 🔴 `types.go:HasPositiveLength` | 36.4% |


---

*Généré le: 2025-12-03 21:45:12*

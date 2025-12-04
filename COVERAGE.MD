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
| 🟡 | **TOTAL (Global)** | **91.6%** |
| 🟢 | `github.com/kodflow/ktn-linter/cmd/ktn-linter` | 100.0% |
| 🔴 | `github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd` | 74.6% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn` | 100.0% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment` | 96.8% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst` | 100.0% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc` | 95.8% |
| 🔴 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktninterface` | 88.7% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnpackage` | 100.0% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnreturn` | 93.8% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct` | 92.1% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest` | 91.3% |
| 🔴 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar` | 87.2% |
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
| 🟡 `001.go:isDocComment` | 91.7% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc` - 95.8%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `004.go:collectPrivateFunctions` | 90.9% |
| 🟡 `004.go:extractPrivateFuncInfo` | 91.7% |
| 🟡 `004.go:collectCalledFunctions` | 90.0% |
| 🔴 `004.go:collectCallbackUsages` | 88.9% |
| 🔴 `004.go:reportUnusedFunc` | 50.0% |
| 🔴 `005.go:runFunc005` | 86.7% |
| 🔴 `005.go:countPureCodeLines` | 88.2% |
| 🔴 `007.go:validateDocFormat` | 86.4% |
| 🟡 `008.go:runFunc008` | 95.0% |
| 🟡 `009.go:runFunc009` | 94.1% |
| 🔴 `009.go:hasCommentBeforeOrInside` | 71.4% |
| 🟡 `010.go:runFunc010` | 94.4% |
| 🔴 `010.go:collectFunctionParams` | 87.5% |
| 🔴 `010.go:collectIgnoredVariables` | 86.7% |
| 🟡 `010.go:findParentAssignToBlank` | 90.0% |
| 🟡 `013.go:runFunc013` | 92.3% |

### 🔴 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktninterface` - 88.7%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `001.go:collectTypeSpecs` | 90.0% |
| 🔴 `001.go:checkEmbeddedInterfaces` | 80.0% |
| 🔴 `001.go:checkType` | 37.5% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnreturn` - 93.8%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `002.go:isSliceOrMapType` | 83.3% |
| 🟡 `002.go:checkNilReturns` | 90.0% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct` - 92.1%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `001.go:runStruct001` | 86.7% |
| 🟡 `001.go:collectInterfaces` | 94.4% |
| 🔴 `002.go:runStruct002` | 84.2% |
| 🟡 `002.go:hasValidDocumentation` | 92.9% |
| 🟡 `003.go:runStruct003` | 93.8% |
| 🔴 `003.go:extractReturnTypeName` | 87.5% |
| 🔴 `004.go:runStruct004` | 80.0% |
| 🔴 `004.go:isSimpleGetter` | 80.0% |
| 🔴 `004.go:getReceiverTypeName` | 85.7% |
| 🔴 `005.go:runStruct005` | 83.3% |
| 🟡 `005.go:collectStructs` | 92.9% |
| 🔴 `006.go:runStruct006` | 83.3% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest` - 91.3%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `002.go:runTest002` | 92.3% |
| 🔴 `003.go:runTest003` | 83.3% |
| 🟡 `004.go:collectFunctions` | 92.3% |
| 🔴 `004.go:buildTestNames` | 75.0% |
| 🔴 `004.go:hasMatchingTest` | 85.7% |
| 🔴 `004.go:isPublicFunction` | 83.3% |
| 🔴 `004.go:extractReceiverTypeName` | 80.0% |
| 🔴 `004.go:collectExternalTestFunctions` | 88.0% |
| 🟡 `005.go:runTest005` | 91.7% |
| 🔴 `006.go:runTest006` | 87.5% |
| 🔴 `007.go:runTest007` | 88.9% |
| 🔴 `008.go:runTest008` | 83.3% |
| 🔴 `008.go:extractReceiverTypeString` | 85.7% |
| 🟡 `008.go:checkVariables` | 90.9% |
| 🔴 `008.go:reportTestFileIssues` | 75.0% |
| 🔴 `008.go:reportMixedFunctionsIssues` | 88.9% |
| 🔴 `008.go:reportPublicOnlyIssues` | 66.7% |
| 🔴 `008.go:reportPrivateOnlyIssues` | 66.7% |
| 🟡 `009.go:addPublicFunction` | 90.0% |
| 🔴 `009.go:checkAndReportPublicFunctionTest` | 80.0% |
| 🔴 `009.go:extractReceiverType` | 80.0% |
| 🔴 `010.go:checkAndReportPrivateFunctionTest` | 83.3% |
| 🔴 `012.go:runTest012` | 83.3% |
| 🔴 `012.go:isTestingMethodCall` | 75.0% |
| 🔴 `012.go:isAssertLibraryCall` | 75.0% |
| 🔴 `013.go:runTest013` | 83.3% |
| 🟡 `013.go:collectFuncSignatures` | 90.0% |
| 🔴 `013.go:addFuncSignature` | 75.0% |
| 🔴 `013.go:collectExternalSourceSignatures` | 72.7% |
| 🟡 `013.go:scanSourceFile` | 93.3% |
| 🔴 `013.go:extractFuncInfo` | 60.0% |
| 🔴 `013.go:functionReturnsError` | 87.5% |
| 🔴 `013.go:isErrorType` | 66.7% |
| 🔴 `013.go:extractReceiverName` | 80.0% |
| 🔴 `013.go:analyzeTestFunction` | 80.0% |
| 🔴 `013.go:checkErrorInKeyValue` | 75.0% |

### 🔴 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar` - 87.2%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `001.go:runVar001` | 93.3% |
| 🟡 `002.go:runVar002` | 95.0% |
| 🔴 `003.go:isTypeVisible` | 87.5% |
| 🔴 `003.go:isTypedCall` | 40.0% |
| 🔴 `004.go:runVar004` | 87.5% |
| 🟡 `004.go:checkStatement` | 90.0% |
| 🔴 `004.go:checkNestedBlocks` | 80.0% |
| 🔴 `004.go:checkIfStmt` | 75.0% |
| 🟡 `005.go:collectAppendVariables` | 91.7% |
| 🔴 `005.go:checkMakeCalls` | 83.3% |
| 🔴 `005.go:checkMakeCall` | 85.7% |
| 🟡 `005.go:checkEmptySliceLiterals` | 92.9% |
| 🔴 `005.go:checkCompositeLit` | 62.5% |
| 🔴 `006.go:checkMakeCallVar008` | 77.8% |
| 🔴 `007.go:checkValueSpec` | 71.4% |
| 🔴 `007.go:reportMissingGrow` | 83.3% |
| ⚫ `007.go:extractTypeString` | 0.0% |
| 🔴 `007.go:extractAssignTypeString` | 87.5% |
| 🟡 `008.go:checkStringConcatInLoop` | 91.7% |
| 🔴 `008.go:isStringConcatenation` | 75.0% |
| 🟡 `009.go:runVar009` | 90.0% |
| 🔴 `009.go:checkLoopBodyForAlloc` | 75.0% |
| 🔴 `009.go:checkStmtForAlloc` | 66.7% |
| ⚫ `009.go:checkDeclForAlloc` | 0.0% |
| 🔴 `010.go:runVar010` | 87.5% |
| 🔴 `010.go:checkDeclForLargeStruct` | 75.0% |
| 🔴 `010.go:checkTypeForLargeStruct` | 76.9% |
| 🔴 `010.go:isExternalType` | 71.4% |
| 🔴 `011.go:checkLoopBodyVar015` | 83.3% |
| 🔴 `011.go:checkMakeCallForByteSlice` | 80.0% |
| 🔴 `012.go:checkShortVarDecl` | 88.9% |
| 🔴 `012.go:isShadowing` | 77.8% |
| 🔴 `013.go:checkFuncForRepeatedConversions` | 80.0% |
| 🟡 `015.go:runVar015` | 92.9% |
| 🟡 `016.go:runVar016` | 91.7% |
| 🔴 `017.go:shouldUseArray` | 75.0% |
| 🔴 `017.go:getConstantSize` | 66.7% |
| 🟡 `018.go:collectTypesWithValueReceivers` | 92.9% |
| 🟡 `018.go:checkStructsWithMutex` | 92.3% |
| 🟡 `018.go:checkValueReceivers` | 92.3% |
| 🔴 `018.go:checkValueParams` | 88.9% |
| 🔴 `018.go:checkAssignments` | 88.9% |
| 🔴 `018.go:getMutexType` | 75.0% |
| 🔴 `018.go:hasMutex` | 75.0% |
| 🟡 `018.go:hasMutexInType` | 91.7% |
| 🔴 `018.go:getMutexTypeFromType` | 73.3% |
| 🔴 `018.go:isMutexCopy` | 85.7% |

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

*Généré le: 2025-12-04 20:59:27*

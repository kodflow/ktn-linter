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
| 🔴 | **TOTAL (Global)** | **88.2%** |
| 🟢 | `github.com/kodflow/ktn-linter/cmd/ktn-linter` | 100.0% |
| 🔴 | `github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd` | 72.6% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn` | 100.0% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment` | 90.0% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst` | 94.1% |
| 🔴 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc` | 89.6% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktninterface` | 93.0% |
| 🔴 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnreturn` | 87.3% |
| 🔴 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct` | 88.9% |
| 🔴 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest` | 84.7% |
| 🔴 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar` | 87.0% |
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

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment` - 90.0%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `001.go:runComment001` | 80.0% |
| 🟡 `001.go:isDocComment` | 91.7% |
| 🔴 `002.go:runComment002` | 78.6% |
| 🟡 `003.go:runComment003` | 90.9% |
| 🔴 `004.go:runComment004` | 83.3% |
| 🔴 `004.go:checkFileDeclarations` | 85.7% |
| 🔴 `005.go:runComment005` | 80.0% |
| 🔴 `006.go:runComment006` | 76.0% |
| 🔴 `006.go:validateDocFormat` | 86.4% |
| 🟡 `007.go:runComment007` | 90.9% |
| 🔴 `007.go:shouldSkipFunction` | 70.0% |
| 🔴 `007.go:hasCommentBeforeOrInside` | 71.4% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst` - 94.1%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `001.go:runConst001` | 89.5% |
| 🔴 `002.go:runConst002` | 80.0% |
| 🟡 `003.go:runConst003` | 90.5% |

### 🔴 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc` - 89.6%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `001.go:runFunc001` | 86.7% |
| 🔴 `001.go:isErrorType` | 76.9% |
| 🔴 `001.go:isBuiltinError` | 80.0% |
| 🔴 `002.go:runFunc002` | 85.7% |
| 🟡 `002.go:analyzeContextParams` | 93.3% |
| 🔴 `002.go:reportMultipleContexts` | 50.0% |
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
| 🔴 `006.go:runFunc006` | 83.3% |
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
| 🔴 `010.go:runFunc010` | 88.2% |
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

### 🔴 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct` - 88.9%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `001.go:runStruct001` | 80.0% |
| 🟡 `001.go:collectInterfaces` | 94.4% |
| 🔴 `002.go:runStruct002` | 85.7% |
| 🔴 `002.go:extractReturnTypeName` | 87.5% |
| 🔴 `003.go:runStruct003` | 78.6% |
| 🔴 `003.go:isSimpleGetter` | 80.0% |
| 🔴 `003.go:getReceiverTypeName` | 85.7% |
| 🔴 `004.go:runStruct004` | 76.5% |
| 🟡 `004.go:collectStructs` | 92.9% |
| 🔴 `005.go:runStruct005` | 83.3% |
| 🔴 `006.go:runStruct006` | 82.4% |
| 🔴 `006.go:checkPrivateFieldsWithTags` | 87.5% |
| 🔴 `007.go:runStruct007` | 87.5% |
| 🟡 `007.go:collectStructPrivateFields` | 95.0% |
| 🔴 `007.go:collectMethodsDetailed` | 86.7% |
| 🔴 `007.go:extractSimpleReturnType` | 62.5% |
| 🔴 `007.go:extractReceiverType` | 80.0% |
| 🔴 `007.go:checkNamingConventions` | 85.7% |
| 🔴 `007.go:checkGetterFieldMismatch` | 73.3% |

### 🔴 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest` - 84.7%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `001.go:runTest001` | 82.4% |
| ⚫ `001.go:verifyBenchFile` | 0.0% |
| 🔴 `002.go:runTest002` | 83.3% |
| 🔴 `003.go:runTest003` | 81.2% |
| 🔴 `004.go:collectFunctions` | 87.5% |
| 🔴 `004.go:runTest004` | 88.9% |
| 🔴 `004.go:checkFunctionsHaveTests` | 84.6% |
| 🔴 `004.go:buildTestNames` | 83.3% |
| 🔴 `004.go:hasMatchingTest` | 85.7% |
| 🟡 `004.go:reportMissingTest` | 91.7% |
| 🔴 `004.go:findPackageDir` | 83.3% |
| 🔴 `004.go:collectExternalTestFunctions` | 83.3% |
| 🔴 `004.go:collectExternalTestedFunctions` | 86.7% |
| 🔴 `005.go:runTest005` | 82.4% |
| 🔴 `006.go:runTest006` | 82.8% |
| 🔴 `007.go:runTest007` | 82.6% |
| 🔴 `008.go:runTest008` | 73.7% |
| 🔴 `008.go:classifyFunction` | 85.7% |
| 🔴 `008.go:extractReceiverTypeString` | 85.7% |
| 🔴 `008.go:checkVariables` | 84.6% |
| 🔴 `008.go:reportTestFileIssues` | 75.0% |
| 🔴 `008.go:reportMixedFunctionsIssues` | 88.9% |
| 🔴 `008.go:reportPublicOnlyIssues` | 66.7% |
| 🔴 `008.go:reportPrivateOnlyIssues` | 66.7% |
| 🔴 `009.go:runTest009` | 85.7% |
| 🟡 `009.go:collectPublicFunctions` | 90.9% |
| 🔴 `009.go:addPublicFunction` | 75.0% |
| 🔴 `009.go:checkInternalTestsForPublicFunctions` | 84.6% |
| 🔴 `009.go:checkAndReportPublicFunctionTest` | 72.7% |
| 🔴 `010.go:runTest010` | 85.7% |
| 🟡 `010.go:collectPrivateFunctions` | 90.9% |
| 🔴 `010.go:addPrivateFunction` | 75.0% |
| 🔴 `010.go:checkExternalTestsForPrivateFunctions` | 84.6% |
| 🔴 `010.go:checkAndReportPrivateFunctionTest` | 76.9% |
| 🟡 `011.go:runTest011` | 90.9% |
| 🔴 `012.go:runTest012` | 73.9% |
| 🔴 `012.go:isAssertLibraryCall` | 87.5% |
| 🔴 `013.go:runTest013` | 69.6% |
| 🔴 `013.go:collectFuncSignatures` | 64.3% |
| 🔴 `013.go:addFuncSignature` | 66.7% |
| 🔴 `013.go:collectExternalSourceSignatures` | 72.7% |
| 🔴 `013.go:scanSourceFile` | 84.2% |
| 🔴 `013.go:functionReturnsError` | 87.5% |
| 🔴 `013.go:isErrorType` | 66.7% |
| 🔴 `013.go:analyzeTestFunction` | 76.9% |
| 🔴 `013.go:checkErrorInKeyValue` | 75.0% |

### 🔴 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar` - 87.0%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `001.go:runVar001` | 85.0% |
| 🔴 `002.go:runVar002` | 89.5% |
| 🔴 `002.go:checkVarSpec` | 83.3% |
| 🔴 `003.go:runVar003` | 76.9% |
| 🔴 `003.go:checkStatement` | 80.0% |
| 🔴 `003.go:checkNestedBlocks` | 80.0% |
| 🔴 `003.go:checkIfStmt` | 75.0% |
| 🔴 `004.go:runVar004` | 87.5% |
| 🟡 `004.go:collectAppendVariables` | 91.7% |
| 🔴 `004.go:checkMakeCalls` | 77.8% |
| 🔴 `004.go:checkMakeCall` | 77.8% |
| 🔴 `004.go:checkEmptySliceLiterals` | 88.2% |
| 🔴 `004.go:checkCompositeLit` | 68.8% |
| 🔴 `005.go:checkMakeCallVar008` | 72.7% |
| 🔴 `005.go:runVar005` | 81.8% |
| 🔴 `006.go:runVar006` | 80.0% |
| 🔴 `006.go:reportMissingGrow` | 83.3% |
| 🔴 `006.go:extractAssignTypeString` | 87.5% |
| 🔴 `007.go:runVar007` | 80.0% |
| 🟡 `007.go:checkStringConcatInLoop` | 91.7% |
| 🔴 `007.go:isStringConcatenation` | 75.0% |
| 🔴 `008.go:runVar008` | 80.0% |
| 🔴 `008.go:checkLoopBodyForAlloc` | 75.0% |
| 🔴 `008.go:checkStmtForAlloc` | 66.7% |
| 🟡 `008.go:checkDeclForAlloc` | 90.9% |
| 🔴 `009.go:runVar009` | 78.6% |
| 🔴 `009.go:checkDeclForLargeStruct` | 75.0% |
| 🔴 `009.go:checkTypeForLargeStruct` | 76.9% |
| 🔴 `009.go:isExternalType` | 71.4% |
| 🔴 `010.go:runVar010` | 83.3% |
| 🔴 `010.go:checkLoopBodyVar015` | 83.3% |
| 🔴 `010.go:checkMakeCallForByteSlice` | 80.0% |
| 🔴 `011.go:runVar011` | 80.0% |
| 🔴 `011.go:checkShortVarDecl` | 88.9% |
| 🔴 `011.go:isShadowing` | 81.8% |
| 🔴 `012.go:runVar012` | 81.8% |
| 🔴 `012.go:checkFuncForRepeatedConversions` | 80.0% |
| 🔴 `013.go:runVar013` | 77.8% |
| 🔴 `014.go:runVar014` | 84.2% |
| 🔴 `015.go:runVar015` | 82.4% |
| 🔴 `016.go:runVar016` | 86.7% |
| 🔴 `016.go:shouldUseArray` | 75.0% |
| 🔴 `016.go:getConstantSize` | 66.7% |
| 🟡 `017.go:runVar017` | 90.0% |
| 🟡 `017.go:collectTypesWithValueReceivers` | 92.9% |
| 🔴 `017.go:checkStructsWithMutex` | 81.2% |
| 🔴 `017.go:checkValueReceivers` | 87.5% |
| 🔴 `017.go:checkValueParams` | 83.3% |
| 🔴 `017.go:checkAssignments` | 83.3% |
| 🔴 `017.go:getMutexType` | 75.0% |
| 🔴 `017.go:hasMutex` | 75.0% |
| 🔴 `017.go:hasMutexInType` | 83.3% |
| 🔴 `017.go:getMutexTypeFromType` | 73.3% |
| 🔴 `017.go:isMutexCopy` | 85.7% |
| 🔴 `018.go:runVar018` | 86.7% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/shared` - 98.2%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `classify.go:ExtractReceiverTypeName` | 83.3% |
| 🔴 `types.go:hasSerializationTags` | 83.3% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/utils` - 94.4%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `ident.go:ExtractVarName` | 76.9% |
| 🔴 `slice.go:IsSmallConstantSize` | 85.7% |
| 🔴 `types.go:HasPositiveLength` | 36.4% |

### 🔴 `github.com/kodflow/ktn-linter/pkg/config` - 79.7%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `config.go:Get` | 38.5% |
| ⚫ `config.go:IsFileExcludedGlobally` | 0.0% |
| 🔴 `config.go:matchesAnyPattern` | 88.9% |
| 🔴 `config.go:matchDoubleStarPattern` | 66.7% |
| 🔴 `config.go:Merge` | 85.7% |
| 🔴 `loader.go:loadFromDefaultLocations` | 87.5% |
| 🔴 `loader.go:validateConfig` | 68.8% |
| 🔴 `loader.go:LoadAndSet` | 80.0% |
| 🔴 `loader.go:MustLoad` | 75.0% |
| 🔴 `loader.go:SaveToFile` | 66.7% |


---

*Généré le: 2025-12-10 07:54:23*

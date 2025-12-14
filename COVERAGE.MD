# Coverage Report

Rapport de couverture généré automatiquement.

**Légende:**
- 🔵 100% - Couverture parfaite
- 🟢 >90% - Excellente couverture
- 🟡 80-90% - Bonne couverture
- 🔴 <80% - Couverture insuffisante
- ⚫ 0% - Pas de tests

---

## Coverage par Package

| Icon | Package | Coverage |
|:----:|---------|----------|
| 🟢 | **TOTAL (Global)** | **94.4%** |
| 🔵 | `github.com/kodflow/ktn-linter/cmd/ktn-linter` | 100.0% |
| 🟢 | `github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd` | 97.7% |
| 🔵 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn` | 100.0% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment` | 98.5% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst` | 98.4% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc` | 95.3% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktninterface` | 95.2% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnreturn` | 98.4% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct` | 94.0% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest` | 88.3% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar` | 91.8% |
| 🔵 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper` | 100.0% |
| 🔵 | `github.com/kodflow/ktn-linter/pkg/analyzer/modernize` | 100.0% |
| 🔵 | `github.com/kodflow/ktn-linter/pkg/analyzer/shared` | 100.0% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/utils` | 98.5% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/config` | 98.1% |
| 🔵 | `github.com/kodflow/ktn-linter/pkg/formatter` | 100.0% |
| 🔵 | `github.com/kodflow/ktn-linter/pkg/severity` | 100.0% |

---

## Détail des packages incomplets

### 🟢 `github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd` - 97.7%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟢 `lint.go:runLint` | 92.9% |
| 🔴 `lint.go:selectSingleRule` | 57.1% |
| 🟡 `lint.go:analyzePackage` | 90.9% |

### 🟢 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment` - 98.5%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟢 `001.go:runComment001` | 96.0% |
| 🟢 `001.go:isDocComment` | 91.7% |
| 🟢 `002.go:runComment002` | 92.9% |
| 🟢 `005.go:runComment005` | 92.0% |

### 🟢 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst` - 98.4%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `002.go:extractConstTypeName` | 87.5% |
| 🟡 `002.go:checkScatteredConstBlocks` | 87.5% |

### 🟢 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc` - 95.3%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `001.go:isErrorType` | 84.6% |
| 🟢 `001.go:isBuiltinError` | 93.3% |
| 🟢 `002.go:analyzeContextParams` | 93.3% |
| 🟡 `002.go:isContextTypeByType` | 88.9% |
| 🔴 `002.go:isContextUnderlying` | 60.0% |
| 🟢 `003.go:runFunc003` | 95.0% |
| 🟢 `004.go:collectPrivateFunctions` | 92.3% |
| 🟢 `004.go:extractPrivateFuncInfo` | 91.7% |
| 🟡 `004.go:collectCalledFunctions` | 83.3% |
| 🟡 `004.go:collectCallbackUsages` | 81.8% |
| 🟡 `005.go:runFunc005` | 90.9% |
| 🟡 `005.go:countPureCodeLines` | 88.2% |
| 🟢 `007.go:runFunc007` | 94.7% |
| 🟢 `007.go:collectLazyLoadFields` | 92.9% |
| 🟡 `008.go:analyzeFunc008` | 88.9% |
| 🟡 `008.go:findImplementedInterface` | 85.7% |
| 🟡 `008.go:getReceiverType` | 83.3% |
| 🟡 `008.go:findInterfaceForMethod` | 88.9% |
| 🟡 `008.go:collectFunctionParams008` | 87.5% |
| 🟢 `008.go:collectIgnoredVariables008` | 93.3% |
| 🟢 `008.go:findParentAssignToBlank008` | 95.0% |
| 🟢 `009.go:checkMagicNumbers` | 92.3% |
| 🟢 `011.go:runFunc011` | 95.0% |

### 🟢 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktninterface` - 95.2%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `001.go:collectInterfaces` | 85.7% |
| 🟡 `001.go:extractTypesFromNode` | 88.9% |
| 🟡 `001.go:extractEmbeddedTypes` | 85.7% |
| 🔴 `001.go:extractCompositeType` | 66.7% |

### 🟢 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnreturn` - 98.4%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟢 `002.go:checkNilReturns` | 92.3% |

### 🟢 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct` - 94.0%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟢 `001.go:runStruct001` | 93.3% |
| 🟡 `001.go:checkStructs` | 80.0% |
| 🟢 `001.go:collectInterfaceChecksWithTypes` | 92.9% |
| 🔴 `001.go:extractInterfaceCheckWithTypes` | 62.5% |
| 🔴 `001.go:extractStructNameFromValue` | 75.0% |
| 🟡 `001.go:extractStructNameFromType` | 88.9% |
| 🟡 `001.go:extractStructNameFromAST` | 80.0% |
| 🟡 `001.go:extractFromCallExpr` | 80.0% |
| 🔴 `001.go:extractFromUnaryExpr` | 66.7% |
| 🔴 `001.go:extractFromCompositeLit` | 75.0% |
| 🟡 `001.go:hasMatchingInterfaceCheck` | 83.3% |
| 🟡 `001.go:interfaceCoversAllPublicMethods` | 90.0% |
| 🟡 `001.go:signaturesMatch` | 83.3% |
| 🟢 `001.go:collectInterfaces` | 94.4% |
| 🟢 `002.go:runStruct002` | 95.2% |
| 🟢 `003.go:isValidGetterToReport` | 92.3% |
| 🟢 `003.go:isSimpleGetter` | 93.3% |
| 🟢 `004.go:collectStructs` | 92.9% |
| 🟢 `005.go:runStruct005` | 94.4% |
| 🟢 `006.go:runStruct006` | 94.1% |
| 🟢 `007.go:collectMethodsDetailed` | 93.3% |
| 🟢 `007.go:checkGetterFieldMismatch` | 93.3% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest` - 88.3%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `001.go:runTest001` | 82.4% |
| 🟡 `002.go:runTest002` | 83.3% |
| 🟡 `003.go:runTest003` | 81.2% |
| 🟡 `004.go:collectFunctions` | 87.5% |
| 🟡 `004.go:runTest004` | 88.9% |
| 🟢 `004.go:reportMissingTest` | 91.7% |
| 🟡 `004.go:findPackageDir` | 83.3% |
| 🟡 `004.go:collectExternalTestFunctions` | 83.3% |
| 🟡 `004.go:collectExternalTestedFunctions` | 86.7% |
| 🟡 `005.go:runTest005` | 82.4% |
| 🔴 `006.go:runTest006` | 77.8% |
| 🟢 `006.go:collectFiles006` | 92.3% |
| 🔴 `006.go:extractBaseName006` | 77.8% |
| 🟡 `007.go:runTest007` | 82.6% |
| 🔴 `008.go:runTest008` | 73.7% |
| 🟡 `008.go:classifyFunction` | 85.7% |
| 🟡 `008.go:extractReceiverTypeString` | 85.7% |
| 🔴 `008.go:checkVariables` | 78.6% |
| 🟡 `008.go:checkTypes` | 80.0% |
| 🟡 `008.go:checkConsts` | 81.8% |
| 🟡 `009.go:runTest009` | 85.7% |
| 🟡 `009.go:collectPublicFunctions` | 90.9% |
| 🟡 `009.go:addPublicFunction` | 83.3% |
| 🟡 `009.go:checkInternalTestsForPublicFunctions` | 84.6% |
| 🟡 `009.go:checkAndReportPublicFunctionTest` | 81.8% |
| 🟡 `010.go:runTest010` | 85.7% |
| 🟡 `010.go:collectPrivateFunctions` | 90.9% |
| 🟡 `010.go:addPrivateFunction` | 83.3% |
| 🟡 `010.go:checkExternalTestsForPrivateFunctions` | 84.6% |
| 🟡 `010.go:checkAndReportPrivateFunctionTest` | 84.6% |
| 🟡 `011.go:runTest011` | 87.5% |
| 🟡 `011.go:validateTestFiles011` | 88.9% |
| 🔴 `012.go:runTest012` | 73.9% |
| 🟡 `012.go:isAssertLibraryCall` | 87.5% |
| 🔴 `013.go:runTest013` | 69.6% |
| 🟡 `013.go:collectFuncSignatures` | 85.7% |
| 🟡 `013.go:addFuncSignature` | 83.3% |
| 🟡 `013.go:collectExternalSourceSignatures` | 81.8% |
| 🟡 `013.go:scanSourceFile` | 84.2% |
| 🟡 `013.go:functionReturnsError` | 87.5% |
| 🔴 `013.go:analyzeTestFunction` | 76.9% |
| 🔴 `013.go:checkErrorInKeyValue` | 75.0% |

### 🟢 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar` - 91.8%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟢 `001.go:runVar001` | 95.0% |
| 🟡 `002.go:checkVarSpec` | 83.3% |
| 🟡 `003.go:runVar003` | 84.6% |
| 🟢 `004.go:collectAppendVariables` | 91.7% |
| 🔴 `004.go:checkMakeCalls` | 77.8% |
| 🟡 `004.go:checkMakeCall` | 88.9% |
| 🟡 `004.go:checkEmptySliceLiterals` | 88.2% |
| 🔴 `004.go:checkCompositeLit` | 75.0% |
| 🟡 `005.go:checkMakeCallVar008` | 81.8% |
| 🟡 `005.go:runVar005` | 90.9% |
| 🟡 `006.go:reportMissingGrow` | 83.3% |
| 🟡 `006.go:extractAssignTypeString` | 87.5% |
| 🟡 `007.go:runVar007` | 90.0% |
| 🟢 `007.go:checkStringConcatInLoop` | 91.7% |
| 🔴 `007.go:isStringConcatenation` | 75.0% |
| 🟡 `008.go:runVar008` | 86.7% |
| 🔴 `008.go:checkStmtForAlloc` | 66.7% |
| 🟡 `008.go:checkDeclForAlloc` | 90.9% |
| 🟡 `009.go:runVar009` | 85.7% |
| 🔴 `009.go:checkDeclForLargeStruct` | 75.0% |
| 🔴 `009.go:checkTypeForLargeStruct` | 76.9% |
| 🟡 `009.go:isExternalType` | 85.7% |
| 🟢 `010.go:runVar010` | 91.7% |
| 🟡 `010.go:checkLoopBodyVar015` | 83.3% |
| 🟡 `010.go:checkMakeCallForByteSlice` | 80.0% |
| 🟡 `011.go:runVar011` | 90.0% |
| 🟡 `011.go:checkShortVarDecl` | 88.9% |
| 🟡 `011.go:isShadowing` | 81.8% |
| 🟡 `012.go:runVar012` | 90.9% |
| 🟡 `012.go:checkFuncForRepeatedConversions` | 80.0% |
| 🟡 `013.go:runVar013` | 88.9% |
| 🟢 `014.go:runVar014` | 94.7% |
| 🟡 `015.go:runVar015` | 88.2% |
| 🟢 `016.go:runVar016` | 93.3% |
| 🟡 `016.go:shouldUseArray` | 87.5% |
| 🟡 `016.go:getConstantSize` | 83.3% |
| 🟢 `017.go:collectTypesWithValueReceivers` | 92.9% |
| 🟡 `017.go:checkStructsWithMutex` | 81.2% |
| 🟡 `017.go:checkValueReceivers` | 87.5% |
| 🟡 `017.go:checkValueParams` | 83.3% |
| 🟡 `017.go:checkAssignments` | 83.3% |
| 🟡 `017.go:hasMutexInType` | 83.3% |
| 🟡 `017.go:getMutexTypeFromType` | 86.7% |
| 🟡 `017.go:isMutexCopy` | 85.7% |

### 🟢 `github.com/kodflow/ktn-linter/pkg/analyzer/utils` - 98.5%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `ident.go:ExtractVarName` | 76.9% |

### 🟢 `github.com/kodflow/ktn-linter/pkg/config` - 98.1%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟢 `config.go:matchesAnyPattern` | 94.4% |
| 🟢 `loader.go:loadFromDefaultLocations` | 93.8% |
| 🟡 `loader.go:SaveToFile` | 83.3% |


---

*Généré le: 2025-12-14 11:15:17*

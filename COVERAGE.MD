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
| 🟡 | **TOTAL (Global)** | **88.3%** |
| 🔵 | `github.com/kodflow/ktn-linter/cmd/ktn-linter` | 100.0% |
| 🟡 | `github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd` | 86.9% |
| 🔵 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn` | 100.0% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnapi` | 93.9% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment` | 98.1% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst` | 94.6% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc` | 96.4% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktngeneric` | 96.9% |
| 🔵 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktninterface` | 100.0% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnreturn` | 98.4% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct` | 97.5% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest` | 90.8% |
| 🔴 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar` | 76.9% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper` | 98.0% |
| 🔵 | `github.com/kodflow/ktn-linter/pkg/analyzer/modernize` | 100.0% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/shared` | 99.2% |
| 🔵 | `github.com/kodflow/ktn-linter/pkg/analyzer/utils` | 100.0% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/config` | 98.1% |
| 🔵 | `github.com/kodflow/ktn-linter/pkg/formatter` | 100.0% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/messages` | 99.1% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/orchestrator` | 82.9% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/prompt` | 97.0% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/rules` | 96.0% |
| 🔵 | `github.com/kodflow/ktn-linter/pkg/severity` | 100.0% |
| 🔴 | `github.com/kodflow/ktn-linter/pkg/updater` | 53.1% |

---

## Détail des packages incomplets

### 🟡 `github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd` - 86.9%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟢 `lint.go:runLint` | 91.7% |
| 🟡 `lint.go:loadConfiguration` | 90.0% |
| 🟡 `lint.go:formatAndDisplay` | 88.9% |
| 🔴 `prompt.go:runPrompt` | 75.0% |
| 🔴 `rules.go:displayRuleDetails` | 66.7% |
| 🔴 `rules.go:DisplayRuleDetails` | 54.5% |
| 🔴 `rules.go:DisplayRuleDetails` | 44.4% |
| 🔴 `rules.go:encodeJSON` | 60.0% |
| 🟡 `upgrade.go:runUpgrade` | 83.3% |
| 🔴 `upgrade.go:handleCheckOnly` | 60.0% |
| ⚫ `upgrade.go:handleUpgrade` | 0.0% |

### 🟢 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnapi` - 93.9%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟢 `001.go:runAPI001` | 91.7% |
| 🟡 `001.go:analyzeFunction` | 90.0% |
| 🟡 `001.go:collectParams` | 84.6% |
| 🟢 `001.go:getExternalConcreteType` | 94.4% |
| 🟡 `001.go:isAllowedType` | 87.5% |
| 🟢 `001.go:scanBodyForMethodCalls` | 92.9% |
| 🟢 `001.go:scanBodyForFieldAccess` | 95.2% |
| 🟡 `001.go:getMethodSignature` | 87.5% |
| 🟡 `001.go:getBaseIdent` | 85.7% |
| 🟡 `001.go:reportDiagnostics` | 90.9% |

### 🟢 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment` - 98.1%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟢 `001.go:runComment001` | 96.2% |
| 🟢 `001.go:isDocComment` | 91.7% |
| 🟢 `002.go:runComment002` | 93.3% |
| 🟢 `005.go:runComment005` | 92.3% |
| 🟢 `007.go:hasCommentBefore` | 92.9% |
| 🟢 `007.go:hasInlineComment` | 93.3% |

### 🟢 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst` - 94.6%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟢 `001.go:runConst001` | 95.0% |
| 🟡 `002.go:extractConstTypeName` | 87.5% |
| 🟡 `002.go:checkConstOrder` | 86.7% |
| 🟡 `004.go:runConst004` | 85.0% |
| 🟡 `005.go:runConst005` | 85.0% |
| 🟡 `006.go:runConst006` | 85.0% |

### 🟢 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc` - 96.4%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `001.go:isErrorType` | 84.6% |
| 🟢 `001.go:isBuiltinError` | 93.3% |
| 🟡 `002.go:isContextTypeByType` | 88.9% |
| 🔴 `002.go:isContextUnderlying` | 60.0% |
| 🟢 `003.go:runFunc003` | 95.2% |
| 🟢 `004.go:collectPrivateFunctions` | 92.3% |
| 🟡 `004.go:collectCalledFunctions` | 83.3% |
| 🟡 `004.go:collectCallbackUsages` | 81.8% |
| 🟢 `005.go:runFunc005` | 91.3% |
| 🟡 `005.go:countStmtComplexity` | 88.9% |
| 🟢 `007.go:runFunc007` | 94.7% |
| 🟢 `007.go:collectLazyLoadFields` | 92.9% |
| 🟡 `008.go:findInterfaceForMethod` | 88.9% |
| 🟢 `008.go:collectIgnoredVariables008` | 93.3% |
| 🟢 `008.go:findParentAssignToBlank008` | 95.0% |
| 🟢 `009.go:checkMagicNumbers` | 92.9% |
| 🟢 `011.go:runFunc011` | 95.2% |

### 🟢 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktngeneric` - 96.9%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟢 `001.go:checkEqualityUsage` | 92.9% |
| 🔴 `002.go:isSingleInterfaceConstraint` | 75.0% |
| 🟡 `002.go:reportUnnecessaryGeneric` | 83.3% |
| 🟡 `003.go:checkDeprecatedConstraintsImport` | 90.0% |
| 🟡 `005.go:reportShadowing` | 80.0% |
| 🟢 `006.go:checkOrderedUsage` | 93.8% |
| 🟡 `006.go:reportIfUsesAnyTypeParamOrdered` | 84.6% |

### 🟢 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnreturn` - 98.4%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟢 `001.go:checkNilReturns` | 92.9% |

### 🟢 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct` - 97.5%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟢 `001.go:collectMethodsDetailed` | 93.3% |
| 🟢 `001.go:checkGetterFieldMismatch` | 93.8% |
| 🟢 `002.go:runStruct002` | 95.5% |
| 🟢 `003.go:isValidGetterToReport` | 92.3% |
| 🟢 `003.go:isSimpleGetter` | 93.3% |
| 🟢 `004.go:collectStructs` | 92.9% |
| 🟢 `005.go:runStruct005` | 94.4% |
| 🟢 `006.go:runStruct006` | 94.1% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest` - 90.8%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `001.go:runTest001` | 82.4% |
| 🟡 `002.go:runTest002` | 81.2% |
| 🟡 `003.go:collectFunctions` | 87.5% |
| 🟡 `003.go:runTest003` | 88.9% |
| 🟡 `003.go:reportMissingTest` | 88.9% |
| 🟡 `003.go:findPackageDir` | 83.3% |
| 🟡 `003.go:collectExternalTestFunctions` | 83.3% |
| 🟡 `003.go:collectExternalTestedFunctions` | 86.7% |
| 🟡 `004.go:runTest004` | 83.3% |
| 🟡 `005.go:runTest005` | 83.3% |
| 🔴 `006.go:runTest006` | 73.7% |
| 🟢 `006.go:checkVariables` | 92.9% |
| 🟡 `006.go:checkTypes` | 90.0% |
| 🟡 `006.go:checkConsts` | 90.9% |
| 🟡 `007.go:runTest007` | 85.7% |
| 🟡 `007.go:collectPublicFunctions` | 90.9% |
| 🟡 `007.go:addPublicFunction` | 83.3% |
| 🟡 `007.go:checkInternalTestsForPublicFunctions` | 84.6% |
| 🟡 `007.go:checkAndReportPublicFunctionTest` | 81.8% |
| 🟡 `008.go:runTest008` | 85.7% |
| 🟡 `008.go:collectPrivateFunctions` | 90.9% |
| 🟡 `008.go:addPrivateFunction` | 83.3% |
| 🟡 `008.go:checkExternalTestsForPrivateFunctions` | 84.6% |
| 🟡 `008.go:checkAndReportPrivateFunctionTest` | 84.6% |
| 🟡 `009.go:runTest009` | 87.5% |
| 🟡 `009.go:validateTestFiles011` | 88.9% |
| 🔴 `010.go:runTest010` | 75.0% |
| 🔴 `011.go:runTest011` | 69.6% |
| 🟡 `011.go:collectFuncSignatures` | 85.7% |
| 🟡 `011.go:collectExternalSourceSignatures` | 90.9% |
| 🟢 `011.go:scanSourceFile` | 94.7% |
| 🟢 `011.go:analyzeTestFunction` | 92.9% |

### 🔴 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar` - 76.9%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `001.go:runVar001` | 88.5% |
| 🔴 `001.go:checkVarSpec` | 70.0% |
| 🟡 `002.go:runVar002` | 80.0% |
| 🔴 `003.go:runVar003` | 72.7% |
| 🔴 `003.go:checkVar003Names` | 70.0% |
| 🔴 `004.go:runVar004` | 75.0% |
| 🟡 `004.go:checkVar004PackageLevel` | 83.3% |
| 🔴 `004.go:checkVar004LocalVars` | 75.0% |
| 🔴 `004.go:checkVar004Node` | 66.7% |
| 🟡 `004.go:checkVar004AssignStmt` | 85.7% |
| ⚫ `004.go:checkVar004DeclStmt` | 0.0% |
| 🔴 `004.go:checkVar004Name` | 71.4% |
| 🔴 `005.go:runVar005` | 75.0% |
| 🔴 `005.go:checkVar005PackageLevel` | 75.0% |
| 🔴 `005.go:checkVar005LocalVars` | 16.7% |
| ⚫ `005.go:checkVar005Node` | 0.0% |
| ⚫ `005.go:checkVar005AssignStmt` | 0.0% |
| ⚫ `005.go:checkVar005RangeStmt` | 0.0% |
| ⚫ `005.go:checkVar005DeclStmt` | 0.0% |
| 🔴 `005.go:checkVar005Name` | 70.0% |
| 🔴 `006.go:runVar006` | 73.3% |
| 🔴 `007.go:runVar007` | 75.0% |
| 🟡 `007.go:checkControlFlowStmt` | 83.3% |
| 🔴 `007.go:checkCaseClause` | 50.0% |
| 🔴 `007.go:checkCommClause` | 50.0% |
| 🔴 `007.go:reportVarErrors` | 66.7% |
| 🟡 `008.go:runVar008` | 84.6% |
| 🟢 `008.go:collectAppendVariables` | 91.7% |
| 🟡 `008.go:checkMakeCalls` | 88.9% |
| 🔴 `008.go:checkMakeCall` | 76.9% |
| 🟡 `008.go:checkEmptySliceLiterals` | 88.2% |
| 🟡 `008.go:getAppendVariableName008` | 87.5% |
| ⚫ `008.go:reportVar008Error` | 0.0% |
| 🟡 `008.go:checkCompositeLit` | 80.0% |
| 🔴 `009.go:checkMakeCallVar008` | 73.3% |
| 🔴 `009.go:runVar009` | 77.8% |
| 🟡 `010.go:runVar010` | 86.7% |
| 🔴 `010.go:reportMissingGrow` | 70.0% |
| 🟡 `010.go:extractAssignTypeString` | 87.5% |
| 🟡 `011.go:runVar011` | 82.4% |
| 🟡 `011.go:checkStringConcatInLoop` | 81.2% |
| 🟡 `011.go:isStringConcatenation` | 87.5% |
| 🟡 `012.go:runVar012` | 80.0% |
| 🔴 `012.go:checkAssignForAlloc` | 71.4% |
| 🟡 `012.go:checkDeclForAlloc` | 86.7% |
| 🟡 `013.go:runVar013` | 84.6% |
| 🟡 `013.go:checkFuncParams009` | 90.0% |
| 🔴 `013.go:checkParamType009` | 75.0% |
| 🔴 `013.go:getStructSize009` | 73.7% |
| 🟡 `013.go:isExternalType009` | 85.7% |
| 🟡 `014.go:runVar014` | 82.4% |
| 🟡 `014.go:checkLoopBodyVar015` | 83.3% |
| 🔴 `014.go:checkMakeCallForByteSlice` | 66.7% |
| 🟡 `015.go:runVar015` | 81.2% |
| 🟡 `015.go:checkFuncForRepeatedConversions` | 80.0% |
| 🟡 `015.go:checkLoopsForStringConversion` | 81.8% |
| 🟡 `015.go:checkMultipleConversions` | 90.0% |
| 🟡 `016.go:runVar016` | 88.9% |
| 🔴 `016.go:checkVarGrouping` | 75.0% |
| 🔴 `017.go:runVar017` | 75.0% |
| 🟡 `018.go:runVar018` | 82.6% |
| 🟡 `018.go:shouldUseArray` | 90.0% |
| 🟡 `018.go:getConstantSize` | 83.3% |
| 🔴 `018.go:isTotalSizeSmall` | 71.4% |
| 🔴 `018.go:reportArraySuggestion` | 60.0% |
| 🟡 `019.go:runVar019` | 82.4% |
| 🟢 `019.go:collectTypesWithValueReceivers` | 92.9% |
| 🔴 `019.go:checkStructsWithMutex` | 75.0% |
| 🔴 `019.go:checkValueReceivers` | 78.9% |
| 🔴 `019.go:checkValueParams` | 75.0% |
| 🔴 `019.go:checkAssignments` | 75.0% |
| 🟡 `019.go:hasMutexInType` | 84.6% |
| 🟡 `019.go:getMutexTypeFromType` | 87.5% |
| 🟡 `019.go:isMutexCopy` | 88.9% |
| 🔴 `020.go:runVar020` | 76.5% |
| 🔴 `020.go:checkEmptySliceLiteral` | 70.0% |
| 🔴 `020.go:checkMakeSliceZero` | 68.8% |
| 🔴 `020.go:isMakeCall` | 75.0% |
| 🔴 `020.go:isZeroLiteral` | 75.0% |
| 🔴 `020.go:formatSliceType` | 75.0% |
| 🔴 `020.go:formatElementType` | 25.0% |
| ⚫ `020.go:formatSelectorType` | 0.0% |
| ⚫ `020.go:formatMapType` | 0.0% |
| 🔴 `021.go:runVar021` | 75.0% |
| 🔴 `021.go:collectReceivers` | 75.0% |
| 🔴 `021.go:extractReceiverType` | 71.4% |
| 🟡 `021.go:checkReceiverConsistency` | 85.7% |
| 🔴 `021.go:reportInconsistency` | 75.0% |
| 🟡 `022.go:runVar022` | 86.7% |
| 🟡 `022.go:checkFuncDecls` | 81.8% |
| 🟡 `022.go:checkVarDecls` | 80.0% |
| 🟡 `022.go:checkStructFields` | 83.3% |
| 🟡 `022.go:checkPointerToInterface` | 83.3% |
| 🔴 `023.go:runVar023` | 71.4% |
| 🟡 `023.go:checkMathRandUsage` | 85.7% |
| 🟡 `023.go:checkFuncForMathRand` | 81.2% |
| 🔴 `023.go:checkAssignForMathRand` | 40.0% |
| ⚫ `023.go:checkLocalVarDecl` | 0.0% |
| ⚫ `023.go:processLocalVarSpec` | 0.0% |
| 🔴 `023.go:hasSecurityVarName` | 75.0% |
| 🟡 `023.go:checkValuesForMathRand` | 83.3% |
| 🟡 `023.go:processGenDeclSpec023` | 85.7% |
| 🔴 `023.go:isMathRandCall` | 71.4% |
| 🔴 `023.go:reportMathRandInSecurity` | 60.0% |
| 🔴 `024.go:runVar024` | 72.2% |
| 🔴 `024.go:checkEmptyInterface` | 71.4% |
| 🔴 `024.go:isEmptyInterface` | 66.7% |
| 🔴 `025.go:runVar025` | 70.0% |
| 🔴 `025.go:checkMapDeletePattern` | 73.3% |
| 🔴 `025.go:isDeleteCallWithKeyAndMap` | 69.2% |
| 🔴 `025.go:checkSliceZeroPattern` | 66.7% |
| 🟡 `025.go:checkIndexAssignZero` | 80.0% |
| 🟡 `025.go:isMatchingSliceIndex` | 85.7% |
| 🔴 `025.go:isZeroValue` | 50.0% |
| 🔴 `025.go:isBasicLitZero025` | 22.2% |
| 🔴 `025.go:isIdentZero025` | 60.0% |
| ⚫ `025.go:isCompositeLitZero025` | 0.0% |
| ⚫ `025.go:isZeroConversion` | 0.0% |
| 🔴 `025.go:getRangeCollectionIdent` | 50.0% |
| 🔴 `025.go:reportClearPattern` | 60.0% |
| 🔴 `026.go:runVar026` | 76.5% |
| 🔴 `026.go:checkMathMinMax` | 68.8% |
| 🔴 `026.go:checkIfMinMaxPattern` | 36.4% |
| 🔴 `026.go:isMinMaxCondition` | 75.0% |
| 🔴 `026.go:hasReturnInBody` | 75.0% |
| 🔴 `026.go:hasMatchingReturn` | 66.7% |
| ⚫ `026.go:hasReturnInElse` | 0.0% |
| 🔴 `027.go:runVar027` | 69.6% |
| 🔴 `027.go:checkInitIsZero` | 77.8% |
| 🔴 `027.go:extractAssignFromInit027` | 62.5% |
| 🔴 `027.go:extractVarNameFromAssign027` | 75.0% |
| 🔴 `027.go:validateInitZero027` | 75.0% |
| 🔴 `027.go:checkPostIsIncrement` | 77.8% |
| 🔴 `027.go:checkConditionIsLessThan` | 77.8% |
| 🔴 `028.go:runVar028` | 72.2% |
| 🔴 `028.go:checkLoopVarCopyPattern` | 71.4% |
| 🟡 `028.go:checkShortVarDecl` | 85.7% |
| 🔴 `028.go:checkAssignmentPair` | 76.9% |
| 🔴 `028.go:reportLoopVarCopy` | 60.0% |
| 🔴 `029.go:runVar029` | 72.2% |
| 🔴 `029.go:extractSliceFromCapLenCondition` | 66.7% |
| 🟡 `029.go:extractSliceFromCapMinusLen` | 87.5% |
| 🔴 `029.go:extractSliceFromCapCall` | 66.7% |
| 🔴 `029.go:extractSliceFromLenCall` | 66.7% |
| 🔴 `029.go:extractSliceName` | 75.0% |
| 🟡 `029.go:hasGrowPatternInBody` | 83.3% |
| 🟡 `029.go:analyzeAssignForGrow` | 83.3% |
| 🟡 `029.go:isMakeSliceCall` | 80.0% |
| 🔴 `029.go:checkSliceReassign` | 75.0% |
| 🔴 `029.go:analyzeCopyForGrow` | 72.7% |
| 🔴 `029.go:reportGrowPattern` | 60.0% |
| 🔴 `030.go:runVar030` | 75.0% |
| 🔴 `030.go:checkAppendNilPattern` | 71.4% |
| 🟡 `030.go:isAppendNilPattern` | 87.5% |
| 🟡 `030.go:isNilSliceConversion` | 80.0% |
| 🔴 `030.go:checkMakeCopyPattern` | 75.0% |
| 🔴 `030.go:extractMakeAssign030` | 77.8% |
| 🟡 `030.go:validateMakeCall030` | 87.5% |
| 🔴 `030.go:extractLenSource` | 66.7% |
| 🔴 `030.go:exprToString` | 40.0% |
| 🔴 `030.go:checkCopyCall` | 75.0% |
| 🔴 `030.go:extractCopyCall030` | 70.0% |
| 🔴 `030.go:findMatchingMake030` | 66.7% |
| 🔴 `030.go:reportMakeCopy030` | 60.0% |
| 🔴 `031.go:runVar031` | 71.4% |
| 🟡 `031.go:analyzeBlockForMapClone` | 85.7% |
| 🟢 `031.go:extractMakeMapAssign` | 92.9% |
| 🟡 `031.go:isMakeMapCallExpr` | 85.7% |
| 🟡 `031.go:isSimpleMapClone` | 85.7% |
| 🔴 `031.go:extractRangeKeyValue031` | 71.4% |
| 🟡 `031.go:extractRangeBodyAssign031` | 87.5% |
| 🟡 `031.go:validateCloneAssignment031` | 88.9% |
| 🔴 `031.go:validateIndexExpr031` | 71.4% |
| 🔴 `033.go:runVar033` | 72.7% |
| 🔴 `033.go:checkCmpOrPattern` | 66.7% |
| 🟡 `033.go:extractVarFromCondition033` | 85.7% |
| 🔴 `033.go:validateIfBodyReturns033` | 75.0% |
| 🔴 `033.go:extractVarAndZeroFor033` | 40.0% |
| 🔴 `033.go:isSimpleZeroValueFor033` | 75.0% |
| 🔴 `033.go:returnsVariableFor033` | 77.8% |
| 🔴 `033.go:reportCmpOrPattern` | 60.0% |
| 🔴 `034.go:runVar034` | 75.0% |
| 🟡 `034.go:extractWaitGroupAdd1` | 88.9% |
| 🔴 `034.go:extractWaitGroupMethodCall` | 62.5% |
| 🔴 `034.go:isWaitGroupReceiver` | 75.0% |
| 🔴 `034.go:isWaitGroupType` | 66.7% |
| 🟡 `034.go:isGoWithDeferDone` | 85.7% |
| 🟡 `034.go:hasDeferDoneFirst` | 85.7% |
| 🔴 `034.go:reportVar034` | 60.0% |
| 🔴 `035.go:runVar035` | 72.7% |
| 🔴 `035.go:checkContainsPattern` | 66.7% |
| 🟡 `035.go:analyzeBodyForContainsPattern` | 90.9% |
| 🔴 `035.go:matchesContainsPattern` | 75.0% |
| 🔴 `035.go:isBlankOrNil` | 75.0% |
| 🔴 `035.go:matchesIfReturnTruePattern` | 73.3% |
| 🔴 `035.go:containsRangeValue` | 50.0% |
| 🟡 `035.go:returnsLiteralBool` | 87.5% |
| 🔴 `035.go:reportContainsPattern` | 66.7% |
| 🔴 `036.go:runVar036` | 72.2% |
| 🟢 `036.go:analyzeFunctionForIndexPattern` | 93.8% |
| 🔴 `036.go:functionReturnsInt` | 77.8% |
| 🔴 `036.go:isReturnMinusOne` | 75.0% |
| 🔴 `036.go:extractRangeVariables036` | 77.8% |
| 🟡 `036.go:findIndexPatternInBody036` | 83.3% |
| 🔴 `036.go:checkIfEqualReturnIndex` | 77.8% |
| 🔴 `036.go:ifBodyReturnsIndex` | 66.7% |
| 🔴 `036.go:reportIndexPattern` | 60.0% |
| 🔴 `037.go:runVar037` | 69.2% |
| 🔴 `037.go:checkMapCollectionPatterns` | 76.9% |
| 🟡 `037.go:isRangingOverMap` | 80.0% |
| 🟡 `037.go:extractAppendCall037` | 81.8% |
| 🔴 `037.go:isVar037AppendCall` | 75.0% |
| 🟡 `037.go:isKeyCollection` | 80.0% |
| 🟡 `037.go:isValueCollection` | 80.0% |
| 🟡 `037.go:isSameIdent` | 85.7% |
| 🔴 `037.go:reportMapCollectionPattern` | 60.0% |

### 🟢 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper` - 98.0%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `helper.go:createPass` | 83.3% |
| 🟡 `helper.go:createPassForPackage` | 88.9% |

### 🟢 `github.com/kodflow/ktn-linter/pkg/analyzer/shared` - 99.2%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `methods.go:ExtractReceiverName` | 87.5% |

### 🟢 `github.com/kodflow/ktn-linter/pkg/config` - 98.1%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟢 `config.go:matchesAnyPattern` | 94.4% |
| 🟢 `loader.go:loadFromDefaultLocations` | 93.8% |
| 🟡 `loader.go:SaveToFile` | 83.3% |

### 🟢 `github.com/kodflow/ktn-linter/pkg/messages` - 99.1%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `messages.go:Format` | 83.3% |

### 🟡 `github.com/kodflow/ktn-linter/pkg/orchestrator` - 82.9%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `discovery.go:FindModules` | 90.9% |
| 🟡 `discovery.go:findInPath` | 88.9% |
| 🔴 `discovery.go:findModuleForFile` | 71.4% |
| 🟡 `discovery.go:searchDirectory` | 87.5% |
| 🟡 `loader.go:LoadFromDir` | 85.7% |
| 🟡 `orchestrator.go:Run` | 90.0% |
| ⚫ `orchestrator.go:DiscoverModules` | 0.0% |
| ⚫ `orchestrator.go:LoadPackagesFromDir` | 0.0% |
| ⚫ `orchestrator.go:RunMultiModule` | 0.0% |
| ⚫ `orchestrator.go:runSingleModule` | 0.0% |
| 🟢 `runner.go:Run` | 94.4% |
| 🟢 `runner.go:analyzePackageParallel` | 91.7% |
| 🟡 `runner.go:runAnalyzerGroup` | 83.3% |
| 🟡 `runner.go:shouldIncludeFile` | 90.0% |

### 🟢 `github.com/kodflow/ktn-linter/pkg/prompt` - 97.0%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟢 `formatter.go:Format` | 92.9% |
| 🟡 `formatter.go:writePhase` | 90.9% |
| 🟡 `formatter.go:writeRule` | 87.5% |
| 🟡 `generator.go:collectViolations` | 87.5% |
| 🟡 `phases.go:ClassifyRule` | 88.9% |

### 🟢 `github.com/kodflow/ktn-linter/pkg/rules` - 96.0%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `examples.go:findProjectRoot` | 60.0% |

### 🔴 `github.com/kodflow/ktn-linter/pkg/updater` - 53.1%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🔴 `updater.go:CheckForUpdate` | 28.6% |
| 🔴 `updater.go:Upgrade` | 33.3% |
| 🟡 `updater.go:getLatestVersion` | 90.9% |
| 🔴 `updater.go:downloadAndReplace` | 21.9% |
| 🔴 `updater.go:getBinaryName` | 75.0% |


---

*Généré le: 2026-01-02 09:47:41*

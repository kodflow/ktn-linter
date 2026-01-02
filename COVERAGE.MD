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
| 🟢 | **TOTAL (Global)** | **94.1%** |
| 🔵 | `github.com/kodflow/ktn-linter/cmd/ktn-linter` | 100.0% |
| 🟡 | `github.com/kodflow/ktn-linter/cmd/ktn-linter/cmd` | 86.9% |
| 🔵 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn` | 100.0% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnapi` | 93.9% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktncomment` | 98.1% |
| 🔵 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnconst` | 100.0% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnfunc` | 96.4% |
| 🔵 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktngeneric` | 100.0% |
| 🔵 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktninterface` | 100.0% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnreturn` | 98.4% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnstruct` | 97.5% |
| 🟡 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktntest` | 90.8% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar` | 99.2% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/testhelper` | 98.0% |
| 🔵 | `github.com/kodflow/ktn-linter/pkg/analyzer/modernize` | 100.0% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/analyzer/shared` | 99.2% |
| 🔵 | `github.com/kodflow/ktn-linter/pkg/analyzer/utils` | 100.0% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/config` | 98.1% |
| 🔵 | `github.com/kodflow/ktn-linter/pkg/formatter` | 100.0% |
| 🟢 | `github.com/kodflow/ktn-linter/pkg/messages` | 94.2% |
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

### 🟢 `github.com/kodflow/ktn-linter/pkg/analyzer/ktn/ktnvar` - 99.2%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟢 `001.go:runVar001` | 91.7% |
| 🟡 `002.go:runVar002` | 88.0% |
| 🟡 `003.go:runVar003` | 85.0% |
| 🟡 `004.go:runVar004` | 90.0% |
| 🟡 `004.go:checkVar004PackageLevel` | 88.9% |
| 🟢 `004.go:checkVar004LocalVars` | 91.7% |
| 🟢 `004.go:checkVar004Name` | 93.3% |
| 🟡 `005.go:runVar005` | 90.0% |
| 🟢 `005.go:checkVar005PackageLevel` | 93.8% |
| 🟢 `005.go:checkVar005LocalVars` | 91.7% |
| 🟢 `006.go:runVar006` | 92.0% |
| 🟡 `007.go:runVar007` | 88.9% |
| 🟡 `008.go:runVar008` | 90.9% |
| 🟢 `008.go:collectAppendVariables` | 91.7% |
| 🟡 `008.go:checkMakeCalls` | 88.9% |
| 🟡 `008.go:checkEmptySliceLiterals` | 88.2% |
| 🟡 `009.go:runVar009` | 81.2% |
| 🟢 `010.go:runVar010` | 92.3% |
| 🟡 `010.go:reportMissingGrow` | 85.7% |
| 🟡 `011.go:runVar011` | 84.6% |
| 🟡 `012.go:runVar012` | 83.3% |
| 🟡 `013.go:runVar013` | 87.5% |
| 🟡 `013.go:checkParamType009` | 88.9% |
| 🟢 `013.go:getStructSize009` | 94.7% |
| 🟡 `014.go:runVar014` | 86.7% |
| 🟡 `015.go:runVar015` | 85.7% |
| 🟡 `016.go:runVar016` | 88.9% |
| 🔴 `016.go:checkVarGrouping` | 75.0% |
| 🟡 `017.go:runVar017` | 82.6% |
| 🟡 `018.go:runVar018` | 85.7% |
| 🟡 `018.go:getConstantSize` | 83.3% |
| 🟢 `018.go:isTotalSizeSmall` | 92.9% |
| 🟡 `019.go:runVar019` | 86.7% |
| 🟡 `019.go:checkStructsWithMutex` | 82.4% |
| 🟢 `019.go:checkValueReceivers` | 93.8% |
| 🟡 `019.go:checkValueParams` | 84.6% |
| 🟡 `019.go:checkAssignments` | 84.6% |
| 🟢 `019.go:getMutexTypeFromType` | 93.8% |
| 🟡 `020.go:runVar020` | 88.2% |
| 🟡 `020.go:checkEmptySliceLiteral` | 80.0% |
| 🟡 `020.go:checkMakeSliceZero` | 87.5% |
| 🟡 `021.go:runVar021` | 83.3% |
| 🟡 `021.go:collectReceivers` | 87.5% |
| 🔴 `021.go:reportInconsistency` | 75.0% |
| 🟡 `022.go:runVar022` | 86.7% |
| 🟡 `022.go:checkFuncDecls` | 81.8% |
| 🟡 `022.go:checkVarDecls` | 80.0% |
| 🟡 `022.go:checkStructFields` | 83.3% |
| 🟡 `022.go:checkPointerToInterface` | 83.3% |
| 🔴 `023.go:runVar023` | 75.0% |
| 🟡 `023.go:checkMathRandUsage` | 85.7% |
| 🟡 `023.go:checkAssignForMathRand` | 80.0% |
| 🔴 `024.go:runVar024` | 75.0% |
| 🔴 `025.go:runVar025` | 72.2% |
| 🔴 `025.go:isZeroValue` | 66.7% |
| 🔴 `025.go:isZeroConversion` | 38.5% |
| 🟡 `026.go:runVar026` | 80.0% |
| 🔴 `026.go:checkIfMinMaxPattern` | 75.0% |
| 🔴 `027.go:runVar027` | 77.8% |
| 🔴 `028.go:runVar028` | 75.0% |
| 🔴 `029.go:runVar029` | 75.0% |
| 🟡 `030.go:runVar030` | 80.0% |
| 🟡 `030.go:checkAppendNilPattern` | 81.8% |
| 🔴 `030.go:checkMakeCopyPattern` | 75.0% |
| 🔴 `030.go:checkCopyCall` | 75.0% |
| 🔴 `031.go:runVar031` | 73.7% |
| 🔴 `033.go:runVar033` | 77.8% |
| 🔴 `033.go:checkCmpOrPattern` | 66.7% |
| 🟡 `034.go:runVar034` | 81.2% |
| 🟡 `034.go:isWaitGroupType` | 88.9% |
| 🔴 `035.go:runVar035` | 77.8% |
| 🔴 `035.go:checkContainsPattern` | 66.7% |
| 🔴 `036.go:runVar036` | 75.0% |
| 🟢 `036.go:analyzeFunctionForIndexPattern` | 93.8% |
| 🔴 `037.go:runVar037` | 77.8% |
| 🔴 `037.go:checkMapCollectionPatterns` | 76.9% |

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

### 🟢 `github.com/kodflow/ktn-linter/pkg/messages` - 94.2%

| Fichier:Fonction | Couverture |
|------------------|------------|
| 🟡 `messages.go:Format` | 83.3% |
| ⚫ `messages.go:GetAll` | 0.0% |
| ⚫ `messages.go:Clear` | 0.0% |
| ⚫ `messages.go:Unregister` | 0.0% |

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

*Généré le: 2026-01-02 12:09:37*

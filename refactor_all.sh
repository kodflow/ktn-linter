#!/bin/bash

# Liste des fichiers à refactoriser avec les fonctions cibles
FILES=(
    "src/pkg/analyzer/ktn/control_flow/if_001.go:runRuleIf001"
    "src/pkg/analyzer/ktn/control_flow/range_001.go:runRuleRange001"
    "src/pkg/analyzer/ktn/control_flow/switch_001.go:runRuleSwitch001"
    "src/pkg/analyzer/ktn/data_structures/array_001.go:runRuleArray001"
    "src/pkg/analyzer/ktn/data_structures/map_001.go:runRuleMap001"
    "src/pkg/analyzer/ktn/data_structures/slice_001.go:runRuleSlice001"
    "src/pkg/analyzer/ktn/interface/001.go:RunRule001"
    "src/pkg/analyzer/ktn/interface/001.go:CollectPackageInfo"
    "src/pkg/analyzer/ktn/interface/002.go:RunRule002"
    "src/pkg/analyzer/ktn/interface/003.go:RunRule003"
    "src/pkg/analyzer/ktn/interface/004.go:RunRule004"
    "src/pkg/analyzer/ktn/interface/005.go:RunRule005"
    "src/pkg/analyzer/ktn/method/001.go:methodModifiesReceiver"
    "src/pkg/analyzer/ktn/mock/002.go:runRule002"
    "src/pkg/analyzer/ktn/ops/chan_001.go:runRuleChan001"
    "src/pkg/analyzer/ktn/ops/comp_001.go:runRuleComp001"
    "src/pkg/analyzer/ktn/ops/conv_001.go:runRuleConv001"
    "src/pkg/analyzer/ktn/ops/pointer_001.go:runRulePointer001"
    "src/pkg/analyzer/ktn/ops/predecl_001.go:runRulePredecl001"
    "src/pkg/analyzer/ktn/ops/return_001.go:runRuleReturn001"
    "src/pkg/analyzer/ktn/registry.go:GetRulesByCategory"
    "src/pkg/analyzer/ktn/test/002.go:runRule002"
    "src/pkg/formatter/formatter.go:groupByFile"
)

echo "Starting refactoring process..."
for item in "${FILES[@]}"; do
    file="${item%%:*}"
    func="${item##*:}"
    echo "Processing $file :: $func"
done

echo "All files identified for refactoring"

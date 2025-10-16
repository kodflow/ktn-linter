#!/bin/bash

echo "🧪 Validation des scénarios de test KTN-Linter"
echo "=============================================="
echo ""

# Test 1: Code parfait
echo "✅ Test 1: Code parfait (doit avoir 0 violations)"
./builds/ktn-linter ./tests/target/rules_interface/ktn_interface_008_only_interfaces/... > /tmp/test1.out 2>&1
if grep -q "No issues found" /tmp/test1.out; then
    echo "   PASS ✅"
else
    echo "   FAIL ❌"
    cat /tmp/test1.out
    exit 1
fi

# Test 2: Structs dans interfaces.go
echo "✅ Test 2: Structs dans interfaces.go (doit avoir 6 violations)"
./builds/ktn-linter ./tests/source/rules_interface/ktn_interface_008_structs_in_interfaces/... > /tmp/test2.out 2>&1
count=$(grep "Total:" /tmp/test2.out | sed 's/\x1b\[[0-9;]*m//g' | grep -oP '\d+' | head -1)
if [ "$count" = "6" ]; then
    echo "   PASS ✅ (6 violations détectées)"
else
    echo "   FAIL ❌ (attendu 6, obtenu $count)"
    exit 1
fi

# Test 3: mock.go manquant
echo "✅ Test 3: mock.go manquant (doit avoir 3 violations)"
./builds/ktn-linter ./tests/source/rules_interface/ktn_mock_001_missing_mock/... > /tmp/test3.out 2>&1
count=$(grep "Total:" /tmp/test3.out | sed 's/\x1b\[[0-9;]*m//g' | grep -oP '\d+' | head -1)
if [ "$count" = "3" ]; then
    echo "   PASS ✅ (3 violations détectées)"
else
    echo "   FAIL ❌ (attendu 3, obtenu $count)"
    exit 1
fi

# Test 4: Global target
echo "✅ Test 4: Tous les targets (doit avoir 0 violations)"
./builds/ktn-linter ./tests/target/... > /tmp/test4.out 2>&1
if grep -q "No issues found" /tmp/test4.out; then
    echo "   PASS ✅"
else
    echo "   FAIL ❌"
    cat /tmp/test4.out
    exit 1
fi

# Test 5: Global source
echo "✅ Test 5: Tous les sources (doit avoir 420 violations)"
./builds/ktn-linter ./tests/source/... > /tmp/test5.out 2>&1
count=$(grep "Total:" /tmp/test5.out | sed 's/\x1b\[[0-9;]*m//g' | grep -oP '\d+' | head -1)
if [ "$count" = "420" ]; then
    echo "   PASS ✅ (420 violations détectées)"
else
    echo "   FAIL ❌ (attendu 420, obtenu $count)"
    cat /tmp/test5.out | tail -20
    exit 1
fi

echo ""
echo "🎉 Tous les tests de validation passent avec succès !"

# Validation des scénarios de test

## Commandes de validation

### Validation complète
```bash
# Tous les targets doivent passer (0 violations)
./builds/ktn-linter ./tests/target/...
# Attendu: ✅ No issues found! Code is compliant.

# Tous les sources doivent avoir 420 violations
./builds/ktn-linter ./tests/source/...
# Attendu: 📊 Total: 420 issue(s) to fix
```

### Validation par scénario

#### Scénario 1: Code parfait (target)
```bash
./builds/ktn-linter ./tests/target/rules_interface/ktn_interface_008_only_interfaces/...
```
**Attendu:** ✅ 0 violations

**Fichiers:**
- `interfaces.go` - Interfaces uniquement
- `impl.go` - Implémentations
- `impl_test.go` - Tests
- `mock.go` - Mocks avec build tag

---

#### Scénario 2: Structs dans interfaces.go (source)
```bash
./builds/ktn-linter ./tests/source/rules_interface/ktn_interface_008_structs_in_interfaces/...
```
**Attendu:** ❌ 6 violations

1. `KTN-MOCK-001` - mock.go manquant
2. `KTN-TEST-002` - interfaces.go sans test (car contient structs)
3. `KTN-INTERFACE-006` - Interface 'Service' sans constructeur
4. `KTN-INTERFACE-002` - Type public 'ServiceImpl' comme struct
5. `KTN-INTERFACE-008` - Structs dans interfaces.go ⚠️ **RÈGLE CIBLÉE**
6. `KTN-INTERFACE-006` - Interface 'Repository' sans constructeur

---

#### Scénario 3: mock.go manquant (source)
```bash
./builds/ktn-linter ./tests/source/rules_interface/ktn_mock_001_missing_mock/...
```
**Attendu:** ❌ 3 violations

1. `KTN-MOCK-001` - mock.go manquant ⚠️ **RÈGLE CIBLÉE**
2. `KTN-INTERFACE-006` - Interface 'Service' sans constructeur
3. `KTN-INTERFACE-006` - Interface 'Repository' sans constructeur

**Note:** Pas de KTN-TEST-002 car interfaces.go contient uniquement des interfaces (exemption).

---

## Résultats de validation attendus

| Scénario | Chemin | Violations | Status |
|----------|--------|------------|--------|
| Code parfait | `tests/target/rules_interface/ktn_interface_008_only_interfaces/` | 0 | ✅ |
| Structs dans interfaces.go | `tests/source/rules_interface/ktn_interface_008_structs_in_interfaces/` | 6 | ❌ |
| mock.go manquant | `tests/source/rules_interface/ktn_mock_001_missing_mock/` | 3 | ❌ |

## Règles testées

### KTN-INTERFACE-008
**Le fichier interfaces.go doit contenir UNIQUEMENT des interfaces**

Testé dans:
- ✅ `tests/target/rules_interface/ktn_interface_008_only_interfaces/` (respect)
- ❌ `tests/source/rules_interface/ktn_interface_008_structs_in_interfaces/` (violation)

### KTN-MOCK-001
**Si interfaces.go contient des interfaces, mock.go doit exister**

Testé dans:
- ✅ `tests/target/rules_interface/ktn_interface_008_only_interfaces/` (respect)
- ❌ `tests/source/rules_interface/ktn_interface_008_structs_in_interfaces/` (violation)
- ❌ `tests/source/rules_interface/ktn_mock_001_missing_mock/` (violation)

### Exception KTN-TEST-002
**interfaces.go contenant uniquement des interfaces est exempté**

Testé dans:
- ✅ `tests/target/rules_interface/ktn_interface_008_only_interfaces/` (interfaces.go sans test = OK)
- ✅ `tests/source/rules_interface/ktn_mock_001_missing_mock/` (interfaces.go sans test = OK)
- ❌ `tests/source/rules_interface/ktn_interface_008_structs_in_interfaces/` (interfaces.go avec structs = test requis)

## Script de validation automatique

```bash
#!/bin/bash
set -e

echo "🧪 Validation des scénarios de test KTN-Linter"
echo "=============================================="
echo ""

# Test 1: Code parfait
echo "✅ Test 1: Code parfait (doit avoir 0 violations)"
result=$(./builds/ktn-linter ./tests/target/rules_interface/ktn_interface_008_only_interfaces/... 2>&1)
if echo "$result" | grep -q "No issues found"; then
    echo "   PASS ✅"
else
    echo "   FAIL ❌"
    exit 1
fi

# Test 2: Structs dans interfaces.go
echo "✅ Test 2: Structs dans interfaces.go (doit avoir 6 violations)"
count=$(./builds/ktn-linter ./tests/source/rules_interface/ktn_interface_008_structs_in_interfaces/... 2>&1 | grep -oP '\d+(?= issue)' || echo "0")
if [ "$count" = "6" ]; then
    echo "   PASS ✅ (6 violations détectées)"
else
    echo "   FAIL ❌ (attendu 6, obtenu $count)"
    exit 1
fi

# Test 3: mock.go manquant
echo "✅ Test 3: mock.go manquant (doit avoir 3 violations)"
count=$(./builds/ktn-linter ./tests/source/rules_interface/ktn_mock_001_missing_mock/... 2>&1 | grep -oP '\d+(?= issue)' || echo "0")
if [ "$count" = "3" ]; then
    echo "   PASS ✅ (3 violations détectées)"
else
    echo "   FAIL ❌ (attendu 3, obtenu $count)"
    exit 1
fi

# Test 4: Global target
echo "✅ Test 4: Tous les targets (doit avoir 0 violations)"
result=$(./builds/ktn-linter ./tests/target/... 2>&1)
if echo "$result" | grep -q "No issues found"; then
    echo "   PASS ✅"
else
    echo "   FAIL ❌"
    exit 1
fi

# Test 5: Global source
echo "✅ Test 5: Tous les sources (doit avoir 420 violations)"
count=$(./builds/ktn-linter ./tests/source/... 2>&1 | grep -oP '\d+(?= issue)' || echo "0")
if [ "$count" = "420" ]; then
    echo "   PASS ✅ (420 violations détectées)"
else
    echo "   FAIL ❌ (attendu 420, obtenu $count)"
    exit 1
fi

echo ""
echo "🎉 Tous les tests de validation passent avec succès !"
```

Sauvegardez ce script dans `tests/validate.sh` et exécutez:
```bash
chmod +x tests/validate.sh
./tests/validate.sh
```

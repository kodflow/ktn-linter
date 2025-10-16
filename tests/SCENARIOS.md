# Scénarios de test KTN-Linter

## Vue d'ensemble

Ce document référence tous les scénarios de test pour les nouvelles règles **KTN-INTERFACE-008** et **KTN-MOCK-001**.

## Architecture des tests

### tests/target/ - Code Parfait ✅
Contient des exemples de code **conformes** à toutes les règles KTN.
**Résultat attendu : 0 violations**

### tests/source/ - Anti-patterns ❌
Contient des exemples de code **violant** les règles KTN (anti-patterns).
**Résultat attendu : N violations documentées**

## Scénarios pour KTN-INTERFACE-008 & KTN-MOCK-001

### ✅ Code Parfait (target/)

#### `tests/target/rules_interface/ktn_interface_008_only_interfaces/`

**Architecture idéale** démontrant la séparation interfaces/implémentations/mocks :

- `interfaces.go` → Interfaces uniquement
- `impl.go` → Implémentations privées + constructeurs
- `impl_test.go` → Tests des implémentations
- `mock.go` → Mocks réutilisables avec build tag `test`

**Violations : 0** ✅

**Règles respectées :**
- ✅ KTN-INTERFACE-008 : interfaces.go contient uniquement des interfaces
- ✅ KTN-MOCK-001 : mock.go existe et contient tous les mocks
- ✅ KTN-TEST-002 : mock.go et interfaces.go exemptés de tests
- ✅ Toutes les autres règles KTN

**Voir :** `tests/target/rules_interface/ktn_interface_008_only_interfaces/.README.md`

---

### ❌ Anti-patterns (source/)

#### 1. `tests/source/rules_interface/ktn_interface_008_structs_in_interfaces/`

**Violation principale :** Structs dans interfaces.go

**Problème :**
- Le fichier interfaces.go contient des structs (ServiceImpl, repositoryImpl)
- Violation de la règle KTN-INTERFACE-008

**Violations détectées : 6**
1. KTN-MOCK-001 : mock.go manquant
2. KTN-TEST-002 : interfaces.go sans test
3. KTN-INTERFACE-006 : Interface 'Service' sans constructeur
4. KTN-INTERFACE-002 : Type public 'ServiceImpl' comme struct
5. KTN-INTERFACE-008 : Structs dans interfaces.go ← **RÈGLE CIBLÉE**
6. KTN-INTERFACE-006 : Interface 'Repository' sans constructeur

**Solution :** Séparer en interfaces.go / impl.go / mock.go

**Voir :** `tests/source/rules_interface/ktn_interface_008_structs_in_interfaces/.README.md`

---

#### 2. `tests/source/rules_interface/ktn_mock_001_missing_mock/`

**Violation principale :** mock.go manquant

**Problème :**
- Le fichier interfaces.go contient 2 interfaces
- Aucun fichier mock.go n'existe
- Violation de la règle KTN-MOCK-001

**Violations détectées : 3**
1. KTN-MOCK-001 : mock.go manquant ← **RÈGLE CIBLÉE**
2. KTN-INTERFACE-006 : Interface 'Service' sans constructeur
3. KTN-INTERFACE-006 : Interface 'Repository' sans constructeur

**Note importante :** Pas de KTN-TEST-002 car interfaces.go contient uniquement des interfaces (exemption).

**Solution :** Créer mock.go avec build tag `//go:build test`

**Voir :** `tests/source/rules_interface/ktn_mock_001_missing_mock/.README.md`

---

## Règles d'exemption KTN-TEST-002

### Fichiers exemptés de l'obligation de tests

1. **mock.go** - Toujours exempté
   - Contient des mocks avec build tag `test`
   - Utilisé par d'autres tests

2. **interfaces.go** - Exempté conditionnellement
   - ✅ Exempté SI contient **uniquement** des interfaces
   - ❌ Requis SI contient structs, fonctions, ou autres types

### Logique d'exemption

```
interfaces.go
  ├─ Contient uniquement interfaces → ✅ Pas de interfaces_test.go requis
  └─ Contient structs/fonctions → ❌ interfaces_test.go requis
```

## Statistiques globales

### Target (code parfait)
```bash
./builds/ktn-linter ./tests/target/...
✅ No issues found! Code is compliant.
```

### Source (anti-patterns)
```bash
./builds/ktn-linter ./tests/source/...
❌ 420 issue(s) to fix
```

## Commandes de test

### Tester un scénario spécifique

```bash
# Code parfait
./builds/ktn-linter ./tests/target/rules_interface/ktn_interface_008_only_interfaces/...

# Anti-pattern : structs dans interfaces.go
./builds/ktn-linter ./tests/source/rules_interface/ktn_interface_008_structs_in_interfaces/...

# Anti-pattern : mock.go manquant
./builds/ktn-linter ./tests/source/rules_interface/ktn_mock_001_missing_mock/...
```

### Tester toutes les règles interface

```bash
# Tous les tests target interface
./builds/ktn-linter ./tests/target/rules_interface/...

# Tous les tests source interface
./builds/ktn-linter ./tests/source/rules_interface/...
```

## Résumé des nouvelles règles

### KTN-INTERFACE-008
**Le fichier interfaces.go doit contenir UNIQUEMENT des interfaces**

- ✅ Accepté : Interfaces avec godoc
- ❌ Refusé : Structs, types alias, implémentations, fonctions

### KTN-MOCK-001
**Si interfaces.go contient des interfaces, mock.go doit exister**

- ✅ mock.go avec build tag `//go:build test`
- ✅ Contient tous les mocks du package
- ✅ Vérification de compilation : `var _ Interface = (*Mock)(nil)`

### Exception KTN-TEST-002
**Fichiers exemptés de l'obligation de test**

1. mock.go (toujours)
2. interfaces.go contenant uniquement des interfaces (conditionnel)

# Guide Pratique de Corrections des Violations KTN

## Vue d'Ensemble

Ce guide fournit des exemples concrets pour corriger les 586 violations CRITICAL et 528 violations WARNING identifiées par le linter KTN.

## Table des Matières

1. [Violations CRITICAL](#violations-critical)
2. [Violations WARNING](#violations-warning)
3. [Scripts d'Aide](#scripts-daide)
4. [Procédure Recommandée](#procédure-recommandée)

---

## Violations CRITICAL

### 1. KTN-VAR-001: Variables Individuelles (221 violations)

**Problème:** Variables déclarées individuellement au lieu d'être regroupées.

**❌ Avant:**
```go
var Rule001 = &analysis.Analyzer{...}
var Rule002 = &analysis.Analyzer{...}
var Rule003 = &analysis.Analyzer{...}
```

**✅ Après:**
```go
var (
    // Rule001 checks for...
    Rule001 = &analysis.Analyzer{...}

    // Rule002 checks for...
    Rule002 = &analysis.Analyzer{...}

    // Rule003 checks for...
    Rule003 = &analysis.Analyzer{...}
)
```

**Commande de recherche:**
```bash
# Trouver les fichiers concernés
go run ./src/cmd/ktn-linter/main.go ./src/... 2>&1 | grep "KTN-VAR-001"
```

---

### 2. KTN-FUNC-001: Noms de Test avec Underscores (188 violations)

**Problème:** Les noms de fonctions de test utilisent des underscores au lieu de MixedCaps.

**❌ Avant:**
```go
func TestRule001_ConstGrouping(t *testing.T) {...}
func TestRule002_MakeSlicePrealloc(t *testing.T) {...}
```

**✅ Après:**
```go
func TestRule001ConstGrouping(t *testing.T) {...}
func TestRule002MakeSlicePrealloc(t *testing.T) {...}
```

**Correction automatique:**
```bash
# Utiliser le script fourni
./quick_fixes.sh
```

---

### 3. KTN-FUNC-009: Profondeur d'Imbrication Élevée (71 violations)

**Problème:** Fonctions avec profondeur d'imbrication > 3.

**❌ Avant:**
```go
func processData(data []Item) error {
    for _, item := range data {           // Niveau 1
        if item.Valid {                    // Niveau 2
            for _, sub := range item.Subs { // Niveau 3
                if sub.Active {             // Niveau 4 - TROP PROFOND
                    // ...
                }
            }
        }
    }
}
```

**✅ Après:**
```go
func processData(data []Item) error {
    for _, item := range data {
        if err := processItem(item); err != nil {
            return err
        }
    }
    return nil
}

func processItem(item Item) error {
    if !item.Valid {
        return nil
    }
    for _, sub := range item.Subs {
        if err := processSub(sub); err != nil {
            return err
        }
    }
    return nil
}

func processSub(sub SubItem) error {
    if !sub.Active {
        return nil
    }
    // Traitement du sub-item
    return nil
}
```

**Stratégie:**
1. Identifier les blocs profondément imbriqués
2. Extraire chaque niveau en fonction helper
3. Utiliser early returns pour réduire l'imbrication

---

### 4. KTN-ERROR-001: Mauvaise Gestion des Erreurs (30 violations)

**Problème:** Erreurs non gérées ou mal gérées.

**❌ Avant:**
```go
func loadFile(path string) Data {
    data, _ := os.ReadFile(path)  // Erreur ignorée
    return parseData(data)
}
```

**✅ Après:**
```go
func loadFile(path string) (Data, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return Data{}, fmt.Errorf("failed to read file %s: %w", path, err)
    }

    result, err := parseData(data)
    if err != nil {
        return Data{}, fmt.Errorf("failed to parse data: %w", err)
    }

    return result, nil
}
```

**Règles:**
1. Toujours vérifier les erreurs
2. Utiliser `%w` pour wrapper les erreurs
3. Ajouter du contexte aux erreurs
4. Propager les erreurs vers le haut

---

### 5. KTN-FUNC-006: Fonctions Trop Longues (27 violations)

**Problème:** Fonctions de plus de 35 lignes.

**Stratégie de correction:**
1. Identifier les blocs logiques distincts
2. Extraire chaque bloc en fonction helper
3. Garder la fonction principale comme orchestrateur

**Exemple:**
```go
// Avant: runAnalyzers() - 38 lignes
func runAnalyzers(pkgs []*packages.Package) []diagWithFset {
    analyzers := selectAnalyzers()
    var allDiagnostics []diagWithFset
    for _, pkg := range pkgs {
        analyzePackage(pkg, analyzers, &allDiagnostics)
    }
    return allDiagnostics
}

// Les détails sont dans selectAnalyzers() et analyzePackage()
```

---

### 6. KTN-GOROUTINE-002: Mauvaise Gestion des Goroutines (22 violations)

**Problème:** Goroutines sans synchronisation ou gestion d'erreur.

**❌ Avant:**
```go
func processItems(items []Item) {
    for _, item := range items {
        go processItem(item)  // Pas de synchronisation
    }
}
```

**✅ Après:**
```go
func processItems(ctx context.Context, items []Item) error {
    var wg sync.WaitGroup
    errCh := make(chan error, len(items))

    for _, item := range items {
        item := item  // Capture de la variable
        wg.Add(1)
        go func() {
            defer wg.Done()
            if err := processItem(ctx, item); err != nil {
                errCh <- err
            }
        }()
    }

    // Attendre la fin
    wg.Wait()
    close(errCh)

    // Vérifier les erreurs
    for err := range errCh {
        if err != nil {
            return err
        }
    }

    return nil
}
```

---

### 7. KTN-FUNC-007: Complexité Cyclomatique Élevée (15 violations)

**Problème:** Complexité cyclomatique > 10.

**Solution:** Décomposer en fonctions plus petites avec des responsabilités uniques.

---

### 8. KTN-GOROUTINE-001: Goroutines sans Gestion d'Erreur (11 violations)

Voir KTN-GOROUTINE-002 ci-dessus.

---

### 9. KTN-DEFER-001: Defer dans une Boucle (1 violation)

**Statut:** ✅ **INTENTIONNEL** - Fichier de test bad_usage

Cette violation est dans un fichier de test pour démontrer le mauvais usage.

**Localisation:** `src/tests/bad_usage/ktn/rules_func/ktn_func_edge_defer_panic/ktn_func_edge_defer_panic.go:52`

**Aucune correction nécessaire.**

---

## Violations WARNING

### 1. KTN-FUNC-002: Fonction sans Godoc (325 violations)

**❌ Avant:**
```go
func parseFlags() {
    flag.BoolVar(&aiMode, "ai", false, "Enable AI-friendly output")
    flag.Parse()
}
```

**✅ Après:**
```go
// parseFlags initializes and parses command-line flags.
func parseFlags() {
    flag.BoolVar(&aiMode, "ai", false, "Enable AI-friendly output")
    flag.Parse()
}
```

---

### 2. KTN-VAR-005: Variable Non Utilisée (142 violations)

**Solution:** Supprimer ou utiliser la variable.

**❌ Avant:**
```go
func process() {
    result := compute()
    // result non utilisé
}
```

**✅ Option 1 - Utiliser:**
```go
func process() {
    result := compute()
    log.Println("Result:", result)
}
```

**✅ Option 2 - Supprimer:**
```go
func process() {
    _ = compute()  // Explicitement ignoré
}
```

---

### 3. KTN-STRUCT-002: Struct sans Godoc (21 violations)

**❌ Avant:**
```go
type Config struct {
    Host string
    Port int
}
```

**✅ Après:**
```go
// Config holds the server configuration.
type Config struct {
    Host string
    Port int
}
```

---

## Scripts d'Aide

### Script de Correction Automatique

```bash
#!/bin/bash
# quick_fixes.sh - Corrections automatiques sûres

# 1. Corriger les noms de test
find ./src -name "*_test.go" -type f -exec \
    sed -i 's/func Test\([A-Za-z0-9]*\)_\([A-Za-z0-9]*\)(/func Test\1\2(/g' {} \;

# 2. Vérifier l'impact
go run ./src/cmd/ktn-linter/main.go ./src/... | head -20
```

### Analyser une Violation Spécifique

```bash
# Rechercher toutes les occurrences d'une règle
go run ./src/cmd/ktn-linter/main.go ./src/... 2>&1 | grep "KTN-VAR-001"

# Compter les violations par fichier
go run ./src/cmd/ktn-linter/main.go ./src/... 2>&1 | \
    grep "File:" | sort | uniq -c | sort -rn
```

### Tester les Corrections

```bash
# Avant correction
go run ./src/cmd/ktn-linter/main.go ./src/pkg/analyzer/ktn/const/001.go

# Après correction
# ... appliquer les corrections ...

go run ./src/cmd/ktn-linter/main.go ./src/pkg/analyzer/ktn/const/001.go
```

---

## Procédure Recommandée

### Phase 1: CRITICAL (Priorité Haute)

**Jour 1-2 (4-6 heures):**
1. ✅ Corriger KTN-VAR-001 (regrouper variables) - Semi-automatique
2. ✅ Corriger KTN-FUNC-001 (noms de test) - Automatique via script
3. ✅ Corriger KTN-FUNC-009 (imbrication) - Refactoring manuel

**Objectif:** Réduire de 586 → ~200 violations CRITICAL

### Phase 2: WARNING (Priorité Moyenne)

**Jour 3-4 (3-4 heures):**
1. ✅ Ajouter commentaires godoc (KTN-FUNC-002)
2. ✅ Nettoyer variables inutilisées (KTN-VAR-005)
3. ✅ Ajouter commentaires struct (KTN-STRUCT-002)

**Objectif:** Réduire de 528 → ~100 violations WARNING

### Phase 3: INFO (Optionnel)

**Si temps disponible:**
1. Compléter commentaires godoc (KTN-FUNC-003)
2. Commenter les champs struct (KTN-STRUCT-003)
3. Commenter les variables (KTN-VAR-003)

### Phase 4: MINOR (Optionnel - Non Recommandé)

Les 1228 violations MINOR (surtout KTN-FUNC-008) sont stylistiques et peuvent rester.

---

## Commandes Utiles

```bash
# Linter complet
go run ./src/cmd/ktn-linter/main.go ./src/...

# Par catégorie
go run ./src/cmd/ktn-linter/main.go -category=func ./src/...
go run ./src/cmd/ktn-linter/main.go -category=var ./src/...

# Mode simple (pour parsing)
go run ./src/cmd/ktn-linter/main.go -simple ./src/...

# Tests
go test ./src/...

# Coverage
go test -cover ./src/...
```

---

## Résumé

- **586 CRITICAL** → Corriger en priorité (4-6h)
- **528 WARNING** → Corriger ensuite (3-4h)
- **745 INFO** → Optionnel (6-8h)
- **1228 MINOR** → Ignorer (trop verbeux)

**Total réaliste:** 10-15 heures pour les violations importantes
**Impact:** Code plus maintenable, moins de bugs, meilleure documentation

---

**Bonne chance avec les corrections !** 🚀

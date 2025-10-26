# Go Expert Plugin - Exemples Pratiques

## 🎯 Cas d'Usage Réels

### Exemple 1 : Nouveau Projet Go

**Contexte** : Vous démarrez un nouveau projet Go

```bash
# Vous créez un nouveau dossier
mkdir my-api && cd my-api
go mod init github.com/user/my-api
```

**Avec le plugin activé** :

```
🤖 Claude détecte :
"Nouveau projet Go détecté. Je vais configurer les linters et best practices."

Création de :
  ✅ .golangci.yml (configuration stricte)
  ✅ Makefile (targets lint/test/build)
  ✅ CLAUDE.md (instructions projet)
  ✅ Installation golangci-lint

Vérifications :
  ✅ Go 1.25.3 installé
  ✅ golangci-lint v1.62.2 installé
  ✅ Structure projet standard créée

Prêt à coder ! 🚀
```

### Exemple 2 : Refactoring Fonction Trop Longue

**Code initial** :

```go
func ProcessOrder(orderID string, userID string, items []Item,
    discount float64, shippingAddress Address, paymentMethod string,
    giftWrap bool, giftMessage string) error {

    // Validation
    if orderID == "" {
        return errors.New("order ID required")
    }
    if userID == "" {
        return errors.New("user ID required")
    }
    if len(items) == 0 {
        return errors.New("items required")
    }

    // Calculate total
    var total float64
    for _, item := range items {
        total += item.Price * float64(item.Quantity)
    }

    // Apply discount
    if discount > 0 {
        total = total * (1 - discount/100)
    }

    // Calculate shipping
    shipping := calculateShipping(shippingAddress, items)
    total += shipping

    // Add gift wrap
    if giftWrap {
        total += 5.99
    }

    // Process payment
    payment := &Payment{
        Amount: total,
        Method: paymentMethod,
        UserID: userID,
    }

    if err := processPayment(payment); err != nil {
        return fmt.Errorf("payment failed: %w", err)
    }

    // Create order
    order := &Order{
        ID:      orderID,
        UserID:  userID,
        Items:   items,
        Total:   total,
        Status:  "pending",
    }

    // Save to database
    if err := saveOrder(order); err != nil {
        return fmt.Errorf("failed to save: %w", err)
    }

    // Send confirmation email
    sendConfirmationEmail(userID, order)

    return nil
}
```

**Plugin détecte** :

```
📁 File: order.go (4 issues)
────────────────────────────────────────────────

[1] order.go:12:1
  ✖ Code: KTN-FUNC-002
  ▶ function has 8 parameters (max 5). Use Functional Options or config struct

[2] order.go:12:1
  ⚠ Code: KTN-FUNC-001
  ▶ function is 48 lines (max 35). Extract to smaller functions

[3] order.go:12:1
  ℹ Code: KTN-FUNC-005
  ▶ function has side effects (email). Consider separating pure logic

🎯 Design Pattern Suggestion:
- Use Functional Options Pattern (8 params → options)
- Extract validation to separate function
- Extract calculation to OrderCalculator
- Use Builder for Order creation
```

**Refactoring suggéré par Claude** :

```go
// order.go

// OrderConfig configure une commande.
type OrderConfig struct {
    orderID         string
    userID          string
    items           []Item
    discount        float64
    shippingAddress Address
    paymentMethod   string
    giftWrap        bool
    giftMessage     string
}

// Option configure OrderConfig.
type Option func(*OrderConfig)

// WithDiscount applique une réduction.
func WithDiscount(discount float64) Option {
    return func(c *OrderConfig) { c.discount = discount }
}

// WithGiftWrap active l'emballage cadeau.
func WithGiftWrap(message string) Option {
    return func(c *OrderConfig) {
        c.giftWrap = true
        c.giftMessage = message
    }
}

// ProcessOrder traite une commande avec options.
//
// Params:
//   - orderID: identifiant unique commande
//   - userID: identifiant utilisateur
//   - items: articles commandés
//   - opts: options supplémentaires
//
// Returns:
//   - error: erreur éventuelle
func ProcessOrder(orderID, userID string, items []Item, opts ...Option) error {
    // Configuration
    config := &OrderConfig{
        orderID: orderID,
        userID:  userID,
        items:   items,
    }

    // Application options
    for _, opt := range opts {
        opt(config)
    }

    // Validation
    if err := validateOrder(config); err != nil {
        // Retour erreur validation
        return err
    }

    // Calcul du total
    total := calculateTotal(config)

    // Traitement paiement
    if err := processOrderPayment(config, total); err != nil {
        // Retour erreur paiement
        return err
    }

    // Création et sauvegarde
    order := buildOrder(config, total)
    if err := saveOrder(order); err != nil {
        // Retour erreur sauvegarde
        return fmt.Errorf("failed to save: %w", err)
    }

    // Notification (async)
    go sendConfirmationEmail(userID, order)

    // Retour succès
    return nil
}

// validateOrder valide la configuration commande.
//
// Params:
//   - config: configuration à valider
//
// Returns:
//   - error: erreur de validation
func validateOrder(config *OrderConfig) error {
    // Vérification orderID
    if config.orderID == "" {
        // Retour erreur
        return errors.New("order ID required")
    }

    // Vérification userID
    if config.userID == "" {
        // Retour erreur
        return errors.New("user ID required")
    }

    // Vérification items
    if len(config.items) == 0 {
        // Retour erreur
        return errors.New("items required")
    }

    // Retour succès
    return nil
}

// calculateTotal calcule le total de la commande.
//
// Params:
//   - config: configuration commande
//
// Returns:
//   - float64: montant total
func calculateTotal(config *OrderConfig) float64 {
    // Calcul prix items
    var total float64
    for _, item := range config.items {
        total += item.Price * float64(item.Quantity)
    }

    // Application réduction
    if config.discount > 0 {
        total = total * (1 - config.discount/100)
    }

    // Ajout frais de port
    shipping := calculateShipping(config.shippingAddress, config.items)
    total += shipping

    // Ajout emballage cadeau
    if config.giftWrap {
        const GIFT_WRAP_PRICE float64 = 5.99
        total += GIFT_WRAP_PRICE
    }

    // Retour total
    return total
}
```

**Résultat du plugin** :

```
✅ Refactoring terminé !

Avant :
  ✖ 8 paramètres (max 5)
  ⚠ 48 lignes (max 35)
  ℹ Side effects non isolés

Après :
  ✅ Functional Options Pattern
  ✅ 4 fonctions < 25 lignes chacune
  ✅ Validation isolée
  ✅ Calculs purs séparés
  ✅ Email async (pas de blocage)

make lint: ✅ 0 issues
make test: ✅ PASS (coverage 94.2%)
```

### Exemple 3 : Détection Worker Pool

**Code initial** :

```go
func ProcessFiles(files []string) error {
    for _, file := range files {
        data, err := os.ReadFile(file)
        if err != nil {
            return err
        }

        result := process(data)

        if err := save(result); err != nil {
            return err
        }
    }
    return nil
}
```

**Plugin détecte** :

```
🎯 Performance Suggestion:

  ⚠ Sequential processing detected for I/O operations

  Suggestion: Use Worker Pool Pattern
  - Current: 1 file at a time (slow for 1000+ files)
  - Optimized: 10 workers concurrent (10x faster)

  Estimated improvement: 85% faster
```

**Refactoring suggéré** :

```go
// Task représente une tâche de traitement fichier.
type Task struct {
    File string
}

// Result représente le résultat d'un traitement.
type Result struct {
    File  string
    Data  []byte
    Error error
}

// ProcessFiles traite des fichiers en parallèle.
//
// Params:
//   - files: liste des fichiers à traiter
//
// Returns:
//   - error: erreur éventuelle
func ProcessFiles(files []string) error {
    const NUM_WORKERS int = 10

    // Création des channels
    tasks := make(chan Task, len(files))
    results := make(chan Result, len(files))

    // Démarrage workers
    var wg sync.WaitGroup
    for i := 0; i < NUM_WORKERS; i++ {
        wg.Add(1)
        go worker(tasks, results, &wg)
    }

    // Envoi des tâches
    for _, file := range files {
        tasks <- Task{File: file}
    }
    close(tasks)

    // Attente completion
    go func() {
        wg.Wait()
        close(results)
    }()

    // Collecte résultats
    for result := range results {
        // Vérification erreur
        if result.Error != nil {
            // Retour première erreur
            return fmt.Errorf("failed to process %s: %w", result.File, result.Error)
        }
    }

    // Retour succès
    return nil
}

// worker traite les tâches du channel.
//
// Params:
//   - tasks: channel des tâches
//   - results: channel des résultats
//   - wg: waitgroup pour synchronisation
func worker(tasks <-chan Task, results chan<- Result, wg *sync.WaitGroup) {
    // Signaler completion
    defer wg.Done()

    // Traiter les tâches
    for task := range tasks {
        // Lecture fichier
        data, err := os.ReadFile(task.File)
        if err != nil {
            results <- Result{File: task.File, Error: err}
            // Continuer avec tâche suivante
            continue
        }

        // Traitement
        processed := process(data)

        // Sauvegarde
        if err := save(processed); err != nil {
            results <- Result{File: task.File, Error: err}
            // Continuer avec tâche suivante
            continue
        }

        // Résultat succès
        results <- Result{File: task.File, Data: processed}
    }
}
```

**Résultat** :

```
✅ Worker Pool Pattern appliqué

Performance:
  Before: 45.2s (1000 files sequential)
  After:  4.8s (1000 files, 10 workers)
  Gain: 89.4% faster ⚡

make lint: ✅ 0 issues
make test: ✅ PASS (coverage 92.1%)
```

### Exemple 4 : Intégration Projet Existant

**Scénario** : Vous héritez d'un vieux projet Go

```bash
cd legacy-project
```

**Plugin s'active** :

```
🔍 Analyse du projet...

Détections:
  ⚠ Go 1.18 (outdated, 1.25 available)
  ❌ golangci-lint non installé
  ❌ Pas de Makefile
  ⚠ 847 violations détectées

🤖 Je propose:
  1. Upgrade Go 1.18 → 1.25
  2. Installer golangci-lint
  3. Créer Makefile et .golangci.yml
  4. Fixer les violations par priorité (ERROR → WARNING → INFO)

Accepter ? (y/n)
```

**Après acceptation** :

```
✅ Configuration terminée

Étape 1/4: Upgrade Go
  ✅ go.mod: go 1.18 → go 1.25
  ✅ Dépendances mises à jour

Étape 2/4: Installation linters
  ✅ golangci-lint v1.62.2 installé

Étape 3/4: Configuration
  ✅ .golangci.yml créé
  ✅ Makefile créé

Étape 4/4: Premier lint
  ✖ 127 ERRORS (blocants)
  ⚠ 453 WARNINGS (maintenabilité)
  ℹ 267 INFO (style)

📊 Plan de correction:
  - Jour 1: Fix 127 ERRORS (priorité critique)
  - Jour 2-3: Fix 453 WARNINGS (par batch de 100)
  - Jour 4-5: Fix 267 INFO (amélioration continue)

Commencer par les ERRORS ? (y/n)
```

## 🎓 Best Practices Appliquées

### Naming Conventions

```go
// ❌ AVANT
const api_key = "secret"
var user_name string
func get_user() {}

// ✅ APRÈS (plugin auto-fix)
const API_KEY string = "secret"
var userName string
func getUser() {}
```

### Error Handling

```go
// ❌ AVANT
func ReadFile(path string) []byte {
    data, _ := os.ReadFile(path)  // erreur ignorée !
    return data
}

// ✅ APRÈS (plugin détecte + suggère)
func ReadFile(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    // Vérification erreur
    if err != nil {
        // Retour erreur wrappée
        return nil, fmt.Errorf("failed to read %s: %w", path, err)
    }
    // Retour données
    return data, nil
}
```

### Context Propagation

```go
// ❌ AVANT
func FetchUser(userID string) (*User, error) {
    // Pas de timeout, peut bloquer indéfiniment
}

// ✅ APRÈS (plugin suggère)
func FetchUser(ctx context.Context, userID string) (*User, error) {
    // Vérification cancellation
    select {
    case <-ctx.Done():
        // Retour erreur context
        return nil, ctx.Err()
    default:
    }

    // Suite du traitement avec ctx
}
```

## 🚀 Workflow Complet

```bash
# 1. Créer feature
git checkout -b feat/user-service

# 2. Coder (plugin auto-lint à chaque save)
vim user_service.go

# 3. Plugin exécute automatiquement
make lint  # ✅ 0 issues
make test  # ✅ PASS

# 4. Commit (pre-commit hook s'active)
git add .
git commit -m "feat: add user service"
# → Hook vérifie: lint + test + build
# → Bloque si échecs

# 5. Push
git push origin feat/user-service

# 6. CI/CD GitHub Actions
# → Même vérifications (lint, test, build)
# → Merge si ✅
```

---

**Le plugin transforme Claude en expert Go qui ne laisse RIEN passer !** 🎯

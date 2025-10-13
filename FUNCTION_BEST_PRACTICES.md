# Bonnes pratiques Go pour les FONCTIONS

Ce document compile TOUTES les bonnes pratiques officielles et communautaires concernant les fonctions en Go.

## 📋 Sources
- Effective Go (go.dev/doc/effective_go)
- Google Go Style Guide (google.github.io/styleguide/go/)
- Clean Code principles
- Cognitive Complexity research (SonarSource, gocognit)

---

## 1. 🏷️ NAMING (Nommage des fonctions)

### 1.1 Convention MixedCaps
- **Règle** : Utiliser `MixedCaps` (exportées) ou `mixedCaps` (non-exportées)
- ❌ **INTERDIT** : snake_case, kebab-case, UPPER_CASE
- ✅ **Correct** : `ParseHTTPRequest`, `calculateTotal`
- ❌ **Incorrect** : `parse_http_request`, `Calculate_Total`

### 1.2 Initialismes
- **Règle** : Les initialismes restent en majuscules dans MixedCaps
- ✅ **Correct** : `HTTPServer`, `URLParser`, `IDGenerator`, `XMLParser`
- ❌ **Incorrect** : `HttpServer`, `UrlParser`, `IdGenerator`

### 1.3 Noms descriptifs selon le rôle
- **Fonctions qui retournent quelque chose** → Nom de type nom (noun-like)
  - ✅ `UserName()`, `Config()`, `Token()`

- **Fonctions qui font quelque chose** → Nom de type verbe (verb-like)
  - ✅ `ParseRequest()`, `ValidateInput()`, `SendEmail()`

### 1.4 Éviter les préfixes redondants
- ❌ **ÉVITER** : `GetUserName()`, `SetUserName()`
- ✅ **PRÉFÉRER** : `UserName()`, `SetUserName()` (setter OK avec Set)

### 1.5 Getters sans "Get"
- Si le champ est `owner`, le getter doit être `Owner()`, pas `GetOwner()`
- ✅ `user.Email()` pas `user.GetEmail()`

### 1.6 Nom basé sur le résultat, pas l'implémentation
- ✅ `CalculateTotal()` (décrit le résultat)
- ❌ `LoopThroughItemsAndSum()` (décrit l'implémentation)

### 1.7 Fonctions qui diffèrent par le type
- Inclure le nom du type à la fin
- ✅ `ParseInt()`, `ParseFloat()`, `ParseBool()`

### 1.8 Interfaces à une méthode
- Suffixe `-er` pour les interfaces à une méthode
- ✅ `Reader`, `Writer`, `Formatter`, `CloseNotifier`

---

## 2. 📝 DOCUMENTATION

### 2.1 Commentaire obligatoire
- **Règle** : Toute fonction exportée DOIT avoir un commentaire de documentation
- Format : `// FunctionName does something...`

### 2.2 Format godoc
```go
// CalculateTotal calcule le total des articles dans le panier.
// Retourne une erreur si le panier est vide ou si un article est invalide.
func CalculateTotal(items []Item) (float64, error) {
    // ...
}
```

### 2.3 Documenter les comportements non-évidents
- Documenter les paramètres error-prone
- Expliquer POURQUOI un paramètre est important (pas juste ce qu'il est)
- Documenter les cas particuliers et edge cases

### 2.4 Documenter les valeurs de retour d'erreur
- Documenter les erreurs sentinelles significatives
- Expliquer quand chaque type d'erreur peut survenir

---

## 3. 📥 PARAMÈTRES

### 3.1 Ordre des paramètres
1. **Context en premier** (si présent)
   ```go
   func Process(ctx context.Context, data string) error
   ```

2. **Paramètres requis avant optionnels**
   ```go
   func NewServer(addr string, opts ...Option) *Server
   ```

### 3.2 Nombre de paramètres
- **Limite recommandée** : ≤ 5 paramètres
- Si > 5 paramètres → Envisager une struct d'options

  ❌ **Trop de paramètres** :
  ```go
  func NewServer(addr string, port int, timeout time.Duration,
                 maxConn int, readBuf int, writeBuf int) *Server
  ```

  ✅ **Struct d'options** :
  ```go
  type ServerConfig struct {
      Addr      string
      Port      int
      Timeout   time.Duration
      MaxConn   int
      ReadBuf   int
      WriteBuf  int
  }

  func NewServer(cfg ServerConfig) *Server
  ```

### 3.3 Options variadiques
- Pour configuration flexible :
  ```go
  type Option func(*Server)

  func WithTimeout(d time.Duration) Option {
      return func(s *Server) { s.timeout = d }
  }

  func NewServer(addr string, opts ...Option) *Server
  ```

### 3.4 Éviter `testing.T` en paramètre
- Les helpers de test ne doivent PAS prendre `*testing.T` en paramètre
- Retourner les valeurs et laisser l'appelant gérer les assertions

---

## 4. 📤 VALEURS DE RETOUR

### 4.1 Valeurs multiples
- **Utiliser** les valeurs de retour multiples, surtout pour les erreurs
  ```go
  func ParseUser(data []byte) (User, error)
  ```

### 4.2 Named returns (Retours nommés)
- **Utiliser** pour la documentation quand utile
- **ÉVITER** pour les fonctions courtes ou évidentes

  ✅ **Bon usage** :
  ```go
  func Divide(a, b float64) (quotient float64, remainder float64, err error) {
      if b == 0 {
          err = errors.New("division by zero")
          return
      }
      quotient = a / b
      remainder = math.Mod(a, b)
      return
  }
  ```

  ❌ **Naked return dans fonction longue** :
  ```go
  func ComplexCalculation(...) (result int, err error) {
      // ... 50 lignes de code ...
      return // Pas clair ce qui est retourné
  }
  ```

### 4.3 Ordre des valeurs de retour
- Convention Go : **(résultat, erreur)**
  ```go
  func ReadFile(path string) ([]byte, error)
  ```

### 4.4 Éviter les retours bool pour les erreurs
- ❌ `func Process() bool` (succès/échec)
- ✅ `func Process() error` (retourne l'erreur)

---

## 5. ⚠️ GESTION DES ERREURS

### 5.1 Retourner les erreurs, ne pas paniquer
- **Règle** : `return error` plutôt que `panic()`
- `panic` réservé aux erreurs non récupérables

### 5.2 Toujours vérifier les erreurs
```go
result, err := DoSomething()
if err != nil {
    return fmt.Errorf("doing something: %w", err)
}
```

### 5.3 Ajouter du contexte aux erreurs
- **Wrapper les erreurs** avec `%w`
  ```go
  if err != nil {
      return fmt.Errorf("reading config file: %w", err)
  }
  ```

- **Position du %w** : Toujours à la FIN du message
  ```go
  // ✅ Correct
  fmt.Errorf("cannot process user %s: %w", userID, err)

  // ❌ Incorrect
  fmt.Errorf("%w: cannot process user %s", err, userID)
  ```

### 5.4 Erreurs sentinelles
```go
var (
    ErrNotFound = errors.New("not found")
    ErrInvalid  = errors.New("invalid input")
)
```

### 5.5 Origine des erreurs
- Les messages d'erreur doivent identifier leur origine
  ```go
  return fmt.Errorf("database: failed to connect: %w", err)
  ```

---

## 6. 🧩 COMPLEXITÉ

### 6.1 Complexité cognitive
- **Cible** : < 15 par fonction
- **Idéal** : < 10 en moyenne

### 6.2 Longueur des fonctions
- **Recommandation** : < 50 lignes
- Si > 50 lignes → Envisager de découper

### 6.3 Réduire la complexité

#### Extraire des méthodes
```go
// ❌ Trop complexe
func ProcessOrder(order Order) error {
    // validation (10 lignes)
    // calcul (15 lignes)
    // sauvegarde (10 lignes)
    // notification (5 lignes)
}

// ✅ Décomposé
func ProcessOrder(order Order) error {
    if err := validateOrder(order); err != nil {
        return err
    }
    total := calculateTotal(order)
    if err := saveOrder(order, total); err != nil {
        return err
    }
    return notifyCustomer(order)
}
```

#### Éviter les imbrications profondes
```go
// ❌ Imbrication profonde
func Process(data []Data) error {
    for _, d := range data {
        if d.Valid {
            if d.Type == "A" {
                if d.Value > 0 {
                    // ...
                }
            }
        }
    }
}

// ✅ Guard clauses
func Process(data []Data) error {
    for _, d := range data {
        if !d.Valid {
            continue
        }
        if d.Type != "A" {
            continue
        }
        if d.Value <= 0 {
            continue
        }
        // traitement principal
    }
}
```

#### Utiliser switch plutôt que if-else en chaîne
```go
// ❌ Chaîne if-else
if status == "pending" {
    // ...
} else if status == "processing" {
    // ...
} else if status == "completed" {
    // ...
}

// ✅ Switch
switch status {
case "pending":
    // ...
case "processing":
    // ...
case "completed":
    // ...
}
```

---

## 7. 🧹 CLEAN CODE

### 7.1 Single Responsibility
- Une fonction = une responsabilité
- Si le nom contient "And", c'est suspect

### 7.2 Defer pour le cleanup
```go
func ReadFile(path string) ([]byte, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close() // ✅ Garantit la fermeture

    return io.ReadAll(f)
}
```

### 7.3 Éviter les side effects cachés
- Les fonctions doivent être prévisibles
- Documenter tout side effect non-évident

### 7.4 Testabilité
- Concevoir les fonctions pour être facilement testables
- Éviter les dépendances globales
- Préférer l'injection de dépendances

---

## 8. 📊 MÉTRIQUES RECOMMANDÉES

| Métrique | Valeur cible | Critique |
|----------|--------------|----------|
| Complexité cognitive | < 15 | > 25 |
| Longueur (lignes) | < 50 | > 100 |
| Nombre de paramètres | ≤ 5 | > 7 |
| Profondeur imbrication | < 4 | > 6 |
| Valeurs de retour | ≤ 3 | > 4 |

---

## 9. ✅ CHECKLIST FONCTION PARFAITE

- [ ] Nom en MixedCaps/mixedCaps
- [ ] Initialismes en majuscules (HTTP, URL, ID...)
- [ ] Nom descriptif du résultat/comportement
- [ ] Commentaire godoc complet
- [ ] ≤ 5 paramètres
- [ ] Context en premier (si présent)
- [ ] Retour (résultat, error) pour opérations faillibles
- [ ] Toutes les erreurs sont vérifiées
- [ ] Erreurs wrappées avec contexte (%w)
- [ ] Complexité cognitive < 15
- [ ] Longueur < 50 lignes
- [ ] Single responsibility
- [ ] Defer pour cleanup
- [ ] Testable sans dépendances globales

---

## 10. 🚫 ANTI-PATTERNS À ÉVITER

1. ❌ Fonctions avec snake_case
2. ❌ Préfixes Get inutiles
3. ❌ Trop de paramètres (> 5)
4. ❌ Naked returns dans fonctions longues
5. ❌ Ignorer les erreurs
6. ❌ Panic pour erreurs récupérables
7. ❌ Fonctions > 100 lignes
8. ❌ Complexité cognitive > 25
9. ❌ Side effects non documentés
10. ❌ Noms basés sur l'implémentation

---

## 11. 📚 EXEMPLES COMPLETS

### Exemple PARFAIT
```go
// ParseUserConfig parse et valide la configuration utilisateur depuis un fichier JSON.
// Retourne une erreur si le fichier n'existe pas, est malformé, ou contient des valeurs invalides.
func ParseUserConfig(ctx context.Context, path string) (*UserConfig, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("reading config file %s: %w", path, err)
    }

    var cfg UserConfig
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("parsing config JSON: %w", err)
    }

    if err := cfg.Validate(); err != nil {
        return nil, fmt.Errorf("invalid config: %w", err)
    }

    return &cfg, nil
}
```

### Exemple avec Options
```go
// ServerOption configure un Server.
type ServerOption func(*Server)

// WithTimeout définit le timeout du serveur.
func WithTimeout(d time.Duration) ServerOption {
    return func(s *Server) { s.timeout = d }
}

// WithMaxConnections définit le nombre maximum de connexions.
func WithMaxConnections(n int) ServerOption {
    return func(s *Server) { s.maxConn = n }
}

// NewServer crée un nouveau serveur avec les options spécifiées.
func NewServer(addr string, opts ...ServerOption) *Server {
    s := &Server{
        addr:    addr,
        timeout: 30 * time.Second, // défaut
        maxConn: 100,              // défaut
    }

    for _, opt := range opts {
        opt(s)
    }

    return s
}

// Usage
server := NewServer("localhost:8080",
    WithTimeout(60*time.Second),
    WithMaxConnections(200),
)
```

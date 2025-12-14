# KTN-VAR-018

**Sévérité**: WARNING

## Description

Les variables doivent utiliser `camelCase`, pas `snake_case`.

## Pourquoi

Go utilise **camelCase** pour toutes les variables:
- Variables locales: `userName`, `maxSize`
- Variables package: `defaultTimeout`, `httpClient`
- Constantes: `MaxRetries` (PascalCase si exportées)

`snake_case` n'est **pas** idiomatique en Go et rend le code inconsistant avec l'écosystème.

## Exemple incorrect

```go
var user_name string      // snake_case
var max_retry_count int   // snake_case
var http_client *Client   // snake_case

func process() {
    file_path := "/tmp/data"  // snake_case
    error_count := 0          // snake_case
}
```

## Exemple correct

```go
var userName string      // camelCase
var maxRetryCount int    // camelCase
var httpClient *Client   // camelCase

func process() {
    filePath := "/tmp/data"  // camelCase
    errorCount := 0          // camelCase
}
```

## Exceptions

Les seuls cas où `_` est accepté:
- `_` seul pour ignorer une valeur: `_, err := f()`
- Préfixe `_` pour paramètres non utilisés: `func f(_unused int)`

## Configuration

Cette règle peut être désactivée:

```yaml
rules:
  KTN-VAR-018:
    enabled: false
```

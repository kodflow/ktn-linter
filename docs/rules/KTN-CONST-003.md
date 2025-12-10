# KTN-CONST-003

**Sévérité**: INFO

## Description

Les constantes doivent utiliser la convention de nommage CamelCase standard de Go.

- **Exportées** : PascalCase (ex: `MaxSize`, `HttpTimeout`, `APIKey`)
- **Non exportées** : camelCase (ex: `maxSize`, `httpTimeout`, `apiKey`)

Les underscores (`_`) ne sont pas autorisés dans les noms de constantes (sauf pour le blank identifier `_`).

## Exemple non conforme

```go
const (
    // ❌ SCREAMING_SNAKE_CASE
    MAX_SIZE int = 100
    API_KEY string = "secret"
    HTTP_TIMEOUT int = 30

    // ❌ snake_case
    max_size int = 100
    api_key string = "secret"

    // ❌ Mixed_Case
    Max_Size int = 100
)
```

## Exemple conforme

```go
const (
    // ✅ PascalCase pour exportées
    MaxSize int = 100
    ApiKey string = "secret"
    HttpTimeout int = 30

    // ✅ camelCase pour non exportées
    maxInternalSize int = 50
    defaultTimeout int = 10

    // ✅ Acronymes
    APIEndpoint string = "/api"
    HTTPStatus int = 200

    // ✅ Avec chiffres
    Http2Protocol string = "h2"
    Version100 string = "1.0.0"
)
```

## Pourquoi CamelCase ?

Cette règle suit les conventions idiomatiques de Go :
- Cohérence avec le reste de l'écosystème Go
- Les constantes exportées/non-exportées sont clairement identifiables
- Meilleure lisibilité dans le code Go standard

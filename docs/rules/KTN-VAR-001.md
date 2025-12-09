# KTN-VAR-001

**Sévérité**: ERROR

## Description

Les variables de package doivent utiliser la convention de nommage Go standard :
- **Variables privées** : `camelCase` (ex: `maxRetries`, `httpClient`)
- **Variables exportées** : `PascalCase` (ex: `MaxRetries`, `HTTPClient`)

Les noms en `SCREAMING_SNAKE_CASE` (avec underscores et majuscules) sont interdits pour les variables.

## Gestion des acronymes

Les acronymes suivent la convention Go standard :
- **Exportées** : acronymes en majuscules (`HTTPClient`, `XMLParser`, `APIKey`)
- **Privées** : première lettre de l'acronyme en minuscule (`httpClient`, `xmlParser`, `apiKey`)

Ces formes sont **valides** car elles n'utilisent pas d'underscores.

## Exemple non conforme

```go
var (
    // ❌ SCREAMING_SNAKE_CASE interdit
    MAX_VALUE int = 100
    SERVER_PORT int = 8080
    HTTP_CLIENT string = "client"
    API_KEY string = "secret"
)
```

## Exemple conforme

```go
var (
    // ✅ camelCase pour variables privées
    maxValue int = 100
    serverPort int = 8080
    httpClient string = "client"
    apiKey string = "secret"

    // ✅ PascalCase pour variables exportées
    MaxValue int = 100
    ServerPort int = 8080
    HTTPClient string = "client"
    APIKey string = "secret"
)
```

## Distinction avec les constantes

Depuis KTN-CONST-003, les constantes utilisent également `CamelCase` (pas `SCREAMING_SNAKE_CASE`).
La distinction se fait par le mot-clé (`const` vs `var`), pas par la casse.

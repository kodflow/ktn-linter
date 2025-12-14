# KTN-STRUCT-006

**Sévérité**: INFO

## Description

Les champs privés ne doivent pas avoir de tags de sérialisation (json, xml, yaml).

## Pourquoi

Les champs privés (commençant par une minuscule) ne sont **jamais** sérialisés par les packages Go standard (encoding/json, encoding/xml, gopkg.in/yaml.v3). Ajouter un tag sur un champ privé est donc:
- Inutile (le tag est ignoré)
- Trompeur (suggère que le champ sera sérialisé)
- Source de confusion pour les autres développeurs

## Exemple incorrect

```go
type User struct {
    name     string `json:"name"`     // Tag ignoré!
    password string `json:"password"` // Tag ignoré!
}
```

## Exemple correct

**Option 1: Exporter le champ si nécessaire**
```go
type User struct {
    Name string `json:"name"`
}
```

**Option 2: Garder privé sans tag**
```go
type User struct {
    name     string // Pas de tag, c'est intentionnel
    password string
}
```

## Configuration

Cette règle ne peut pas être désactivée via la configuration.

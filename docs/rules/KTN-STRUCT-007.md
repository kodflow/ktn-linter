# KTN-STRUCT-007

**Sévérité**: INFO

## Description

Les getters et setters doivent suivre la convention Go de nommage.

## Convention Go

| Type | Nom du champ | Getter | Setter |
|------|-------------|--------|--------|
| Standard | `name` | `Name()` | `SetName(v)` |
| Boolean | `valid` | `IsValid()` | `SetValid(v)` |

**Important**: Les getters/setters sont **optionnels**. Cette règle ne s'applique que si vous choisissez d'en créer.

## Pourquoi

- Go n'utilise **pas** le préfixe `Get` pour les getters
- Le nom du getter doit correspondre au nom du champ (en exporté)
- Cette convention est documentée dans [Effective Go](https://golang.org/doc/effective_go#Getters)

## Exemple incorrect

```go
type User struct {
    name string
}

func (u *User) GetName() string { // Préfixe Get inutile
    return u.name
}

func (u *User) GetAge() int { // Retourne 'age' mais s'appelle GetAge
    return u.age
}
```

## Exemple correct

```go
type User struct {
    name string
    age  int
}

// Getter: pas de préfixe Get
func (u *User) Name() string {
    return u.name
}

// Getter: correspond au champ
func (u *User) Age() int {
    return u.age
}

// Setter: préfixe Set
func (u *User) SetName(name string) {
    u.name = name
}
```

## Configuration

Cette règle ne peut pas être désactivée via la configuration.

package struct003

// User champs exportés avant privés - CONFORME.
// Représente un utilisateur avec champs publics et privés.
type User struct {
	Name  string // exporté
	Age   int    // exporté
	id    int    // privé
	email string // privé
}

// AllPublic tous exportés - CONFORME.
// Structure avec uniquement des champs publics.
type AllPublic struct {
	Name    string
	Age     int
	Address string
}

// AllPrivate tous privés - CONFORME.
// Structure avec uniquement des champs privés.
type AllPrivate struct {
	name    string
	age     int
	address string
}

// EmptyStruct struct vide - CONFORME.
// Structure vide utilisée comme marqueur.
type EmptyStruct struct{}

// OnlyPublic un seul champ exporté - CONFORME.
// Structure avec un seul champ public.
type OnlyPublic struct {
	Name string
}

// OnlyPrivate un seul champ privé - CONFORME.
// Structure avec un seul champ privé.
type OnlyPrivate struct {
	name string
}

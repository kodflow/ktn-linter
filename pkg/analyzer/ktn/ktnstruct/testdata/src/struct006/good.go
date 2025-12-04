// Good examples for the struct003 test case.
package struct006

// UserModel champs exportés avant privés - CONFORME.
// Représente un utilisateur avec champs publics et privés.
type UserModel struct {
	Name  string `json:"name"`  // exporté
	Age   int    `json:"age"`   // exporté
	id    int    `json:"-"`     // privé
	email string `json:"-"`     // privé
}

// AllPublicModel tous exportés - CONFORME.
// Structure avec uniquement des champs publics.
type AllPublicModel struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
}

// AllPrivateModel tous privés - CONFORME.
// Structure avec uniquement des champs privés.
type AllPrivateModel struct {
	name    string `json:"-"`
	age     int    `json:"-"`
	address string `json:"-"`
}

// EmptyData struct vide - CONFORME.
// Structure vide utilisée comme marqueur.
type EmptyData struct{}

// OnlyPublicModel un seul champ exporté - CONFORME.
// Structure avec un seul champ public.
type OnlyPublicModel struct {
	Name string `json:"name"`
}

// OnlyPrivateModel un seul champ privé - CONFORME.
// Structure avec un seul champ privé.
type OnlyPrivateModel struct {
	name string `json:"-"`
}

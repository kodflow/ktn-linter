// Good examples for the struct003 test case.
package struct005

// UserModel champs exportés avant privés - CONFORME.
// Représente un utilisateur avec champs publics et privés.
type UserModel struct {
	Name  string `json:"name"` // exporté
	Age   int    `json:"age"`  // exporté
	id    int    // privé sans tag
	email string // privé sans tag
}

// AllPublicModel tous exportés - CONFORME.
// Structure avec uniquement des champs publics.
type AllPublicModel struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
}

// AllPrivateModel tous privés - CONFORME.
// Structure privée avec uniquement des champs privés.
type allPrivateModel struct {
	name    string
	age     int
	address string
}

// EmptyData struct vide - CONFORME.
// Structure vide utilisée comme marqueur.
type EmptyData struct{}

// OnlyPublicModel un seul champ exporté - CONFORME.
// Structure avec un seul champ public.
type OnlyPublicModel struct {
	Name string `json:"name"`
}

// onlyPrivateModel un seul champ privé - CONFORME.
// Structure privée avec un seul champ privé.
type onlyPrivateModel struct {
	name string
}

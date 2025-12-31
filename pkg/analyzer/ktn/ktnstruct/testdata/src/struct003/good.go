// Package struct003 provides good test cases.
// Ce fichier démontre les getters idiomatiques Go (sans préfixe Get).
package struct003

// GoodUser représente un utilisateur avec encapsulation correcte.
// Les getters suivent la convention Go idiomatique sans préfixe "Get".
type GoodUser struct {
	id    int
	name  string
	email string
}

// GoodUserInterface définit le contrat public de GoodUser.
type GoodUserInterface interface {
	Id() int
	Name() string
	Email() string
	SetName(name string)
}

// NewGoodUser crée une nouvelle instance de GoodUser.
//
// Params:
//   - id: identifiant unique de l'utilisateur
//   - name: nom de l'utilisateur
//   - email: adresse email de l'utilisateur
//
// Returns:
//   - *GoodUser: nouvelle instance initialisée
func NewGoodUser(id int, name, email string) *GoodUser {
	// Retourne une nouvelle instance avec les valeurs fournies
	return &GoodUser{
		id:    id,
		name:  name,
		email: email,
	}
}

// Id retourne l'identifiant de l'utilisateur.
//
// Returns:
//   - int: identifiant unique
func (u *GoodUser) Id() int {
	// Retourne le champ id
	return u.id
}

// Name retourne le nom de l'utilisateur.
//
// Returns:
//   - string: nom de l'utilisateur
func (u *GoodUser) Name() string {
	// Retourne le champ name
	return u.name
}

// Email retourne l'adresse email de l'utilisateur.
//
// Returns:
//   - string: adresse email
func (u *GoodUser) Email() string {
	// Retourne le champ email
	return u.email
}

// SetName définit le nom de l'utilisateur.
//
// Params:
//   - name: nouveau nom à définir
func (u *GoodUser) SetName(name string) {
	// Modifie le champ name
	u.name = name
}

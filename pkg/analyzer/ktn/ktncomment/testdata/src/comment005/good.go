// Package comment005 provides good test cases.
// This file demonstrates proper comment length.
package comment005

// User représente un utilisateur du système.
// Stocke les informations de base d'un utilisateur.
type User struct {
	name  string
	age   int
	email string
}

// UserInterface définit le contrat public de User.
type UserInterface interface {
	Name() string
	Age() int
	Email() string
}

// NewUser crée une nouvelle instance de User.
//
// Params:
//   - name: nom de l'utilisateur
//   - age: âge de l'utilisateur
//   - email: email de l'utilisateur
//
// Returns:
//   - *User: nouvelle instance
func NewUser(name string, age int, email string) *User {
	// Retour de la nouvelle instance
	return &User{name: name, age: age, email: email}
}

// Name retourne le nom de l'utilisateur.
//
// Returns:
//   - string: nom de l'utilisateur
func (u *User) Name() string {
	// Retour du nom
	return u.name
}

// Age retourne l'âge de l'utilisateur.
//
// Returns:
//   - int: âge de l'utilisateur
func (u *User) Age() int {
	// Retour de l'âge
	return u.age
}

// Email retourne l'email de l'utilisateur.
//
// Returns:
//   - string: email de l'utilisateur
func (u *User) Email() string {
	// Retour de l'email
	return u.email
}

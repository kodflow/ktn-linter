// Package struct006 contient les exemples de test pour KTN-STRUCT-004.
// Ce fichier démontre les getters non idiomatiques avec préfixe "Get".
package struct004

// BadUser représente un utilisateur avec getters non idiomatiques.
// Les getters utilisent le préfixe "Get" contrairement à la convention Go.
type BadUser struct {
	id    int
	name  string
	email string
}

// BadUserInterface définit le contrat public de BadUser.
type BadUserInterface interface {
	GetID() int
	GetName() string
	GetEmail() string
	Save() error
}

// NewBadUser crée une nouvelle instance de BadUser.
//
// Params:
//   - id: identifiant unique de l'utilisateur
//   - name: nom de l'utilisateur
//   - email: adresse email de l'utilisateur
//
// Returns:
//   - *BadUser: nouvelle instance initialisée
func NewBadUser(id int, name, email string) *BadUser {
	// Retourne une nouvelle instance avec les valeurs fournies
	return &BadUser{
		id:    id,
		name:  name,
		email: email,
	}
}

// GetID retourne l'identifiant de l'utilisateur.
// VIOLATION: devrait être ID() selon la convention Go.
//
// Returns:
//   - int: identifiant unique
func (u *BadUser) GetID() int { // want "KTN-STRUCT-004"
	// Retourne le champ id
	return u.id
}

// GetName retourne le nom de l'utilisateur.
// VIOLATION: devrait être Name() selon la convention Go.
//
// Returns:
//   - string: nom de l'utilisateur
func (u *BadUser) GetName() string { // want "KTN-STRUCT-004"
	// Retourne le champ name
	return u.name
}

// GetEmail retourne l'adresse email de l'utilisateur.
// VIOLATION: devrait être Email() selon la convention Go.
//
// Returns:
//   - string: adresse email
func (u *BadUser) GetEmail() string { // want "KTN-STRUCT-004"
	// Retourne le champ email
	return u.email
}

// Save sauvegarde l'utilisateur.
//
// Returns:
//   - error: erreur éventuelle
func (u *BadUser) Save() error {
	// Retourne nil si succès
	return nil
}

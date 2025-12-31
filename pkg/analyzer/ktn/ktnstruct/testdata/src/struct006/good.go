// Package struct006 provides good test cases.
// Demonstrates correct usage of unexported fields.
package struct006

// GoodUserDTO est un DTO avec champ privé sans tag sérialisable.
// Le champ privé counter n'a pas de tag = OK car c'est un détail interne.
type GoodUserDTO struct {
	id      int    // Champ privé pour l'identifiant
	name    string // Champ privé pour le nom
	counter int    // Champ privé sans tag = OK
}

// GoodUserDTOInterface définit le contrat public de GoodUserDTO.
type GoodUserDTOInterface interface {
	Id() int
	Name() string
	Counter() int
}

// NewGoodUserDTO crée une nouvelle instance de GoodUserDTO.
//
// Params:
//   - id: identifiant de l'utilisateur
//   - name: nom de l'utilisateur
//
// Returns:
//   - *GoodUserDTO: nouvelle instance
func NewGoodUserDTO(id int, name string) *GoodUserDTO {
	// Retour de la nouvelle instance
	return &GoodUserDTO{id: id, name: name, counter: 0}
}

// Id retourne l'identifiant de l'utilisateur.
//
// Returns:
//   - int: identifiant
func (u *GoodUserDTO) Id() int {
	// Retour de l'ID
	return u.id
}

// Name retourne le nom de l'utilisateur.
//
// Returns:
//   - string: nom
func (u *GoodUserDTO) Name() string {
	// Retour du nom
	return u.name
}

// Counter retourne le compteur interne.
//
// Returns:
//   - int: valeur du compteur
func (u *GoodUserDTO) Counter() int {
	// Retour du compteur
	return u.counter
}

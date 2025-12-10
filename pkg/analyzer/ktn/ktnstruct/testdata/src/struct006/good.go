// Good examples for the struct006 test case.
package struct006

// GoodUserDTO est un DTO avec seulement des champs publics tagués.
// Représente un utilisateur pour le transfert de données.
type GoodUserDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// GoodConfigRequest est un DTO correct.
// Contient les paramètres de configuration YAML.
type GoodConfigRequest struct {
	Timeout int    `yaml:"timeout"`
	Secret  string `yaml:"secret"`
}

// GoodResponseData est un DTO avec tous les champs publics.
// Représente la structure de réponse standard.
type GoodResponseData struct {
	Status  int    `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// GoodPrivateWithoutTag est un DTO avec champ privé sans tag.
// Démontre qu'un champ privé sans tag est conforme.
type GoodPrivateWithoutTag struct {
	ID      int `json:"id"`
	counter int // Pas de tag = OK
}

// goodRegularStruct est une struct privée non-DTO avec champs privés.
// Les structs privées ne nécessitent pas de getters.
type goodRegularStruct struct {
	name string
	age  int
}

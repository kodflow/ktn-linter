// Package struct006 contains test cases for KTN rules.
package struct006

// BadUserDTO est un DTO avec un champ privé tagué.
// Démontre la violation STRUCT-006 avec tag json sur champ privé.
type BadUserDTO struct {
	ID   int    `json:"id"`
	name string `json:"name"` // want "KTN-STRUCT-006: le champ privé 'name' du DTO 'BadUserDTO' ne devrait pas avoir de tag"
}

// BadConfigRequest est une struct avec suffixe DTO et champ privé tagué.
// Démontre la violation STRUCT-006 avec tag yaml sur champ privé.
type BadConfigRequest struct {
	Timeout int    `yaml:"timeout"`
	secret  string `yaml:"secret"` // want "KTN-STRUCT-006: le champ privé 'secret' du DTO 'BadConfigRequest' ne devrait pas avoir de tag"
}

// BadResponseData est un DTO avec plusieurs champs privés tagués.
// Démontre la violation STRUCT-006 avec plusieurs tags sur champs privés.
type BadResponseData struct {
	Status  int    `json:"status"`
	code    int    `json:"code"`    // want "KTN-STRUCT-006: le champ privé 'code' du DTO 'BadResponseData' ne devrait pas avoir de tag"
	message string `json:"message"` // want "KTN-STRUCT-006: le champ privé 'message' du DTO 'BadResponseData' ne devrait pas avoir de tag"
}

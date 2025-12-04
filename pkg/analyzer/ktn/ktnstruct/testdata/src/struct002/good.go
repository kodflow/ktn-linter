// Good examples for the struct004 test case.
package struct002

// User représente un utilisateur du système.
// Stocke les informations de base d'un utilisateur.
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// Config représente la configuration de l'application.
// Contient les paramètres de connexion au serveur.
type Config struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

// privateStruct n'a pas besoin de documentation car elle est privée
type privateStruct struct {
	data int `json:"-"`
}

// DataModel représente un modèle de données.
// Utilisé pour le transfert d'informations entre les couches.
type DataModel struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

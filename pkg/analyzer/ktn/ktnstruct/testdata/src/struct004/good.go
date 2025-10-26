package struct004

// User représente un utilisateur du système.
// Stocke les informations de base d'un utilisateur.
type User struct {
	Name string
	Age  int
}

// Config représente la configuration de l'application.
// Contient les paramètres de connexion au serveur.
type Config struct {
	Host string
	Port int
}

// privateStruct n'a pas besoin de documentation car elle est privée
type privateStruct struct {
	data int
}

// DataModel représente un modèle de données.
// Utilisé pour le transfert d'informations entre les couches.
type DataModel struct {
	ID   int
	Name string
}

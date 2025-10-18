package interface007 // want `\[KTN_INTERFACE_005\] Fichier interfaces.go existe mais ne contient aucune interface publique`

// Mauvais : interfaces.go vide ou sans interface publique
// Ce fichier existe mais ne contient pas d'interfaces publiques

// Ceci est une struct privée, pas une interface publique
type privateStruct struct {
	field string
}

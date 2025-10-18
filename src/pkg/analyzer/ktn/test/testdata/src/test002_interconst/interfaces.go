package test002_interconst

// Fichier interfaces.go avec GenDecl const (pas TYPE) - ne devrait pas être ignoré
// containsOnlyInterfaces002 continue à la ligne 96 car GenDecl.Tok != token.TYPE
const MaxSize = 100

var GlobalVar = 200

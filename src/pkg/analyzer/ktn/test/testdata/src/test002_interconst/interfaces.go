package test002_interconst

// Fichier interfaces.go avec GenDecl const (pas TYPE) - ne devrait pas être ignoré
// containsOnlyInterfaces002 continue à la ligne 96 car GenDecl.Tok != token.TYPE

// MaxSize defines the maximum size.
const MaxSize int = 100

// GlobalVar is a global constant.
const GlobalVar int = 200

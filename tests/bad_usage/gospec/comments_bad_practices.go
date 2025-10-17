// Package gospec_bad_comments montre des pratiques de documentation non-idiomatiques.
// Référence: https://go.dev/doc/effective_go
// Référence: https://github.com/golang/go/wiki/CodeReviewComments
package gospec_bad_comments

// ❌ BAD PRACTICE: No comment for exported function
func ProcessData() {}

// ❌ BAD PRACTICE: Comment doesn't start with name
// This processes the user input
func BadCommentStart() {}
// Devrait être: // BadCommentStart processes the user input

// ❌ BAD PRACTICE: Comment is just the signature repeated
// BadRedundant does BadRedundant
func BadRedundant() {}

// ❌ BAD PRACTICE: Comment doesn't add value
// GetUserName gets the user name
func GetUserName() string { return "" }
// Devrait expliquer COMMENT ou POURQUOI, pas répéter le nom

// ❌ BAD PRACTICE: Using // TODO without context
func BadTODO() {
	// TODO: fix this
	// Devrait inclure: qui, quoi, pourquoi
	// TODO(username): Fix race condition in user cache (see issue #123)
}

// ❌ BAD PRACTICE: Commented out code left in
func BadCommentedCode() {
	x := 42
	// y := 43
	// z := 44
	// fmt.Println(y, z)
	_ = x
}

// ❌ BAD PRACTICE: No package comment or incorrect format
// Le package comment devrait être:
// Package gospec_bad_comments provides examples.
// Ou:
/*
Package gospec_bad_comments provides examples of bad practices.

This package demonstrates...
*/

// ❌ BAD PRACTICE: Using /* */ for regular comments
func BadBlockComment() {
	/* This is a block comment */
	/* but single-line comments are preferred */
	x := 42
	_ = x
}

// ❌ BAD PRACTICE: Over-commenting obvious code
func BadOverComment() {
	// Declare variable x
	x := 0
	// Increment x
	x++
	// Check if x equals 1
	if x == 1 {
		// Print x
		print(x)
	}
}

// ❌ BAD PRACTICE: Comment doesn't match code
// This function returns true if valid
func BadMismatch() int { // Retourne int, pas bool!
	return 42
}

// ❌ BAD PRACTICE: No comment for exported type
type User struct {
	Name string
	Age  int
}

// ❌ BAD PRACTICE: No comment for exported field
type Config struct {
	Host string // Host is the host - redondant
	Port int    // No comment
}

// ❌ BAD PRACTICE: Comment with typos and poor grammar
// this function process user and return result without error if everything ok
func BadGrammar() {}

// ❌ BAD PRACTICE: Using non-standard comment markers
func BadMarkers() {
	// FIXME: should be TODO
	// HACK: should explain why
	// XXX: unclear meaning
	// NOTE: should be explained inline
}

// ❌ BAD PRACTICE: Magic number without explanation
func BadMagicNumber() {
	x := 86400 // Pas de commentaire expliquant pourquoi 86400
	_ = x
	// Devrait être: const secondsPerDay = 86400 ou commentaire
}

// ❌ BAD PRACTICE: Comment in wrong language (not English for public code)
// Cette fonction traite les données utilisateur
func BadLanguage() {}

// ❌ BAD PRACTICE: Outdated comment
// This function uses the old API endpoint
func BadOutdated() {
	// But actually uses new API now
	callNewAPI()
}

// ❌ BAD PRACTICE: Comment explaining WHY code is commented out
func BadWhyCommented() {
	// Don't use this, it's broken
	// useOldMethod()

	useNewMethod()
	// Si le code est cassé, supprimez-le avec Git pour historique
}

// ❌ BAD PRACTICE: No comment for complex algorithm
func BadComplexNoComment(items []int) int {
	// Algorithme complexe sans explication
	result := 0
	for i := 0; i < len(items); i++ {
		if i%2 == 0 {
			result += items[i] * 2
		} else {
			result -= items[i]
		}
	}
	return result
}

// ❌ BAD PRACTICE: Comment doesn't explain error conditions
func BadNoErrorDoc() error {
	// No documentation sur quand/pourquoi erreur retournée
	return validateInput()
}

// ❌ BAD PRACTICE: Constructor without doc explaining usage
func NewService() *Service {
	return &Service{}
}
// Devrait documenter: comment utiliser, options, thread-safety, etc.

// ❌ BAD PRACTICE: Interface without package comment explaining purpose
type Processor interface {
	Process()
}
// Devrait expliquer le contrat, expectations, thread-safety

// ❌ BAD PRACTICE: Exported const without explanation
const MaxRetries = 3

// ❌ BAD PRACTICE: Comment uses ASCII art that breaks on small screens
//     ┌──────────────────────────────────────────────────────────────┐
//     │  This is a very long ASCII art comment that might break      │
//     │  formatting or be hard to maintain                           │
//     └──────────────────────────────────────────────────────────────┘
func BadASCIIArt() {}

// ❌ BAD PRACTICE: Comment explains HOW instead of WHY
func BadHowNotWhy() {
	// Loop through items and multiply by 2
	// Comment évident - devrait expliquer POURQUOI on multiplie
	for i := 0; i < 10; i++ {
		_ = i * 2
	}
}

// BadDeprecated is deprecated
// ❌ BAD PRACTICE: Deprecation without alternative
func BadDeprecated() {}
// Devrait être: BadDeprecated is deprecated. Use NewFunction instead.

// Helper types and functions
type Service struct{}

func callNewAPI()      {}
func useNewMethod()    {}
func validateInput() error { return nil }

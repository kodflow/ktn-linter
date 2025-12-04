// Bad examples for the struct003 test case.
package struct003

// BadMixedFields représente une struct avec mauvais ordre des champs.
// Les champs exportés sont mélangés avec les champs privés (violation STRUCT-003).
type BadMixedFields struct {
	id        int
	Name      string // want "KTN-STRUCT-003"
	email     string
	Age       int // want "KTN-STRUCT-003"
	visible   bool
	Public    string // want "KTN-STRUCT-003"
	hidden    int
	Exported  bool // want "KTN-STRUCT-003"
	another   string
	LastField int // want "KTN-STRUCT-003"
}

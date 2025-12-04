// Bad examples for the struct003 test case.
package struct006

// BadMixedFields représente une struct avec mauvais ordre des champs.
// Les champs exportés sont mélangés avec les champs privés (violation STRUCT-003).
type BadMixedFields struct {
	id        int
	Name      string // want "KTN-STRUCT-006"
	email     string
	Age       int // want "KTN-STRUCT-006"
	visible   bool
	Public    string // want "KTN-STRUCT-006"
	hidden    int
	Exported  bool // want "KTN-STRUCT-006"
	another   string
	LastField int // want "KTN-STRUCT-006"
}

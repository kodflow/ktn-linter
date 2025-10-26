package struct003

// BadUser champs exportés après privés - VIOLATION
type BadUser struct {
	id    int
	Name  string // want "KTN-STRUCT-003"
	email string
	Age   int // want "KTN-STRUCT-003"
}

// MixedOrder ordre incorrect - VIOLATION
type MixedOrder struct {
	visible   bool
	Public    string // want "KTN-STRUCT-003"
	hidden    int
	Exported  bool // want "KTN-STRUCT-003"
	another   string
	LastOne   int // want "KTN-STRUCT-003"
}

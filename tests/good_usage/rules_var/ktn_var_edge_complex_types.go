package rules_var

// DataEntry représente une entrée de données.
type dataEntry struct {
	// Name nom de l'entrée
	Name string
	// Value valeur numérique
	Value int
}

// Variables du package.
var (
	// NestedMapGood déclaration de map imbriqué avec type explicite.
	NestedMapGood map[string]map[string][]int = map[string]map[string][]int{
		"level1": {
			"level2": {1, 2, 3},
		},
	}

	// SliceOfStructsGood déclaration de slice de structs avec type explicite.
	SliceOfStructsGood []dataEntry = []dataEntry{
		{"first", 1},
		{"second", 2},
	}

	// ChannelMapGood déclaration de channel de maps avec type explicite et buffer.
	ChannelMapGood chan map[string]int = make(chan map[string]int, 10)
)

package rules_var

// DataEntry représente une entrée de données.
type DataEntry struct {
	// Name nom de l'entrée
	Name string
	// Value valeur numérique
	Value int
}

var (
	// nestedMapGood déclaration de map imbriqué avec type explicite.
	nestedMapGood map[string]map[string][]int = map[string]map[string][]int{
		"level1": {
			"level2": {1, 2, 3},
		},
	}

	// sliceOfStructsGood déclaration de slice de structs avec type explicite.
	sliceOfStructsGood []DataEntry = []DataEntry{
		{"first", 1},
		{"second", 2},
	}

	// channelMapGood déclaration de channel de maps avec type explicite et buffer.
	channelMapGood chan map[string]int = make(chan map[string]int, 10)
)

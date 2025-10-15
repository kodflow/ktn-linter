package rules_var

// nestedMapBad déclaration de map imbriqué sans type explicite (VAR-001, VAR-004).
var nestedMapBad = map[string]map[string][]int{
	"level1": {
		"level2": {1, 2, 3},
	},
}

// sliceOfStructsBad déclaration de slice de structs sans type explicite (VAR-001, VAR-004).
var sliceOfStructsBad = []struct {
	Name  string
	Value int
}{
	{"first", 1},
	{"second", 2},
}

// channelMapBad déclaration de channel de maps sans type explicite (VAR-001, VAR-004, VAR-007).
var channelMapBad = make(chan map[string]int)

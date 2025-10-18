package alloc002

type Item struct {
	ID   int
	Name string
}

func BadMakeAppendPattern() {
	items := make([]Item, 0) // want `\[KTN-ALLOC-002\].*`
	for i := 0; i < 10; i++ {
		items = append(items, Item{ID: i})
	}
	_ = items
}

func BadMakeAppendPatternStrings() {
	names := make([]string, 0) // want `\[KTN-ALLOC-002\].*`
	names = append(names, "test")
	_ = names
}

func GoodPreallocWithCapacity() {
	items := make([]Item, 0, 10)
	for i := 0; i < 10; i++ {
		items = append(items, Item{ID: i})
	}
	_ = items
}

func GoodVarDecl() {
	var items []Item
	items = append(items, Item{ID: 1})
	_ = items
}

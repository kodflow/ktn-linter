package alloc002

type Item struct {
	ID   int
	Name string
}

func BadMakeAppendPattern() {
	// want `\[KTN-ALLOC-002\] Slice 'items' créé avec make\(\[\]T, 0\) puis utilisé avec append\(\)`
	items := make([]Item, 0)
	for i := 0; i < 10; i++ {
		items = append(items, Item{ID: i})
	}
	_ = items
}

func BadMakeAppendPatternStrings() {
	// want `\[KTN-ALLOC-002\] Slice 'names' créé avec make\(\[\]T, 0\) puis utilisé avec append\(\)`
	names := make([]string, 0)
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

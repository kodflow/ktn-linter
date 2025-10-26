package struct003

// User champs exportés avant privés - CONFORME
type User struct {
	Name  string // exporté
	Age   int    // exporté
	id    int    // privé
	email string // privé
}

// AllPublic tous exportés - CONFORME
type AllPublic struct {
	Name    string
	Age     int
	Address string
}

// AllPrivate tous privés - CONFORME
type AllPrivate struct {
	name    string
	age     int
	address string
}

// EmptyStruct struct vide - CONFORME
type EmptyStruct struct{}

// OnlyPublic un seul champ exporté - CONFORME
type OnlyPublic struct {
	Name string
}

// OnlyPrivate un seul champ privé - CONFORME
type OnlyPrivate struct {
	name string
}

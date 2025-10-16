// Package rules_alloc_good contient du code conforme aux règles KTN-ALLOC.
package rules_alloc_good

// ✅ Code conforme KTN-ALLOC-001 : make() pour types référence

// GoodMakeMapString crée une map avec make().
//
// Returns:
//   - map[string]int: map initialisée correctement
func GoodMakeMapString() map[string]int {
	m := make(map[string]int) // ✅ Correct : make() pour map
	m["key"] = 42
	// Retourne la map créée avec make
	return m
}

// GoodMakeMapWithCapacity crée une map avec capacité hint.
//
// Returns:
//   - map[int]string: map avec capacité pré-allouée
func GoodMakeMapWithCapacity() map[int]string {
	numbers := make(map[int]string, 100) // ✅ Correct : capacité spécifiée
	numbers[1] = "one"
	// Retourne la map avec capacité
	return numbers
}

// GoodMakeSlice crée un slice avec make().
//
// Returns:
//   - []int: slice initialisé correctement
func GoodMakeSlice() []int {
	s := make([]int, 0, 10) // ✅ Correct : make() avec capacité
	s = append(s, 1, 2, 3)
	// Retourne le slice créé avec make
	return s
}

// GoodMakeChannel crée un channel avec make().
//
// Returns:
//   - chan int: channel initialisé correctement
func GoodMakeChannel() chan int {
	ch := make(chan int, 5) // ✅ Correct : make() pour channel
	go func() {
		ch <- 42
	}()
	// Retourne le channel créé avec make
	return ch
}

// ✅ Code conforme KTN-ALLOC-002 : make() avec capacité avant append

// GoodMakeAppendWithCapacity fait des append avec capacité pré-allouée.
//
// Params:
//   - source: slice source à transformer
//
// Returns:
//   - []string: slice résultat
func GoodMakeAppendWithCapacity(source []string) []string {
	result := make([]string, 0, len(source)) // ✅ Correct : capacité = len(source)
	for _, v := range source {
		result = append(result, v)
	}
	// Retourne le résultat avec préallocation
	return result
}

// GoodMakeAppendKnownSize crée un slice avec taille connue.
//
// Returns:
//   - []int: slice de nombres doublés
func GoodMakeAppendKnownSize() []int {
	numbers := []int{1, 2, 3, 4, 5}
	doubled := make([]int, 0, len(numbers)) // ✅ Correct : capacité connue
	for _, n := range numbers {
		doubled = append(doubled, n*2)
	}
	// Retourne les nombres doublés
	return doubled
}

// GoodMakeAppendStruct utilise un slice de struct avec capacité.
//
// Returns:
//   - []User: slice d'utilisateurs
func GoodMakeAppendStruct() []User {
	users := make([]User, 0, 10) // ✅ Correct : capacité estimée
	users = append(users, User{ID: 1, Name: "Alice"})
	users = append(users, User{ID: 2, Name: "Bob"})
	// Retourne les utilisateurs
	return users
}

// GoodMakeLiteral crée un slice avec valeurs initiales via literal.
//
// Returns:
//   - []string: slice initialisé avec valeurs
func GoodMakeLiteral() []string {
	// ✅ Correct : composite literal, pas besoin de make
	items := []string{"a", "b", "c"}
	// Retourne le slice créé avec literal
	return items
}

// GoodMakePreallocated crée un slice avec taille initiale.
//
// Returns:
//   - []int: slice pré-alloué
func GoodMakePreallocated() []int {
	// ✅ Correct : make() avec longueur ET capacité
	data := make([]int, 10, 20)
	for i := range data {
		data[i] = i * i
	}
	// Retourne le slice pré-alloué
	return data
}

// ✅ Code conforme KTN-ALLOC-004 : &struct{} au lieu de new()

// User est une struct de test.
type User struct {
	// ID est l'identifiant unique de l'utilisateur.
	ID int
	// Name est le nom de l'utilisateur.
	Name string
	// Age est l'âge de l'utilisateur.
	Age int
}

// Config est une struct de configuration.
type Config struct {
	// Host est l'adresse du serveur.
	Host string
	// Port est le port du serveur.
	Port int
	// Timeout est le délai d'attente en secondes.
	Timeout int
}

// GoodNewUser crée un User avec composite literal.
//
// Returns:
//   - *User: pointeur vers User initialisé
func GoodNewUser() *User {
	u := &User{ // ✅ Correct : composite literal
		ID:   1,
		Name: "Alice",
		Age:  30,
	}
	// Retourne le pointeur vers User
	return u
}

// GoodNewUserEmpty crée un User vide avec composite literal.
//
// Returns:
//   - *User: pointeur vers User avec zero values
func GoodNewUserEmpty() *User {
	u := &User{} // ✅ Correct : composite literal vide
	u.ID = 1
	u.Name = "Bob"
	// Retourne le pointeur vers User
	return u
}

// GoodNewConfig crée une Config avec composite literal.
//
// Returns:
//   - *Config: pointeur vers Config initialisée
func GoodNewConfig() *Config {
	cfg := &Config{ // ✅ Correct : composite literal
		Host:    "localhost",
		Port:    8080,
		Timeout: 30,
	}
	// Retourne le pointeur vers Config
	return cfg
}

// GoodNewInline crée une struct inline avec composite literal.
//
// Returns:
//   - interface{}: pointeur vers Point
func GoodNewInline() interface{} {
	type Point struct {
		X, Y int
	}
	p := &Point{X: 10, Y: 20} // ✅ Correct : composite literal inline
	// Retourne le pointeur vers Point
	return p
}

// GoodMultipleStructs crée plusieurs structs avec composite literals.
//
// Returns:
//   - []*User: slice de pointeurs vers User
func GoodMultipleStructs() []*User {
	u1 := &User{ID: 1, Name: "Alice"} // ✅ Correct
	u2 := &User{ID: 2, Name: "Bob"}   // ✅ Correct
	// Retourne le slice d'utilisateurs
	return []*User{u1, u2}
}

// GoodStructAsReturn retourne directement un composite literal.
//
// Returns:
//   - *User: pointeur vers User créé inline
func GoodStructAsReturn() *User {
	// Retourne directement un composite literal (idiomatique)
	return &User{ID: 1, Name: "Charlie", Age: 25} // ✅ Correct
}

// GoodStructInLoop crée des structs dans une boucle avec composite literal.
//
// Returns:
//   - []*User: slice d'utilisateurs créés en boucle
func GoodStructInLoop() []*User {
	users := make([]*User, 0, 5)
	for i := 0; i < 5; i++ {
		u := &User{ID: i, Name: "User"} // ✅ Correct : composite literal
		users = append(users, u)
	}
	// Retourne les utilisateurs créés en boucle
	return users
}

// GoodMakeNoAppend crée un slice sans append (pas de violation ALLOC-002).
//
// Returns:
//   - []int: slice créé sans append
func GoodMakeNoAppend() []int {
	data := make([]int, 0) // ✅ Correct : pas d'append donc pas de problème
	// Retourne le slice vide sans append
	return data
}

// GoodMakeWithLengthNoAppend crée un slice avec longueur sans append.
//
// Returns:
//   - []string: slice pré-alloué avec longueur
func GoodMakeWithLengthNoAppend() []string {
	items := make([]string, 10) // ✅ Correct : longueur > 0, pas d'append
	items[0] = "first"
	items[9] = "last"
	// Retourne le slice pré-alloué
	return items
}

// GoodVariableLiteral crée un slice avec var et composite literal.
//
// Returns:
//   - []int: slice créé avec var
func GoodVariableLiteral() []int {
	var nums []int // ✅ Correct : nil slice, aucune allocation
	nums = []int{1, 2, 3}
	// Retourne le slice créé avec literal
	return nums
}

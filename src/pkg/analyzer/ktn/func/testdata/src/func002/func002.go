package func002

// goodFormat commence par le nom de la fonction.
func goodFormat() {}

// GoodExported commence par le nom.
func GoodExported() {}

// goodWithParams a un godoc valide.
func goodWithParams(x int) {}

func badNoComment() {} // want "KTN-FUNC-002.*sans commentaire godoc"

// Cette fonction ne commence pas par son nom. // want "KTN-FUNC-002.*doit commencer par le nom"
func badFormat() {}

// bad aussi ne commence pas par son nom. // want "KTN-FUNC-002.*doit commencer par le nom"
func badAnother() {}

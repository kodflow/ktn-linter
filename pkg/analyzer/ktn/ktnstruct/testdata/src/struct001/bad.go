package struct001

// BadProduct première struct dans le fichier
type BadProduct struct { // want "KTN-STRUCT-001"
	ID    int
	Price float64
}

// BadOrder deuxième struct dans le même fichier - violation
type BadOrder struct { // want "KTN-STRUCT-001"
	OrderID   int
	ProductID int
}

// BadCustomer troisième struct - violation encore plus grave
type BadCustomer struct { // want "KTN-STRUCT-001"
	Name  string
	Email string
}

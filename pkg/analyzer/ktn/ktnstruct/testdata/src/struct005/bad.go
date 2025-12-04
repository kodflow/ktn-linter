// Bad examples for the struct001 test case.
package struct005

// BadProduct représente un produit de test.
// Utilisé pour démontrer la violation de une struct par fichier.
type BadProduct struct { // want "KTN-STRUCT-005"
	ID    int
	Price float64
}

// BadOrder représente une commande de test.
// Démontre la violation avec une deuxième struct dans le même fichier.
type BadOrder struct { // want "KTN-STRUCT-005"
	OrderID   int
	ProductID int
}

// BadCustomer représente un client de test.
// Démontre la violation avec une troisième struct dans le même fichier.
type BadCustomer struct { // want "KTN-STRUCT-005"
	Name  string
	Email string
}

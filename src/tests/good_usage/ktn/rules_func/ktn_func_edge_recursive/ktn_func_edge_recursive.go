package rules_func

// Fonctions récursives correctement documentées.

// Factorial calcule la factorielle d'un nombre de manière récursive.
//
// Params:
//   - n: le nombre dont on calcule la factorielle
//
// Returns:
//   - int: la factorielle de n (n!)
func Factorial(n int) int {
	// Cas de base: factorielle de 0 ou 1 est 1
	if n <= 1 {
		// Retourne 1 car c'est le cas de base
		return 1
	}
	// Cas récursif: n! = n * (n-1)!
	// Retourne n multiplié par la factorielle de n-1
	return n * Factorial(n-1)
}

// Fibonacci calcule le n-ième nombre de Fibonacci récursivement.
//
// Params:
//   - n: l'indice du nombre de Fibonacci à calculer
//
// Returns:
//   - int: le n-ième nombre de Fibonacci
func Fibonacci(n int) int {
	// Cas de base: F(0) = 0, F(1) = 1
	if n <= 1 {
		// Retourne n car c'est le cas de base
		return n
	}
	// Cas récursif: F(n) = F(n-1) + F(n-2)
	// Retourne la somme des deux nombres de Fibonacci précédents
	return Fibonacci(n-1) + Fibonacci(n-2)
}

// TreeSearch recherche un nœud dans un arbre binaire de manière récursive.
//
// Params:
//   - node: le nœud racine de l'arbre ou sous-arbre
//   - target: la valeur recherchée
//
// Returns:
//   - *node: le nœud trouvé ou nil si non trouvé
func TreeSearch(node *node, target int) *node {
	// Cas de base: arbre vide
	if node == nil {
		// Retourne nil car l'arbre est vide
		return nil
	}

	// Cas de base: nœud trouvé
	if node.Value == target {
		// Retourne le nœud trouvé
		return node
	}

	// Recherche récursive dans le sous-arbre gauche
	if left := TreeSearch(node.Left, target); left != nil {
		// Retourne le nœud trouvé dans le sous-arbre gauche
		return left
	}

	// Recherche récursive dans le sous-arbre droit
	// Retourne le résultat de la recherche dans le sous-arbre droit
	return TreeSearch(node.Right, target)
}

// Node représente un nœud d'arbre binaire.
type node struct {
	// Value valeur stockée dans le nœud.
	Value int
	// Left enfant gauche.
	Left *node
	// Right enfant droit.
	Right *node
}

// MutuallyRecursive1 première fonction d'une paire mutuellement récursive.
//
// Params:
//   - n: le nombre à traiter
//
// Returns:
//   - int: résultat du calcul mutuel
func MutuallyRecursive1(n int) int {
	// Cas de base
	if n == 0 {
		// Retourne 0 car c'est le cas de base
		return 0
	}
	// Appel mutuel vers MutuallyRecursive2
	// Retourne le résultat de l'appel mutuel
	return MutuallyRecursive2(n - 1)
}

// MutuallyRecursive2 seconde fonction d'une paire mutuellement récursive.
//
// Params:
//   - n: le nombre à traiter
//
// Returns:
//   - int: résultat du calcul mutuel
func MutuallyRecursive2(n int) int {
	// Cas de base
	if n == 0 {
		// Retourne 1 car c'est le cas de base
		return 1
	}
	// Appel mutuel vers MutuallyRecursive1
	// Retourne le résultat de l'appel mutuel
	return MutuallyRecursive1(n - 1)
}

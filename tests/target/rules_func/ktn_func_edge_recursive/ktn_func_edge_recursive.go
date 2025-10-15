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
		return 1
	}
	// Cas récursif: n! = n * (n-1)!
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
		return n
	}
	// Cas récursif: F(n) = F(n-1) + F(n-2)
	return Fibonacci(n-1) + Fibonacci(n-2)
}

// TreeSearch recherche un nœud dans un arbre binaire de manière récursive.
//
// Params:
//   - node: le nœud racine de l'arbre ou sous-arbre
//   - target: la valeur recherchée
//
// Returns:
//   - *Node: le nœud trouvé ou nil si non trouvé
func TreeSearch(node *Node, target int) *Node {
	// Cas de base: arbre vide
	if node == nil {
		return nil
	}

	// Cas de base: nœud trouvé
	if node.Value == target {
		return node
	}

	// Recherche récursive dans le sous-arbre gauche
	if left := TreeSearch(node.Left, target); left != nil {
		return left
	}

	// Recherche récursive dans le sous-arbre droit
	return TreeSearch(node.Right, target)
}

// Node représente un nœud d'arbre binaire.
type Node struct {
	// Value valeur stockée dans le nœud.
	Value int
	// Left enfant gauche.
	Left *Node
	// Right enfant droit.
	Right *Node
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
		return 0
	}
	// Appel mutuel vers MutuallyRecursive2
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
		return 1
	}
	// Appel mutuel vers MutuallyRecursive1
	return MutuallyRecursive1(n - 1)
}

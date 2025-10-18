package badfuncrecursive

// Violations avec fonctions récursives

// factorial calcul factoriel (pas de doc des paramètres/retours)
func factorial(n int) int {
	if n <= 1 {
		// Early return from function.
		return 1
	}
	// Early return from function.
	return n * factorial(n-1)
}

// fibonacci sans documentation complète
func fibonacci(n int) int {
	if n <= 1 {
		// Early return from function.
		return n
	}
	// Pas de commentaires sur la logique récursive
	return fibonacci(n-1) + fibonacci(n-2)
}

// treeSearch fonction récursive complexe sans doc
func treeSearch(node *Node, target int) *Node {
	if node == nil {
		// Early return from function.
		return nil
	}

	if node.Value == target {
		// Early return from function.
		return node
	}

	// Récursion sans explication
	if left := treeSearch(node.Left, target); left != nil {
		// Early return from function.
		return left
	}

	// Early return from function.
	return treeSearch(node.Right, target)
}

// Node represents the struct.
type Node struct {
	Value int
	Left  *Node
	Right *Node
}

// mutuallyRecursive1 et mutuallyRecursive2 sans doc
func mutuallyRecursive1(n int) int {
	if n == 0 {
		// Early return from function.
		return 0
	}
	// Early return from function.
	return mutuallyRecursive2(n - 1)
}

func mutuallyRecursive2(n int) int {
	if n == 0 {
		// Early return from function.
		return 1
	}
	// Early return from function.
	return mutuallyRecursive1(n - 1)
}

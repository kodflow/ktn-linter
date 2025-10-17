package badfuncrecursive

// Violations avec fonctions récursives

// factorial calcul factoriel (pas de doc des paramètres/retours)
func factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

// fibonacci sans documentation complète
func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	// Pas de commentaires sur la logique récursive
	return fibonacci(n-1) + fibonacci(n-2)
}

// treeSearch fonction récursive complexe sans doc
func treeSearch(node *Node, target int) *Node {
	if node == nil {
		return nil
	}

	if node.Value == target {
		return node
	}

	// Récursion sans explication
	if left := treeSearch(node.Left, target); left != nil {
		return left
	}

	return treeSearch(node.Right, target)
}

type Node struct {
	Value int
	Left  *Node
	Right *Node
}

// mutuallyRecursive1 et mutuallyRecursive2 sans doc
func mutuallyRecursive1(n int) int {
	if n == 0 {
		return 0
	}
	return mutuallyRecursive2(n - 1)
}

func mutuallyRecursive2(n int) int {
	if n == 0 {
		return 1
	}
	return mutuallyRecursive1(n - 1)
}

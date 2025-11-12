// Fichier de test pour détecter les vraies fonctions mortes.
package func014_special

// deadFunction n'est jamais utilisée (ni appelée, ni callback).
func deadFunction() {
	// Code mort
}

// PublicFunction est exportée.
func PublicFunction() {
	// Fonction publique
}

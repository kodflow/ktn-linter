// Package var021 provides test cases for receiver consistency rule.
// This file contains non-compliant examples with inconsistent receivers.
package var021

// badServer représente un serveur avec des méthodes incohérentes.
// Cette structure a un mélange de receivers pointeur et valeur.
type badServer struct{ data int }

// Start démarre le serveur (receiver pointeur).
func (s *badServer) Start() {}

// Stop arrête le serveur (receiver valeur - incohérent). // want "KTN-VAR-021"
func (s badServer) Stop() {}

// Restart redémarre le serveur (receiver pointeur).
func (s *badServer) Restart() {}

// Status retourne le statut (receiver valeur - incohérent). // want "KTN-VAR-021"
func (s badServer) Status() int { return 0 }

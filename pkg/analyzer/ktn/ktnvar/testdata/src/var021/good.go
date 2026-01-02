// Package var021 provides test cases for receiver consistency rule.
// This file contains compliant examples with consistent receivers.
package var021

// goodServer représente un serveur avec receivers cohérents (tous pointeurs).
// Cette structure sert de test pour la règle KTN-VAR-021.
type goodServer struct{ data int }

// Start démarre le serveur.
func (s *goodServer) Start() {}

// Stop arrête le serveur.
func (s *goodServer) Stop() {}

// Restart redémarre le serveur.
func (s *goodServer) Restart() {}

// Status retourne le statut.
//
// Returns:
//   - int: données du serveur.
func (s *goodServer) Status() int { return s.data }

// point représente un point 2D avec receivers cohérents (tous valeurs).
// Cette structure est immuable par conception.
type point struct{ X, Y int }

// Add additionne deux points.
//
// Params:
//   - other: point à ajouter.
//
// Returns:
//   - point: nouveau point.
func (p point) Add(other point) point { return point{p.X + other.X, p.Y + other.Y} }

// String retourne une représentation textuelle.
//
// Returns:
//   - string: représentation du point.
func (p point) String() string { return "point" }

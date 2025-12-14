// Package messages provides structured error messages for KTN rules.
// This file contains COMMENT rule messages.
package messages

// registerCommentMessages enregistre les messages COMMENT.
func registerCommentMessages() {
	Register(Message{
		Code:  "KTN-COMMENT-001",
		Short: "commentaire trop long (%d > %d chars)",
		Verbose: `PROBLÈME: Le commentaire fait %d caractères (max %d).

POURQUOI: Les commentaires longs réduisent la lisibilité et cassent
le formatage dans les terminaux et éditeurs standards.

SOLUTIONS:
  1. Raccourcir le commentaire
  2. Passer en multi-ligne avec /* ... */
  3. Déplacer au-dessus de la ligne concernée

EXEMPLE INCORRECT:
  x := getValue() // Ce commentaire est beaucoup trop long

EXEMPLE CORRECT:
  // Récupère la valeur depuis le cache
  x := getValue()`,
	})

	Register(Message{
		Code:  "KTN-COMMENT-002",
		Short: "fichier sans commentaire descriptif avant 'package %s'",
		Verbose: `PROBLÈME: Le fichier n'a pas de commentaire avant package %s.

POURQUOI: Le commentaire de fichier documente le rôle du fichier.
Il apparaît dans godoc et aide à la navigation.

FORMAT ATTENDU:
  // Description courte du fichier.
  // Détails supplémentaires si nécessaire.
  package monpackage

EXEMPLE:
  // repository.go implémente l'accès à la base de données.
  // Il fournit les méthodes CRUD pour les entités User.
  package database`,
	})

	Register(Message{
		Code:  "KTN-COMMENT-003",
		Short: "constante '%s' sans commentaire",
		Verbose: `PROBLÈME: La constante '%s' n'a pas de commentaire explicatif.

POURQUOI: Les constantes définissent des valeurs métier importantes.
Sans documentation, leur signification se perd avec le temps.

FORMAT ATTENDU:
  // NomConstante définit/représente/indique...
  const NomConstante = valeur

EXEMPLE:
  // MaxRetries définit le nombre max de tentatives de connexion
  const MaxRetries = 3`,
	})

	Register(Message{
		Code:  "KTN-COMMENT-004",
		Short: "variable '%s' sans commentaire",
		Verbose: `PROBLÈME: La variable de package '%s' n'a pas de commentaire.

POURQUOI: Les variables de package ont une portée globale.
Leur documentation évite les mauvais usages.

FORMAT ATTENDU:
  // nomVariable stocke/contient/gère...
  var nomVariable Type = valeur

EXEMPLE:
  // defaultTimeout définit le délai d'attente par défaut.
  var defaultTimeout = 30 * time.Second`,
	})

	Register(Message{
		Code:  "KTN-COMMENT-005",
		Short: "struct '%s' sans documentation complète (≥2 lignes requises)",
		Verbose: `PROBLÈME: La struct exportée '%s' n'a pas de doc suffisante.

POURQUOI: Les structs exportées font partie de l'API publique.
Une documentation minimale (≥2 lignes) aide les utilisateurs.

FORMAT ATTENDU:
  // NomStruct représente/gère/contient...
  // Description détaillée du rôle et de l'utilisation.
  type NomStruct struct { ... }

EXEMPLE:
  // User représente un utilisateur du système.
  // Il contient les informations d'identification et de profil.
  type User struct { ... }`,
	})

	Register(Message{
		Code:  "KTN-COMMENT-006",
		Short: "section '%s' manquante dans la doc de '%s'",
		Verbose: `PROBLÈME: La fonction '%s' n'a pas la section '%s'.

POURQUOI: Une documentation structurée permet de comprendre
les entrées/sorties sans lire le code.

FORMAT ATTENDU:
  // NomFonction description courte.
  //
  // Params:
  //   - param1: description
  //
  // Returns:
  //   - Type: description
  func NomFonction(param1 Type) Type

EXEMPLE:
  // GetUserByID récupère un utilisateur par ID.
  //
  // Params:
  //   - id: identifiant unique
  //
  // Returns:
  //   - *User: utilisateur trouvé
  //   - error: ErrNotFound si inexistant
  func GetUserByID(id int) (*User, error)`,
	})

	Register(Message{
		Code:  "KTN-COMMENT-007",
		Short: "bloc '%s' sans commentaire explicatif",
		Verbose: `PROBLÈME: Le bloc '%s' n'a pas de commentaire.

POURQUOI: Les blocs de contrôle contiennent la logique métier.
Un commentaire aide à comprendre l'intention.

BLOCS CONCERNÉS: if, else, switch, case, for, return

EXEMPLE INCORRECT:
  if user.Age < 18 {
      return ErrUnderAge
  }

EXEMPLE CORRECT:
  // Vérifier que l'utilisateur est majeur
  if user.Age < 18 {
      // Mineur - accès refusé
      return ErrUnderAge
  }`,
	})
}

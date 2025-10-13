package rules_func

import (
	"context"
	"errors"
)

// ════════════════════════════════════════════════════════════════════════════
// KTN-FUNC-001 : Nom pas en MixedCaps/mixedCaps (snake_case interdit)
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : snake_case (SEULE ERREUR : KTN-FUNC-001)
// NOTE : Tout est parfait (commentaire + params OK) SAUF nom en snake_case
// ERREUR ATTENDUE : KTN-FUNC-001 sur parse_http_request

// parse_http_request parse une requête HTTP
func parse_http_request(data string) error {
	if data == "" {
		return errors.New("data vide")
	}
	return nil
}

// ❌ CAS INCORRECT 2 : Snake_Case mixte (SEULE ERREUR : KTN-FUNC-001)
// ERREUR ATTENDUE : KTN-FUNC-001 sur Calculate_Total

// Calculate_Total calcule le total
func Calculate_Total(values []int) int {
	total := 0
	for _, v := range values {
		total += v
	}
	return total
}

// ════════════════════════════════════════════════════════════════════════════
// KTN-FUNC-002 : Fonction exportée sans commentaire godoc
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Fonction exportée sans commentaire (SEULE ERREUR : KTN-FUNC-002)
// NOTE : Nom OK, params OK, longueur OK, MAIS pas de commentaire godoc
// ERREUR ATTENDUE : KTN-FUNC-002 sur ProcessOrder

func ProcessOrder(orderID int) error {
	if orderID <= 0 {
		return errors.New("orderID invalide")
	}
	return nil
}

// ❌ CAS INCORRECT 2 : Autre fonction exportée sans commentaire (SEULE ERREUR : KTN-FUNC-002)
// ERREUR ATTENDUE : KTN-FUNC-002 sur ValidateEmail

func ValidateEmail(email string) bool {
	return len(email) > 0
}

// ════════════════════════════════════════════════════════════════════════════
// KTN-FUNC-003 : Commentaire godoc incomplet - paramètres non documentés
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Params non documentés avec > 2 params (SEULE ERREUR : KTN-FUNC-003)
// NOTE : Tout est parfait (nom + commentaire présent + params OK) SAUF params non mentionnés dans doc
// ERREUR ATTENDUE : KTN-FUNC-003 sur CreateUser

// CreateUser crée un nouvel utilisateur
// Retourne l'ID du nouvel utilisateur ou une erreur
func CreateUser(name string, email string, age int) (int, error) {
	if name == "" {
		return 0, errors.New("name requis")
	}
	return 1, nil
}

// ════════════════════════════════════════════════════════════════════════════
// KTN-FUNC-004 : Commentaire godoc incomplet - retours non documentés
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Retours non documentés avec > 1 retour (SEULE ERREUR : KTN-FUNC-004)
// NOTE : Tout est parfait (nom + commentaire + params OK) SAUF retours non documentés
// ERREUR ATTENDUE : KTN-FUNC-004 sur FetchUserData

// FetchUserData récupère les données utilisateur depuis une source externe
func FetchUserData(ctx context.Context, userID int) (string, error) {
	if userID <= 0 {
		return "", errors.New("userID invalide")
	}
	return "user data", nil
}

// ════════════════════════════════════════════════════════════════════════════
// KTN-FUNC-005 : Trop de paramètres (> 5)
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : 6 paramètres (SEULE ERREUR : KTN-FUNC-005)
// NOTE : Tout est parfait (nom + commentaire + longueur OK) SAUF trop de params
// ERREUR ATTENDUE : KTN-FUNC-005 sur CreateUserAccount

// CreateUserAccount crée un nouveau compte utilisateur avec tous les détails
// Les paramètres name, email, age, address, phone et active sont tous requis pour la création
// Retourne l'ID du compte ou une erreur si les données sont invalides
func CreateUserAccount(name string, email string, age int, address string, phone string, active bool) (int, error) {
	if name == "" {
		return 0, errors.New("name requis")
	}
	return 1, nil
}

// ════════════════════════════════════════════════════════════════════════════
// KTN-FUNC-006 : Fonction trop longue (> 35 lignes)
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Fonction avec > 35 lignes (SEULE ERREUR : KTN-FUNC-006)
// NOTE : Tout est parfait (nom + commentaire + params + complexité OK) SAUF longueur
// ERREUR ATTENDUE : KTN-FUNC-006 sur ProcessLargeOrder

// ProcessLargeOrder traite une commande volumineuse avec de nombreuses étapes
// Le paramètre ctx contrôle le timeout et orderID identifie la commande
// Retourne une erreur si le traitement échoue
func ProcessLargeOrder(ctx context.Context, orderID int) error {
	// Simple séquence d'étapes sans conditions (complexité = 1)
	_ = ctx
	_ = orderID

	// Étape 1
	_ = "step 1"

	// Étape 2
	_ = "step 2"

	// Étape 3
	_ = "step 3"

	// Étape 4
	_ = "step 4"

	// Étape 5
	_ = "step 5"

	// Étape 6
	_ = "step 6"

	// Étape 7
	_ = "step 7"

	// Étape 8
	_ = "step 8"

	// Étape 9
	_ = "step 9"

	// Étape 10
	_ = "step 10"

	// Étape 11
	_ = "step 11"

	// Étape 12
	_ = "step 12"

	return nil
}

// ════════════════════════════════════════════════════════════════════════════
// KTN-FUNC-007 : Complexité cyclomatique trop élevée (≥ 10)
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Complexité ≥ 10 (SEULE ERREUR : KTN-FUNC-007)
// NOTE : Tout est parfait (nom + commentaire + params + longueur OK) SAUF complexité
// ERREUR ATTENDUE : KTN-FUNC-007 sur ValidateComplexInput

// ValidateComplexInput valide des données complexes avec de nombreuses conditions
// Les paramètres input et level contrôlent le niveau de validation
// Retourne une erreur si la validation échoue
func ValidateComplexInput(input string, level int) error {
	if input == "" {
		return errors.New("vide")
	}
	if level > 0 && len(input) < 3 {
		return errors.New("court")
	}
	if level > 1 && len(input) > 100 {
		return errors.New("long")
	}
	if level > 2 && input[0] == ' ' {
		return errors.New("espace début")
	}
	if level > 3 && input[len(input)-1] == ' ' {
		return errors.New("espace fin")
	}
	if level > 4 {
		for _, c := range input {
			if c == '\n' {
				return errors.New("newline")
			}
		}
	}
	if level > 5 {
		for _, c := range input {
			if c == '\t' {
				return errors.New("tab")
			}
		}
	}
	return nil
}

// ════════════════════════════════════════════════════════════════════════════
// KTN-FUNC-008 : Préfixe "Get" inutile
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Préfixe Get inutile (SEULE ERREUR : KTN-FUNC-008)
// NOTE : Tout est parfait (commentaire + params OK) SAUF préfixe Get
// ERREUR ATTENDUE : KTN-FUNC-008 sur GetUserName

// GetUserName retourne le nom de l'utilisateur
func GetUserName(userID int) string {
	if userID <= 0 {
		return ""
	}
	return "John Doe"
}

// ❌ CAS INCORRECT 2 : Autre Get inutile (SEULE ERREUR : KTN-FUNC-008)
// ERREUR ATTENDUE : KTN-FUNC-008 sur GetEmail

// GetEmail retourne l'email de l'utilisateur
func GetEmail(userID int) string {
	if userID <= 0 {
		return ""
	}
	return "user@example.com"
}

// ════════════════════════════════════════════════════════════════════════════
// KTN-FUNC-009 : Initialismes incorrects
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Http au lieu de HTTP (SEULE ERREUR : KTN-FUNC-009)
// NOTE : Tout est parfait (commentaire + params OK) SAUF initialisme incorrect
// ERREUR ATTENDUE : KTN-FUNC-009 sur ParseHttpRequest

// ParseHttpRequest parse une requête HTTP
func ParseHttpRequest(request string) error {
	if request == "" {
		return errors.New("requête vide")
	}
	return nil
}

// ❌ CAS INCORRECT 2 : Url au lieu de URL (SEULE ERREUR : KTN-FUNC-009)
// ERREUR ATTENDUE : KTN-FUNC-009 sur ValidateUrlFormat

// ValidateUrlFormat valide le format d'une URL
func ValidateUrlFormat(url string) bool {
	return len(url) > 0
}

// ❌ CAS INCORRECT 3 : Id au lieu de ID (SEULE ERREUR : KTN-FUNC-009)
// ERREUR ATTENDUE : KTN-FUNC-009 sur GenerateUserId

// GenerateUserId génère un ID unique pour un utilisateur
func GenerateUserId() int {
	return 42
}

// ════════════════════════════════════════════════════════════════════════════
// KTN-FUNC-010 : Context pas en premier paramètre
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Context en 2ème position (SEULE ERREUR : KTN-FUNC-010)
// NOTE : Tout est parfait (nom + commentaire + params OK) SAUF position Context
// ERREUR ATTENDUE : KTN-FUNC-010 sur FetchData

// FetchData récupère des données depuis une source externe
// Retourne une erreur si le contexte expire ou si l'ID est invalide
func FetchData(dataID int, ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		if dataID <= 0 {
			return "", errors.New("dataID invalide")
		}
		return "data", nil
	}
}

// ❌ CAS INCORRECT 2 : Context en dernière position (SEULE ERREUR : KTN-FUNC-010)
// ERREUR ATTENDUE : KTN-FUNC-010 sur ProcessRequest

// ProcessRequest traite une requête avec timeout
// Les paramètres requestID, data et ctx contrôlent le traitement
// Retourne une erreur si le traitement échoue
func ProcessRequest(requestID int, data string, ctx context.Context) error {
	_ = ctx
	if requestID <= 0 {
		return errors.New("requestID invalide")
	}
	return nil
}

// ════════════════════════════════════════════════════════════════════════════
// KTN-TEST-003: Fichier _test.go orphelin (sans .go correspondant)
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Un fichier *_test.go ne peut pas exister sans son fichier .go correspondant.
//    Relation 1:1 : orphan_test.go requiert orphan.go
//
//    POURQUOI :
//    - Maintient la cohérence de la structure
//    - Évite les tests qui testent du code supprimé/renommé
//    - Facilite la maintenance (pas de fichiers orphelins)
//    - Détecte les erreurs de refactoring
//
// ❌ CAS INCORRECT : Ce fichier orphan_test.go existe mais pas orphan.go
// ERREUR ATTENDUE: KTN-TEST-003 sur orphan_test.go
//
// ✅ CAS PARFAIT (voir target/) :
//    helper.go + helper_test.go (les deux existent)
//
// ════════════════════════════════════════════════════════════════════════════
package KTN_TEST_003_test

import "testing"

// TestOrphan teste quelque chose qui n'existe pas.
func TestOrphan(t *testing.T) {
	// Ce test est orphelin car orphan.go n'existe pas
	t.Log("Test orphelin")
}

// NOTE: orphan.go n'existe PAS (violation KTN-TEST-003)

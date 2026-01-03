// Package shared provides common utilities for static analysis.
package shared

import (
	"go/ast"
	"go/token"
)

// DeclGroup représente un groupe de déclarations (const ou var).
// Utilisé par ktnconst/002.go et ktnvar/002.go pour vérifier le groupement.
type DeclGroup struct {
	Decl *ast.GenDecl
	Pos  token.Pos
}

// IsSerializableStruct vérifie si une struct est un DTO/serializable.
// Une struct est considérée comme DTO uniquement si:
//   - Elle a AU MOINS UN champ avec un tag de sérialisation (json, yaml, xml, etc.)
//
// Le nom seul (suffixe DTO) ne suffit pas pour éviter les faux positifs
// sur des structs internes qui ne sont jamais sérialisées.
//
// Params:
//   - structType: type de la struct à vérifier
//   - structName: nom de la struct (ignoré, conservé pour compatibilité)
//
// Returns:
//   - bool: true si c'est un DTO/serializable
func IsSerializableStruct(structType *ast.StructType, _ string) bool {
	return hasSerializationTags(structType)
}

// hasSerializationTags vérifie si la struct a des tags (backticks).
// Tout champ avec un tag (json, yaml, xml, gorm, etc.) indique un DTO.
//
// Params:
//   - structType: type de la struct
//
// Returns:
//   - bool: true si des tags sont présents
func hasSerializationTags(structType *ast.StructType) bool {
	// Pas de champs
	if structType.Fields == nil {
		// Retour si pas de champs
		return false
	}

	// Parcourir les champs
	for _, field := range structType.Fields.List {
		// Vérifier si le champ a un tag (backticks)
		if field.Tag != nil && field.Tag.Value != "" {
			// Tag trouvé = DTO
			return true
		}
	}

	// Pas de tag
	return false
}

// IsPureDataStruct vérifie si une struct est un pur conteneur de données.
// Une struct est considérée comme pure data si:
//   - Tous ses champs sont publics (ou tous privés)
//   - Elle n'a pas de méthodes (à part éventuellement String())
//
// Params:
//   - structType: type de la struct
//
// Returns:
//   - bool: true si c'est un pur conteneur de données
func IsPureDataStruct(structType *ast.StructType) bool {
	// Pas de champs = pas une struct de données
	if structType.Fields == nil || len(structType.Fields.List) == 0 {
		// Retour si pas de champs
		return false
	}

	// Vérifier que tous les champs sont publics
	allPublic := true
	// Parcourir les champs
	for _, field := range structType.Fields.List {
		// Parcourir les noms (un champ peut avoir plusieurs noms)
		for _, name := range field.Names {
			// Vérifier si le champ est privé
			if !ast.IsExported(name.Name) {
				allPublic = false
				break
			}
		}
	}

	// Pure data si tous les champs sont publics
	return allPublic
}

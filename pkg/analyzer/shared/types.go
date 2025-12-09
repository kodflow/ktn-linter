// Shared utilities for types handling.
package shared

import (
	"go/ast"
	"go/token"
	"strings"
)

// DeclGroup représente un groupe de déclarations (const ou var).
// Utilisé par ktnconst/002.go et ktnvar/002.go pour vérifier le groupement.
type DeclGroup struct {
	Decl *ast.GenDecl
	Pos  token.Pos
}

// dtoSuffixes contient les suffixes typiques des DTOs.
var dtoSuffixes []string = []string{
	"Config",
	"Settings",
	"Options",
	"Params",
	"Request",
	"Response",
	"DTO",
	"Model",
	"Entity",
	"Spec",
	"Schema",
	"Payload",
	"Input",
	"Output",
	"Args",
	"Result",
	"Data",
	"Info",
	"Details",
	"State",
	"Status",
	"Metadata",
	"Context",
	"Event",
	"Message",
	"Command",
	"Query",
}

// IsSerializableStruct vérifie si une struct est un DTO/serializable.
// Une struct est considérée comme DTO si:
//   - Elle a des champs avec des tags de sérialisation (json, yaml, xml, etc.)
//   - Ou son nom se termine par un suffixe DTO typique (Config, Request, etc.)
//
// Params:
//   - structType: type de la struct à vérifier
//   - structName: nom de la struct
//
// Returns:
//   - bool: true si c'est un DTO/serializable
func IsSerializableStruct(structType *ast.StructType, structName string) bool {
	// Vérifier par le nom de la struct
	if hasSerializableSuffix(structName) {
		// Retour si suffixe DTO trouvé
		return true
	}

	// Vérifier par les tags des champs
	// Retour résultat de la vérification par tags
	return hasSerializationTags(structType)
}

// hasSerializableSuffix vérifie si le nom a un suffixe DTO.
//
// Params:
//   - name: nom de la struct
//
// Returns:
//   - bool: true si suffixe DTO trouvé
func hasSerializableSuffix(name string) bool {
	// Parcourir les suffixes
	for _, suffix := range dtoSuffixes {
		// Vérifier le suffixe
		if strings.HasSuffix(name, suffix) {
			// Retour si suffixe DTO trouvé
			return true
		}
	}
	// Pas de suffixe DTO
	return false
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

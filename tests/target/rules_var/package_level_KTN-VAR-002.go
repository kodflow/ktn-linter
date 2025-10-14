package rules_var

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-002 : Groupe avec commentaire de groupe
// ════════════════════════════════════════════════════════════════════════════

// Application metadata
// Ces variables contiennent les métadonnées (peuvent être modifiées à runtime)
var (
	// ApplicationNameV002Good est le nom de l'application
	ApplicationNameV002Good string = "MyApp"
	// VersionV002Good est la version actuelle de l'application
	VersionV002Good string = "1.0.0"
	// defaultEncodingV002Good est l'encodage par défaut utilisé
	defaultEncodingV002Good string = "UTF-8"
)

// Disk and time values
// Ces variables utilisent int64 pour les grandes valeurs
var (
	// MaxDiskSpaceV002Good est l'espace disque maximum en octets
	MaxDiskSpaceV002Good int64 = 1099511627776
	// UnixEpochV002Good représente le timestamp Unix epoch (intentionnellement 0)
	UnixEpochV002Good int64 = 0
	// nanosPerSecondV002Good est le nombre de nanosecondes par seconde
	nanosPerSecondV002Good int64 = 1000000000
)

// updateMetadataV002Good modifie les métadonnées à runtime
func updateMetadataV002Good() {
	ApplicationNameV002Good = "UpdatedApp"
	VersionV002Good = "2.0.0"
	defaultEncodingV002Good = "UTF-16"
	MaxDiskSpaceV002Good = 2199023255552
	UnixEpochV002Good = 1234567890
	nanosPerSecondV002Good = 2000000000
}

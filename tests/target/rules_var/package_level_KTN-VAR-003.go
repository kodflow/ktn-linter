package rules_var

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-003 : Variable avec commentaire individuel
// ════════════════════════════════════════════════════════════════════════════

// Age limits
// Ces variables définissent les limites d'âge (configurables)
var (
	// MinAgeV003Good est l'âge minimum autorisé
	MinAgeV003Good int8 = 18
	// MaxAgeV003Good est l'âge maximum autorisé
	MaxAgeV003Good int8 = 120
	// defaultPriorityV003Good est la priorité par défaut
	defaultPriorityV003Good int8 = 5
)

// updateLimitsV003Good modifie les limites à runtime
func updateLimitsV003Good() {
	MinAgeV003Good = 21
	MaxAgeV003Good = 100
	defaultPriorityV003Good = 10
}

package rules_var

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-006 : Une variable par ligne (pas de déclaration multiple)
// ════════════════════════════════════════════════════════════════════════════

// Network settings
// Ces variables configurent les paramètres réseau
var (
	// HostNameV006Good est le nom d'hôte par défaut
	HostNameV006Good string = "localhost"
	// PortV006Good est le port réseau par défaut
	PortV006Good int = 8080
)

// updateNetworkV006Good modifie les paramètres réseau à runtime
func updateNetworkV006Good() {
	HostNameV006Good = "example.com"
	PortV006Good = 9090
}

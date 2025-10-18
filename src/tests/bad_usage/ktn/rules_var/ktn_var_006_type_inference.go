package rules_var

// ════════════════════════════════════════════════════════════════════════════
// KTN-VAR-006 : Multiple variables sur une ligne
// ════════════════════════════════════════════════════════════════════════════
// 📝 DESCRIPTION DU PROBLÈME :
//    Plusieurs variables déclarées sur une ligne (HostName, Port = "localhost", 8080)
//    rendent impossible la documentation individuelle de chaque variable.
//
//    POURQUOI :
//    - Impossible de mettre un commentaire par variable
//    - Difficile à lire et à maintenir
//    - Contraire aux bonnes pratiques de documentation
//
// ✅ CAS PARFAIT (une variable par ligne) :
//
//    // Network settings
//    // Ces variables configurent la connexion réseau
//    var (
//        // HostName est le nom d'hôte
//        HostName string = "localhost"
//        // Port est le port réseau
//        Port int = 8080
//    )
//
// ════════════════════════════════════════════════════════════════════════════

// ❌ CAS INCORRECT 1 : Multiple variables sur une ligne
// ERREURS : KTN-VAR-006 sur HostNameV006, PortV006
// Network settings
var (
	// HostNameV006 et PortV006 sont les paramètres réseau
	HostNameV006, PortV006 = "localhost", 8080
)

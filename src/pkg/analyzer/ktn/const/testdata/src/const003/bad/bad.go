package bad

// Status constants - group comment is present
const (
	StatusActive   = "active"   // want `Constante 'StatusActive' sans commentaire individuel`
	StatusInactive = "inactive" // want `Constante 'StatusInactive' sans commentaire individuel`
)

// Configuration - group comment is present
const (
	MaxRetries = 3  // want `Constante 'MaxRetries' sans commentaire individuel`
	Timeout    = 30 // want `Constante 'Timeout' sans commentaire individuel`
)

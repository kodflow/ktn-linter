package interface002

// Mauvais : struct publique au lieu d'interface
type PublicService struct { // want `\[KTN_INTERFACE_002\] Type public 'PublicService' défini comme struct au lieu d'interface`
	name string
}

// Mauvais : autre struct publique
type DataProcessor struct { // want `\[KTN_INTERFACE_002\] Type public 'DataProcessor' défini comme struct au lieu d'interface`
	config string
}

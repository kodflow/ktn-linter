package test011

// PublicFunc est une fonction publique
func PublicFunc() string {
	return privateHelper()
}

// privateHelper est une fonction privée
func privateHelper() string {
	return "helper"
}

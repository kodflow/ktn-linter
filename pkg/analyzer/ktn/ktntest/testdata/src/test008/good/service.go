package test010

// PublicService est une fonction publique
func PublicService() string {
	return privateImplementation()
}

// privateImplementation est une fonction privée
func privateImplementation() string {
	return "impl"
}

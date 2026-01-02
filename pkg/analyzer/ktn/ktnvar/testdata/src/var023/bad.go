package var023

import "math/rand"

func badGenerateKey() int {
	return rand.Intn(1000000) // want "KTN-VAR-023"
}

func badCreateToken() int64 {
	token := rand.Int63() // want "KTN-VAR-023"
	return token
}

var badSecretKey = rand.Uint64() // want "KTN-VAR-023"

func badGenerateSecret() int {
	var secretValue = rand.Intn(100) // want "KTN-VAR-023"
	return secretValue
}

func badNormalFunction() int {
	var tokenValue = rand.Intn(100) // want "KTN-VAR-023" "KTN-VAR-023"
	return tokenValue
}

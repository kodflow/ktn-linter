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

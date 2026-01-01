package var023

import (
	"crypto/rand"
	"math/big"
	mathrand "math/rand"
)

func goodGenerateKey() *big.Int {
	key, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return key
}

func goodShuffleItems(items []int) {
	mathrand.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	}) // OK - not security context
}

func goodRandomIndex() int {
	return mathrand.Intn(100) // OK - not security context
}

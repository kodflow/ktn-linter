package var018

import "bytes"

// badRepeatedConversionInLoop convertit plusieurs fois dans une boucle.
func badRepeatedConversionInLoop(data [][]byte, target string) int {
	count := 0
	for _, item := range data { // want "KTN-VAR-018: conversion string\\(\\) répétée dans la boucle, préallouer hors de la boucle"
		if string(item) == target {
			count++
		}
	}
	return count
}

// badMultipleConversionsInFunction convertit plusieurs fois la même variable.
func badMultipleConversionsInFunction(data []byte) { // want "KTN-VAR-018: conversion string\\(\\) de 'data' répétée 3 fois, stocker dans une variable"
	if string(data) == "hello" {
		println("found hello")
	}
	if string(data) == "world" {
		println("found world")
	}
	println(string(data))
}

// badConversionInForLoop convertit dans un for classique.
func badConversionInForLoop(items [][]byte) { // want "KTN-VAR-018: conversion string\\(\\) répétée dans la boucle, préallouer hors de la boucle"
	for i := 0; i < len(items); i++ {
		if string(items[i]) == "test" {
			println("found")
		}
	}
}

// badNestedLoopConversion convertit dans une boucle imbriquée.
func badNestedLoopConversion(row [][]byte) { // want "KTN-VAR-018: conversion string\\(\\) répétée dans la boucle, préallouer hors de la boucle"
	for _, cell := range row {
		if string(cell) == "x" {
			println("found x")
		}
	}
}

// badMapKeyConversion utilise string() répété pour les clés de map.
func badMapKeyConversion(cache map[string]int, keys [][]byte) { // want "KTN-VAR-018: conversion string\\(\\) répétée dans la boucle, préallouer hors de la boucle"
	for _, key := range keys {
		_ = cache[string(key)]
	}
}

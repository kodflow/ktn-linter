package rules_type_ops

// ✅ GOOD: pas de conversion redondante
func noRedundantConversion() {
	var x int = 42
	y := x // ✅ pas de conversion inutile
	println(y)
}

// ✅ GOOD: conversion nécessaire
func necessaryConversion() {
	var x int32 = 42
	y := int(x) // ✅ int32 -> int nécessaire
	println(y)
}

func necessaryFloat() {
	var i int = 42
	f := float64(i) // ✅ int -> float64 nécessaire
	println(f)
}

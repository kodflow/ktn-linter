// Package interface003 contains test cases for KTN-INTERFACE-003.
package interface003

// Doable is a single-method interface without -er convention. // want "KTN-INTERFACE-003"
type Doable interface {
	Do()
}

// Runnable is a single-method interface without -er convention. // want "KTN-INTERFACE-003"
type Runnable interface {
	Run()
}

// WriteThing is a single-method interface without -er convention. // want "KTN-INTERFACE-003"
type WriteThing interface {
	Write()
}

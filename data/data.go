package data

type Filter int

const (
	// FilterLt represents strictly Less than
	FilterLt Filter = iota
	// FilterLe represents Less than or equal
	FilterLe
	// FilterEq represents equal
	FilterEq
	// FilterGe represent greater than or equal
	FilterGe
	// FilterGt represents strictly greater than
	FilterGt
)

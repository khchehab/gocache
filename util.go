package gocache

import "github.com/DmitriyVTitov/size"

// SizeOf returns the size of a variable of any type.
//
// Parameter:
//   - v: type any variable.
//
// Returns:
//   - uint: The size of the variable in bytes.
func SizeOf(v any) uint64 {
	return uint64(size.Of(v))
}

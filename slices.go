package output

import "slices"

// FilledStrings returns a slice of n strings, each filled with the given value.
func FilledStrings(n int, value string) []string {
	return slices.Repeat([]string{value}, n)
}

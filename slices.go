package output

// FilledStrings returns a slice of n strings, each filled with the given value.
func FilledStrings(n int, value string) []string {
	slice := make([]string, n)
	for i := range slice {
		slice[i] = value
	}

	return slice
}

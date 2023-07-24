package helpers

import "strconv"

// IntToString converts an int to a string.
func IntToString(value int) string {
	return strconv.Itoa(value)
}

// ValueOrDefault handles nil values for strings.
func ValueOrDefault(value *string) string {
	if value == nil {
		return "" // Set it to an empty string if it is nil
	}
	return *value // Otherwise, return the value of the pointer
}

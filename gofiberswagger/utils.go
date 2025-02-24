package gofiberswagger

import "strings"

/// ------------------------------------------------------ ///
/// Function that are usefull for internal & general usage ///
/// ------------------------------------------------------ ///

func replaceNthOccurrence(s, old, new string, n int) string {
	parts := strings.Split(s, old)
	if n <= 0 || n >= len(parts) {
		return s
	}
	result := strings.Join(parts[:n], old) + new + strings.Join(parts[n:], old)
	return result
}

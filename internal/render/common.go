package render

import "strings"

func Spaces(count int) string {
	return strings.Repeat(" ", count)
}

func Underscores(count int) string {
	return strings.Repeat("_", count)
}

// Nlnl (NoLeadingNewline) removes the leading newline from a string.
func Nlnl(s string) string {
	return strings.TrimLeft(s, "\n")
}

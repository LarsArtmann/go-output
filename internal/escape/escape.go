// Package escape provides string escaping utilities for various output formats.
package escape

import "strings"

type mode struct {
	apos string
}

var (
	htmlMode = mode{apos: "&#39;"}
	xmlMode  = mode{apos: "&apos;"}
)

// HTML escapes HTML special characters.
func HTML(s string) string {
	return escape(s, htmlMode)
}

// XML escapes XML special characters.
func XML(s string) string {
	return escape(s, xmlMode)
}

func escape(s string, m mode) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '&':
			b.WriteString("&amp;")
		case '"':
			b.WriteString("&quot;")
		case '\'':
			b.WriteString(m.apos)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

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

// D2 escapes special characters for D2 diagram strings.
func D2(s string) string {
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\t", `\t`)

	return s
}

// DOT escapes special characters for DOT/Graphviz strings.
func DOT(s string) string {
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)

	return s
}

// MermaidID sanitizes a string for use as a Mermaid node identifier.
func MermaidID(id string) string {
	var result strings.Builder

	for _, r := range id {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '_' {
			result.WriteRune(r)
		}
	}

	if result.Len() == 0 {
		return "node"
	}

	return result.String()
}

// MermaidSlug sanitizes a string for use as a Mermaid node identifier fallback.
func MermaidSlug(label string) string {
	s := strings.ReplaceAll(label, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, "/", "_")

	return s
}

// MermaidText escapes special characters for Mermaid display labels.
func MermaidText(s string) string {
	s = strings.ReplaceAll(s, `"`, "'")
	s = strings.ReplaceAll(s, "[", "(")
	s = strings.ReplaceAll(s, "]", ")")
	s = strings.ReplaceAll(s, "{", "(")
	s = strings.ReplaceAll(s, "}", ")")
	s = strings.ReplaceAll(s, "\n", "<br>")

	return s
}

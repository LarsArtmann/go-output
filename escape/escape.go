// Package escape provides string escaping utilities for various output formats.
package escape

import (
	"html"
	"strings"
)

// XML escapes XML special characters.
// Uses &apos; for apostrophe (canonical XML entity) instead of &#39; from html.EscapeString.
func XML(s string) string {
	return strings.ReplaceAll(html.EscapeString(s), "&#39;", "&apos;")
}

// HTML escapes HTML special characters using html.EscapeString from the standard library.
func HTML(s string) string {
	return html.EscapeString(s)
}

// D2 escapes special characters for D2 diagram strings.
func D2(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\t", `\t`)

	return s
}

// DOT escapes special characters for DOT/Graphviz strings.
// DOT and D2 share the same escaping rules: backslash, double quote, newline, and tab.
func DOT(s string) string {
	return D2(s)
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

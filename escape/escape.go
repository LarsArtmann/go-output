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

// d2Replacer escapes backslash, double quote, newline, and tab for D2 strings.
//
//nolint:gochecknoglobals // Reusable strings.Replacer, safe to share.
var d2Replacer = strings.NewReplacer(
	`\`, `\\`,
	`"`, `\"`,
	"\n", `\n`,
	"\t", `\t`,
)

// D2 escapes special characters for D2 diagram strings.
func D2(s string) string {
	return d2Replacer.Replace(s)
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
		if isMermaidIdentRune(r) {
			result.WriteRune(r)
		}
	}

	if result.Len() == 0 {
		return "node"
	}

	return result.String()
}

// isMermaidIdentRune reports whether r is valid in a Mermaid node identifier
// (ASCII letter, digit, or underscore). Centralized so fuzz tests can verify
// MermaidID's output invariant without duplicating the predicate.
func isMermaidIdentRune(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') || r == '_'
}

// slugIDReplacer sanitizes strings for use as identifiers across diagram formats.
// Replaces identifier-hostile characters with underscores.
//
//nolint:gochecknoglobals // Reusable strings.Replacer, safe to share.
var slugIDReplacer = strings.NewReplacer(
	" ", "_",
	"-", "_",
	"/", "_",
	".", "_",
	"*", "_",
	"[", "_",
	"]", "_",
	"{", "_",
	"}", "_",
	"(", "_",
	")", "_",
)

// SlugifyID sanitizes a string for use as a diagram node identifier.
// It replaces spaces, hyphens, slashes, dots, asterisks, and brackets with
// underscores — characters that are problematic in D2, DOT, Mermaid, and
// PlantUML identifiers.
func SlugifyID(s string) string {
	return slugIDReplacer.Replace(s)
}

// MermaidSlug sanitizes a string for use as a Mermaid node identifier fallback.
func MermaidSlug(label string) string {
	return SlugifyID(label)
}

// mermaidTextReplacer escapes brackets, braces, quotes, and newlines for Mermaid labels.
//
//nolint:gochecknoglobals // Reusable strings.Replacer, safe to share.
var mermaidTextReplacer = strings.NewReplacer(
	`"`, "'",
	"[", "(",
	"]", ")",
	"{", "(",
	"}", ")",
	"\n", "<br>",
)

// MermaidText escapes special characters for Mermaid display labels.
func MermaidText(s string) string {
	return mermaidTextReplacer.Replace(s)
}

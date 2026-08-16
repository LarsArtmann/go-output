// Package escape provides string escaping utilities for various output formats.
package escape

import (
	"html"
	"regexp"
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

// d2Replacer escapes backslash, double quote, newline, carriage return, and
// tab for D2 strings.
//
//nolint:gochecknoglobals // Reusable strings.Replacer, safe to share.
var d2Replacer = strings.NewReplacer(
	`\`, `\\`,
	`"`, `\"`,
	"\n", `\n`,
	"\r", `\r`,
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

// fallbackID is returned by MermaidID/PlantUMLID when sanitization strips
// every rune — diagram renderers need a non-empty identifier.
const fallbackID = "node"

// MermaidID sanitizes a string for use as a Mermaid node identifier.
func MermaidID(id string) string {
	var result strings.Builder

	for _, r := range id {
		if isMermaidIdentRune(r) {
			result.WriteRune(r)
		}
	}

	if result.Len() == 0 {
		return fallbackID
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

// mermaidTextReplacer escapes HTML-significant characters, brackets, braces,
// quotes, and newlines for Mermaid labels. Mermaid renders labels as HTML by
// default (htmlLabels), so raw `&`, `<`, and `>` in node labels, edge labels,
// and style values would be interpreted as markup (`<script>`, `<img
// onerror=...>`), not text.
//
//nolint:gochecknoglobals // Reusable strings.Replacer, safe to share.
var mermaidTextReplacer = strings.NewReplacer(
	`&`, "&amp;",
	`<`, "&lt;",
	`>`, "&gt;",
	`"`, "'",
	"[", "(",
	"]", ")",
	"{", "(",
	"}", ")",
	"\n", "<br>",
)

// mermaidEntityGuard neutralizes Mermaid entity codes (`#60;`, `#quot;`):
// Mermaid decodes `#…;` sequences inside labels into raw characters before the
// HTML render, which would re-manufacture the `<`, `>`, and `&` characters the
// replacer above already escaped. Prefixing the `#` with `&#35;` renders
// the sequence literally in HTML contexts. Only a `#` directly followed by
// alphanumeric characters and a semicolon is guarded, so hex color values in
// style directives (`fill:#ff0000`, never semicolon-terminated) pass through
// unchanged. Applied after the replacer; its own output is never rescanned.
//
//nolint:gochecknoglobals // Compiled once, safe to share.
var mermaidEntityGuard = regexp.MustCompile(`#([0-9A-Za-z]+;)`)

// MermaidText escapes special characters for Mermaid display labels.
func MermaidText(s string) string {
	return mermaidEntityGuard.ReplaceAllString(mermaidTextReplacer.Replace(s), "&#35;$1")
}

// MarkdownCell escapes a string for use inside a GitHub-Flavored Markdown
// table cell. A raw `|` would end the cell (corrupting the column layout),
// and a raw newline would break the row; a raw `\` could combine with a
// following `|` into an unintended escape. Newlines become <br> line breaks,
// which Markdown renders inside the cell.
//
//nolint:gochecknoglobals // Reusable strings.Replacer, safe to share.
var markdownCellReplacer = strings.NewReplacer(
	`\`, `\\`,
	"|", `\|`,
	"\n", "<br>",
	"\r", "",
)

// MarkdownCell escapes special characters for Markdown table cells.
func MarkdownCell(s string) string {
	return markdownCellReplacer.Replace(s)
}

// PlantUMLID sanitizes a string for use as a PlantUML node identifier.
// PlantUML has no quoted-ID escape hatch (unlike DOT), so IDs must be
// allowlisted: SlugifyID first (preserving the established underscore
// slugging), then any rune outside ASCII letters, digits, and underscore is
// dropped — this removes newline/`@`/`:`/`;` vectors like `a\n@enduml`
// that could terminate the diagram or forge directives. Empty results fall
// back to "node", matching MermaidID.
func PlantUMLID(s string) string {
	slug := SlugifyID(s)

	var result strings.Builder

	for _, r := range slug {
		if isMermaidIdentRune(r) {
			result.WriteRune(r)
		}
	}

	if result.Len() == 0 {
		return fallbackID
	}

	return result.String()
}

// plantumlReplacer escapes characters that break PlantUML component notation
// and edge labels: right bracket (closes [label]), newline (breaks the line),
// and backslash/quote for general string safety.
//
//nolint:gochecknoglobals // Reusable strings.Replacer, safe to share.
var plantumlReplacer = strings.NewReplacer(
	"]", "\\]",
	"\n", "\\n",
	`\`, `\\`,
	`"`, `\"`,
)

// PlantUML escapes special characters for PlantUML labels and text.
func PlantUML(s string) string {
	return plantumlReplacer.Replace(s)
}

// Standard ANSI escape sequences for terminal styling. Shared across the
// zero-dep renderer modules that emit ANSI styling (markdown, tree).
// Prefixed ANSI to keep their purpose explicit and avoid collision with
// user-defined names at call sites.
const (
	ANSIReturn = "\033[0m"
	ANSIBold   = "\033[1m"
	ANSIDim    = "\033[2m"
)

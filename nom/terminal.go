package nom

import (
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
	"golang.org/x/term"
)

// defaultTerminalWidth is used when the terminal size cannot be detected.
const defaultTerminalWidth = 80

// GetTerminalWidth returns the width of the terminal associated with writer, or
// a sensible default if the writer is not a terminal or size detection fails.
func GetTerminalWidth(writer io.Writer) int {
	if f, ok := writer.(*os.File); ok {
		if width, _, err := term.GetSize(int(f.Fd())); err == nil && width > 0 {
			return width
		}
	}

	return defaultTerminalWidth
}

// VisibleWidth returns the number of visible columns a string occupies after
// stripping ANSI escape sequences.
func VisibleWidth(s string) int {
	plain := StripANSI(s)

	return runewidth.StringWidth(plain)
}

// VisibleLineCount returns how many physical terminal lines a logical line
// occupies given a terminal width. Empty content counts as zero lines.
func VisibleLineCount(line string, termWidth int) int {
	if line == "" || termWidth <= 0 {
		return 0
	}

	w := VisibleWidth(line)
	if w == 0 {
		return 0
	}

	return (w + termWidth - 1) / termWidth
}

// TruncateVisible truncates s so its visible width is at most maxWidth. If
// truncation occurs, "…" is appended (included in the width budget).
func TruncateVisible(s string, maxWidth int) string {
	if maxWidth <= 0 {
		return ""
	}

	if VisibleWidth(s) <= maxWidth {
		return s
	}

	const ellipsis = "…"

	ellipsisWidth := runewidth.StringWidth(ellipsis)
	if maxWidth < ellipsisWidth {
		// Even an ellipsis won't fit; return as many plain characters as fit.
		return truncateToWidth(s, maxWidth)
	}

	return truncateToWidth(s, maxWidth-ellipsisWidth) + ellipsis
}

// truncateToWidth returns a prefix of s whose visible width is at most width.
func truncateToWidth(s string, width int) string {
	if width <= 0 {
		return ""
	}

	// Preserve leading ANSI sequences so styling is retained.
	prefix := ANSIPrefix(s)
	plain := StripANSI(s)

	var b strings.Builder
	b.WriteString(prefix)

	current := 0

	for _, r := range plain {
		w := runewidth.RuneWidth(r)
		if current+w > width {
			break
		}

		b.WriteRune(r)

		current += w
	}

	return b.String()
}

// StripANSI removes ANSI escape sequences from s.
func StripANSI(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	for len(s) > 0 {
		skip, ok := scanANSI(s)
		if ok {
			s = s[skip:]

			continue
		}

		r, size := decodeRune(s)
		b.WriteRune(r)

		s = s[size:]
	}

	return b.String()
}

// ANSIPrefix returns the leading ANSI escape sequences in s, if any.
func ANSIPrefix(s string) string {
	var prefix strings.Builder

	for len(s) > 0 {
		skip, ok := scanANSI(s)
		if !ok {
			break
		}

		prefix.WriteString(s[:skip])
		s = s[skip:]
	}

	return prefix.String()
}

// scanANSI reports whether s starts with an ANSI escape sequence and returns
// the number of bytes to skip past it.
func scanANSI(s string) (int, bool) {
	if len(s) < 2 || s[0] != '\x1b' {
		return 0, false
	}

	if s[1] == '[' {
		i := 2
		for i < len(s) && isANSIPayload(s[i]) {
			i++
		}

		if i < len(s) && isANSIFinal(s[i]) {
			return i + 1, true
		}

		return 0, false
	}

	// Single-byte CSI-alternatives (e.g. \x1b] OSC, \x1b( SCS, \x1b) SCS).
	if s[1] >= 0x40 && s[1] <= 0x5f {
		return 2, true
	}

	return 0, false
}

// isANSIPayload reports whether b is a valid intermediate/parameter byte.
func isANSIPayload(b byte) bool {
	return (b >= 0x30 && b <= 0x3f) || (b >= 0x20 && b <= 0x2f)
}

// isANSIFinal reports whether b is a valid final byte of an ANSI sequence.
func isANSIFinal(b byte) bool {
	return b >= 0x40 && b <= 0x7e
}

// decodeRune decodes the first UTF-8 rune in s using the standard library.
func decodeRune(s string) (rune, int) {
	return utf8.DecodeRuneInString(s)
}

// PhysicalLineCount returns the total physical terminal lines occupied by a
// multi-line frame, accounting for line wrapping at termWidth.
func PhysicalLineCount(frame string, termWidth int) int {
	if frame == "" {
		return 0
	}

	lines := strings.Split(frame, "\n")
	total := 0

	for _, line := range lines {
		total += max(1, VisibleLineCount(line, termWidth))
	}

	return total
}

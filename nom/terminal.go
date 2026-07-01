package nom

import (
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/x/ansi"
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
	return ansi.StringWidth(s)
}

// VisibleLineCount returns how many physical terminal lines a logical line
// occupies given a terminal width. Uses ansi.Hardwrap for grapheme-aware
// wrapping that correctly handles wide-character straddling (a 2-cell CJK
// character at position width-1 wraps to the next line, leaving a gap —
// the old ceiling-division formula undercounted this). Empty content counts
// as zero lines.
func VisibleLineCount(line string, termWidth int) int {
	if line == "" || termWidth <= 0 {
		return 0
	}

	wrapped := ansi.Hardwrap(line, termWidth, true)

	return strings.Count(wrapped, "\n") + 1
}

// TruncateVisible truncates s so its visible width is at most maxWidth. If
// truncation occurs, "…" is appended (included in the width budget).
func TruncateVisible(s string, maxWidth int) string {
	return ansi.Truncate(s, maxWidth, "…")
}

// StripANSI removes ANSI escape sequences from s.
func StripANSI(s string) string {
	return ansi.Strip(s)
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

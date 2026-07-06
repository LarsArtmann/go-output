package markdown

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

func testMarkdownFooter(t *testing.T) {
	t.Helper()

	t.Run("renders footer after separator", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdownTable()
		m.SetHeaders([]string{"Name", "Count"})
		m.AddRow([]string{"Alice", "10"})
		m.AddRow([]string{"Bob", "20"})
		m.SetFooter([]string{"Total", "30"})

		got, err := m.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		assertContains(t, got, "Total", "should contain footer text")
		assertContains(t, got, "30", "should contain footer value")

		lines := strings.Split(strings.TrimSpace(got), "\n")
		if len(lines) != 6 {
			t.Fatalf("expected 6 lines (header + sep + 2 rows + sep + footer), got %d:\n%s", len(lines), got)
		}

		if !strings.Contains(lines[4], "---") {
			t.Errorf("expected separator before footer, got %q", lines[4])
		}

		if !strings.Contains(lines[5], "Total") {
			t.Errorf("expected footer on last line, got %q", lines[5])
		}
	})

	t.Run("no footer by default", func(t *testing.T) {
		t.Parallel()

		m := newMarkdownTableWithData()

		got, err := m.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		lines := strings.Split(strings.TrimSpace(got), "\n")
		if len(lines) != 4 {
			t.Errorf("expected 4 lines without footer, got %d", len(lines))
		}
	})

	t.Run("footer from Table", func(t *testing.T) {
		t.Parallel()

		data := output.NewTable([]string{"Item", "Qty"})
		data.AddRow([]string{"Apple", "5"})
		data.Footer = []string{"Sum", "5"}

		m := NewMarkdownTableFromTable(data)

		got, err := m.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		assertContains(t, got, "Sum", "should contain footer from Table")
	})

	t.Run("footer inherits column alignment", func(t *testing.T) {
		t.Parallel()

		m := NewMarkdownTable()
		m.SetHeaders([]string{"Left", "Right", "Center"})
		m.SetAlign(0, AlignLeft)
		m.SetAlign(1, AlignRight)
		m.SetAlign(2, AlignCenter)
		m.AddRow([]string{"a", "b", "c"})
		m.SetFooter([]string{"L", "R", "C"})

		got, err := m.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		lines := strings.Split(strings.TrimSpace(got), "\n")
		if len(lines) != 5 {
			t.Fatalf("expected 5 lines, got %d:\n%s", len(lines), got)
		}

		footerSeparator := lines[3]
		if !strings.Contains(footerSeparator, "---:") {
			t.Errorf("footer separator should have right-alignment marker '---:', got %q", footerSeparator)
		}

		if !strings.Contains(footerSeparator, ":---") {
			t.Errorf("footer separator should have center-alignment marker ':---', got %q", footerSeparator)
		}

		footerLine := lines[4]
		if !strings.Contains(footerLine, "| L ") {
			t.Errorf("footer left cell should be left-aligned '| L ', got %q", footerLine)
		}

		if !strings.Contains(footerLine, "R |") {
			t.Errorf("footer right cell should be right-aligned 'R |', got %q", footerLine)
		}
	})
}

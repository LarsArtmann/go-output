package table

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

func TestTableColorModeNever(t *testing.T) {
	t.Parallel()

	tbl := New(WithColorMode(output.ColorModeNever))
	tbl.SetHeaders("Name", "Value")
	tbl.AddRow("Alice", "30")

	got, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if strings.Contains(got, "\x1b[") {
		t.Errorf("ColorModeNever should produce no ANSI codes, got: %q", got)
	}

	testhelpers.AssertContains(t, got, "Alice", "should contain data even without colors")
}

func TestTableColorModeAlways(t *testing.T) {
	t.Parallel()

	tbl := New(WithColorMode(output.ColorModeAlways))
	tbl.SetHeaders("Name", "Value")
	tbl.AddRow("Alice", "30")

	got, err := tbl.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if !strings.Contains(got, "\x1b[") {
		t.Errorf("ColorModeAlways should produce ANSI codes, got: %q", got)
	}
}

func TestTableColorModeDefault(t *testing.T) {
	t.Parallel()

	tbl := New()
	if tbl.colorMode != output.ColorModeAuto {
		t.Errorf("default ColorMode = %v, want %v", tbl.colorMode, output.ColorModeAuto)
	}
}

func TestAsTableRenderer(t *testing.T) {
	t.Parallel()

	tbl := New()
	tr := tbl.AsTableRenderer()

	tr.SetHeaders([]string{"A", "B"})
	tr.AddRow([]string{"1", "2"})

	got, err := tr.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(t, got, "A", "should contain header A")
	testhelpers.AssertContains(t, got, "1", "should contain cell 1")

	var _ output.TableRenderer = tr
}

func TestBuildStyleFunc_AllBranches(t *testing.T) {
	t.Parallel()

	t.Run("always color header footer even odd rows", func(t *testing.T) {
		t.Parallel()

		tbl := New(WithColorMode(output.ColorModeAlways))
		tbl.SetHeaders("A", "B")
		tbl.AddRow("1", "2")
		tbl.AddRow("3", "4")
		tbl.SetFooter("T", "6")

		got, err := tbl.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		if !strings.Contains(got, "\x1b[") {
			t.Error("ColorModeAlways should produce ANSI codes")
		}

		testhelpers.AssertContains(t, got, "1", "should contain row")
		testhelpers.AssertContains(t, got, "3", "should contain row")
		testhelpers.AssertContains(t, got, "T", "should contain footer")
	})

	t.Run("never color header footer even odd rows", func(t *testing.T) {
		t.Parallel()

		tbl := New(WithColorMode(output.ColorModeNever))
		tbl.SetHeaders("A", "B")
		tbl.AddRow("1", "2")
		tbl.AddRow("3", "4")
		tbl.SetFooter("T", "6")

		got, err := tbl.Render()
		if err != nil {
			t.Fatalf("Render() error = %v", err)
		}

		if strings.Contains(got, "\x1b[") {
			t.Error("ColorModeNever should produce no ANSI codes")
		}

		testhelpers.AssertContains(t, got, "1", "should contain row")
		testhelpers.AssertContains(t, got, "3", "should contain row")
		testhelpers.AssertContains(t, got, "T", "should contain footer")
	})
}

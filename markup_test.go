package output

import (
	"strings"
	"testing"
)

func TestWriteRowTag(t *testing.T) {
	t.Parallel()

	t.Run("open tag", func(t *testing.T) {
		t.Parallel()

		var buf strings.Builder

		err := writeRowTag(&buf, "  ", "row", false)
		if err != nil {
			t.Fatalf("writeRowTag() error = %v", err)
		}

		assertContains(t, buf.String(), "<row>", "should contain open tag")
	})

	t.Run("close tag", func(t *testing.T) {
		t.Parallel()

		var buf strings.Builder

		err := writeRowTag(&buf, "  ", "row", true)
		if err != nil {
			t.Fatalf("writeRowTag() error = %v", err)
		}

		assertContains(t, buf.String(), "</row>", "should contain close tag")
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		err := writeRowTag(&errorWriter{}, "  ", "row", false)
		if err == nil {
			t.Fatal("expected error from errorWriter")
		}

		assertContains(t, err.Error(), "open", "error should mention open")
	})

	t.Run("close tag error", func(t *testing.T) {
		t.Parallel()

		err := writeRowTag(&errorWriter{}, "  ", "row", true)
		if err == nil {
			t.Fatal("expected error from errorWriter")
		}

		assertContains(t, err.Error(), "close", "error should mention close")
	})
}

func TestWriteMarkupRow(t *testing.T) {
	t.Parallel()

	t.Run("writes row with cells", func(t *testing.T) {
		t.Parallel()

		var buf strings.Builder

		err := writeMarkupRow(
			&buf, []string{"A", "B"}, "row", "cell", "  ",
			func(s string) string { return s },
		)
		if err != nil {
			t.Fatalf("writeMarkupRow() error = %v", err)
		}

		result := buf.String()
		assertContains(t, result, "<row>", "should contain row open")
		assertContains(t, result, "</row>", "should contain row close")
		assertContains(t, result, "<cell>A</cell>", "should contain cell A")
		assertContains(t, result, "<cell>B</cell>", "should contain cell B")
	})

	t.Run("error on open row tag", func(t *testing.T) {
		t.Parallel()

		err := writeMarkupRow(
			&errorWriter{}, []string{"A"}, "row", "cell", "  ",
			func(s string) string { return s },
		)
		if err == nil {
			t.Fatal("expected error from errorWriter")
		}

		assertContains(t, err.Error(), "open row", "error should mention open row")
	})
}

func TestWriteMarkupColumns(t *testing.T) {
	t.Parallel()

	t.Run("writes columns", func(t *testing.T) {
		t.Parallel()

		var buf strings.Builder

		err := writeMarkupColumns(
			&buf, []string{"Name", "Val"}, "  ",
			func(s string) string { return s },
		)
		if err != nil {
			t.Fatalf("writeMarkupColumns() error = %v", err)
		}

		result := buf.String()
		assertContains(t, result, "<column>Name</column>", "should contain column")
		assertContains(t, result, "<column>Val</column>", "should contain column")
	})

	t.Run("error on column open", func(t *testing.T) {
		t.Parallel()

		err := writeMarkupColumns(
			&errorWriter{}, []string{"A"}, "  ",
			func(s string) string { return s },
		)
		if err == nil {
			t.Fatal("expected error from errorWriter")
		}

		assertContains(t, err.Error(), "column tag open", "error should mention column tag open")
	})
}

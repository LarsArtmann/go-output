package markup

import (
	"testing"
)

func TestXMLWriterWriteHeaderError(t *testing.T) {
	t.Parallel()

	x := NewXMLWriter(&errorWriter{})

	err := x.WriteHeader([]string{"Name"})
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}

	assertContains(t, err.Error(), "xml header", "error should mention xml header")
}

func TestXMLWriterWriteHeaderPartialErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		remaining   int
		wantErrPart string
	}{
		{"table open", 1, "table open"},
		{"headers open", 2, "headers open"},
		{"columns", 3, "columns"},
		{"headers close", 6, "headers close"},
		{"rows open", 7, "rows open"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			x := NewXMLWriter(&writeNThenFailWriter{Remaining: tt.remaining})

			err := x.WriteHeader([]string{"Name"})
			if err == nil {
				t.Fatal("expected error")
			}

			assertContains(t, err.Error(), tt.wantErrPart, "error should mention "+tt.wantErrPart)
		})
	}
}

func TestXMLWriterWriteFooterTableCloseError(t *testing.T) {
	t.Parallel()

	x := NewXMLWriter(&writeNThenFailWriter{Remaining: 1})

	err := x.WriteFooter()
	if err == nil {
		t.Fatal("expected error on table close")
	}

	assertContains(t, err.Error(), "table close", "error should mention table close")
}

func TestXMLWriterWriteRowError(t *testing.T) {
	t.Parallel()

	x := NewXMLWriter(&errorWriter{})

	err := x.WriteRow([]string{"test"})
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}

	assertContains(t, err.Error(), "row", "error should mention row")
}

func TestXMLWriterWriteRowsError(t *testing.T) {
	t.Parallel()

	x := NewXMLWriter(&errorWriter{})

	err := x.WriteRows([][]string{{"a", "b"}})
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}

	assertContains(t, err.Error(), "row", "error should mention row")
}

func TestXMLWriterWriteFooterError(t *testing.T) {
	t.Parallel()

	x := NewXMLWriter(&errorWriter{})

	err := x.WriteFooter()
	if err == nil {
		t.Fatal("expected error from errorWriter")
	}

	assertContains(t, err.Error(), "rows close", "error should mention rows close")
}

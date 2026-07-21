package delimited

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

func TestCSVWriter(t *testing.T) {
	t.Parallel()
	t.Run("write header and rows", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer

		w := NewCSVWriter(&buf)

		err := w.WriteHeader([]string{"Name", "Age"})
		if err != nil {
			t.Errorf("WriteHeader() error = %v", err)
		}

		err = w.WriteRow([]string{"Alice", "30"})
		if err != nil {
			t.Errorf("WriteRow() error = %v", err)
		}

		err = w.WriteRow([]string{"Bob", "25"})
		if err != nil {
			t.Errorf("WriteRow() error = %v", err)
		}

		w.Flush()

		output := buf.String()
		if output == "" {
			t.Error("CSVWriter produced empty output")
		}
	})

	t.Run("write multiple rows", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer

		w := NewCSVWriter(&buf)

		rows := [][]string{
			{"Name", "Age"},
			{"Alice", "30"},
			{"Bob", "25"},
		}

		err := w.WriteRows(rows)
		if err != nil {
			t.Errorf("WriteRows() error = %v", err)
		}

		w.Flush()

		output := buf.String()
		if output == "" {
			t.Error("CSVWriter produced empty output")
		}
	})

	t.Run("flush and error", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer

		w := NewCSVWriter(&buf)

		w.Flush()

		err := w.Error()
		if err != nil {
			t.Errorf("Error() should return nil after flush, got %v", err)
		}
	})
}

func TestNewCSVWriter(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	w := NewCSVWriter(&buf)

	if w.writer == nil {
		t.Error("NewCSVWriter() did not initialize writer")
	}
}

func TestCSVWriterErrorPaths(t *testing.T) {
	t.Parallel()

	t.Run("WriteRows error", func(t *testing.T) {
		t.Parallel()

		w := NewCSVWriter(&errorWriter{})

		err := w.WriteRows([][]string{{"Name"}, {"Alice"}})
		if err == nil {
			t.Error("WriteRows() should return error with failing writer")
		}
	})

	t.Run("Error method returns error", func(t *testing.T) {
		t.Parallel()

		w := NewCSVWriter(&errorWriter{})
		_ = w.WriteRow([]string{"test"})
		w.Flush()

		err := w.Error()
		if err == nil {
			t.Error("Error() should return error after failed write")
		}
	})

	t.Run("WriteHeader error on flush", func(t *testing.T) {
		t.Parallel()

		w := NewCSVWriter(&errorWriter{})
		_ = w.WriteHeader([]string{"Name"})
		w.Flush()

		err := w.Error()
		if err == nil {
			t.Error("Error() should return error after WriteHeader with failing writer")
		}
	})
}

func TestMarshalCSVFromTable(t *testing.T) {
	t.Parallel()

	t.Run("with headers and rows", func(t *testing.T) {
		t.Parallel()

		data := output.NewTable([]string{"Name", "Age"})
		data.AddRow([]string{"Alice", "30"})
		data.AddRow([]string{"Bob", "25"})

		b, err := MarshalCSVFromTable(data)
		if err != nil {
			t.Fatalf("MarshalCSVFromTable() error = %v", err)
		}

		result := string(b)
		if !strings.Contains(result, "Name") {
			t.Error("CSV should contain header 'Name'")
		}

		if !strings.Contains(result, "Alice") {
			t.Error("CSV should contain row 'Alice'")
		}
	})

	t.Run("nil data returns nil", func(t *testing.T) {
		t.Parallel()

		b, err := MarshalCSVFromTable(nil)
		if err != nil {
			t.Fatalf("MarshalCSVFromTable(nil) error = %v", err)
		}

		if b != nil {
			t.Errorf("MarshalCSVFromTable(nil) = %q, want nil", b)
		}
	})

	t.Run("empty data", func(t *testing.T) {
		t.Parallel()

		data := output.NewTable(nil)

		b, err := MarshalCSVFromTable(data)
		if err != nil {
			t.Fatalf("MarshalCSVFromTable() error = %v", err)
		}

		if len(b) != 0 {
			t.Errorf("MarshalCSVFromTable(empty) = %q, want empty", b)
		}
	})

	t.Run("headers only no rows", func(t *testing.T) {
		t.Parallel()

		data := output.NewTable([]string{"Name", "Age"})

		b, err := MarshalCSVFromTable(data)
		if err != nil {
			t.Fatalf("MarshalCSVFromTable() error = %v", err)
		}

		result := string(b)
		if !strings.Contains(result, "Name") {
			t.Error("CSV should contain header even with no rows")
		}
	})

	t.Run("with footer row", func(t *testing.T) {
		t.Parallel()

		data := output.NewTable([]string{"Name", "Count"})
		data.AddRow([]string{"Alice", "10"})
		data.AddRow([]string{"Bob", "20"})
		data.Footer = []string{"Total", "30"}

		b, err := MarshalCSVFromTable(data)
		if err != nil {
			t.Fatalf("MarshalCSVFromTable() error = %v", err)
		}

		result := string(b)
		if !strings.Contains(result, "Total") {
			t.Error("CSV should contain footer row")
		}

		assertLineCount(t, "CSV with footer", result, 4)
		assertLastLineContains(t, "CSV with footer", result, "Total")
	})
}

func TestCSVRenderTable(t *testing.T) {
	t.Parallel()

	testRenderTable(t, output.FormatCSV, "CSV")
}

func TestCSVWriter_WriteFooter(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	w := NewCSVWriter(&buf)
	_ = w.WriteHeader([]string{"Name", "Count"})
	_ = w.WriteRow([]string{"Alice", "10"})
	_ = w.WriteFooter([]string{"Total", "10"})
	w.Flush()

	result := buf.String()

	assertLineCount(t, "WriteFooter", result, 3)
	assertLastLineContains(t, "WriteFooter", result, "Total")
}

func TestMarshalFromTable_FlushError(t *testing.T) {
	t.Parallel()

	data := output.NewTableWithRow([]string{"Name"}, "Alice")

	_, err := marshalFromTable(data, "csv", func(w io.Writer) tableDataWriter {
		return NewCSVWriter(&errorWriter{})
	})
	if err == nil {
		t.Fatal("expected flush error from errorWriter")
	}

	if !strings.Contains(err.Error(), "flush csv writer") {
		t.Errorf("error should mention flush, got: %v", err)
	}
}

package output

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

type errorWriter struct{}

func (e *errorWriter) Write(_ []byte) (n int, err error) {
	return 0, errors.New("write error")
}

var _ io.Writer = (*errorWriter)(nil)

func TestCSVWriter(t *testing.T) {
	t.Parallel()
	t.Run("write header and rows", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		w := NewCSVWriter(&buf)

		if err := w.WriteHeader([]string{"Name", "Age"}); err != nil {
			t.Errorf("WriteHeader() error = %v", err)
		}

		if err := w.WriteRow([]string{"Alice", "30"}); err != nil {
			t.Errorf("WriteRow() error = %v", err)
		}

		if err := w.WriteRow([]string{"Bob", "25"}); err != nil {
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

		if err := w.WriteRows(rows); err != nil {
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

		if err := w.Error(); err != nil {
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
		_ = w.WriteRow([]string{"test"}) // Buffer the write
		w.Flush()                        // Error occurs on flush

		err := w.Error()
		if err == nil {
			t.Error("Error() should return error after failed write")
		}
	})
}

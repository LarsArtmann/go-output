package delimited

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/testhelpers"
)

var errMarshalFailed = errors.New("marshal failed")

func TestRenderDelimitedTable(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"Name", "Age"})
	data.AddRow([]string{"Alice", "30"})

	var buf bytes.Buffer

	err := renderDelimitedTable(&buf, data, MarshalCSVFromTable, "csv")
	if err != nil {
		t.Fatalf("renderDelimitedTable csv: %v", err)
	}

	if !strings.Contains(buf.String(), "Alice") {
		t.Errorf("expected Alice in output, got %q", buf.String())
	}
}

func TestRenderDelimitedTable_NilData(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	err := renderDelimitedTable(&buf, nil, MarshalCSVFromTable, "csv")
	if err != nil {
		t.Fatalf("renderDelimitedTable nil: %v", err)
	}

	if buf.String() != "" {
		t.Errorf("expected empty output for nil data, got %q", buf.String())
	}
}

func TestRenderDelimitedTable_WriterError(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"Name"})
	data.AddRow([]string{"Alice"})

	err := renderDelimitedTable(&testhelpers.ErrorWriter{}, data, MarshalCSVFromTable, "csv")
	if err == nil {
		t.Fatal("expected error from failWriter")
	}
}

func TestRenderDelimitedTable_MarshalError(t *testing.T) {
	t.Parallel()

	failMarshal := func(_ *output.Table) ([]byte, error) {
		return nil, errMarshalFailed
	}

	data := output.NewTable([]string{"Name"})

	err := renderDelimitedTable(&bytes.Buffer{}, data, failMarshal, "csv")
	if err == nil {
		t.Fatal("expected error from failMarshal")
	}
}

func TestMarshalDelimitedFromTable_NoHeaders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		marshalFunc func(*output.Table) ([]byte, error)
	}{
		{"CSV", MarshalCSVFromTable},
		{"TSV", MarshalTSVFromTable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data := output.NewTable(nil)
			data.AddRow([]string{"Alice", "30"})

			b, err := tt.marshalFunc(data)
			if err != nil {
				t.Fatalf("marshal %s no headers: %v", tt.name, err)
			}

			if !strings.Contains(string(b), "Alice") {
				t.Errorf("expected Alice in %s output, got %q", tt.name, string(b))
			}
		})
	}
}

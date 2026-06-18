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

func TestRenderDelimitedTableData(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{"Name", "Age"})
	data.AddRow([]string{"Alice", "30"})

	var buf bytes.Buffer

	err := renderDelimitedTableData(&buf, data, MarshalCSVFromTableData, "csv")
	if err != nil {
		t.Fatalf("renderDelimitedTableData csv: %v", err)
	}

	if !strings.Contains(buf.String(), "Alice") {
		t.Errorf("expected Alice in output, got %q", buf.String())
	}
}

func TestRenderDelimitedTableData_NilData(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	err := renderDelimitedTableData(&buf, nil, MarshalCSVFromTableData, "csv")
	if err != nil {
		t.Fatalf("renderDelimitedTableData nil: %v", err)
	}

	if buf.String() != "" {
		t.Errorf("expected empty output for nil data, got %q", buf.String())
	}
}

func TestRenderDelimitedTableData_WriterError(t *testing.T) {
	t.Parallel()

	data := output.NewTableData([]string{"Name"})
	data.AddRow([]string{"Alice"})

	err := renderDelimitedTableData(&testhelpers.ErrorWriter{}, data, MarshalCSVFromTableData, "csv")
	if err == nil {
		t.Fatal("expected error from failWriter")
	}
}

func TestRenderDelimitedTableData_MarshalError(t *testing.T) {
	t.Parallel()

	failMarshal := func(_ *output.TableData) ([]byte, error) {
		return nil, errMarshalFailed
	}

	data := output.NewTableData([]string{"Name"})

	err := renderDelimitedTableData(&bytes.Buffer{}, data, failMarshal, "csv")
	if err == nil {
		t.Fatal("expected error from failMarshal")
	}
}

func TestMarshalDelimitedFromTableData_NoHeaders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		marshalFunc func(*output.TableData) ([]byte, error)
	}{
		{"CSV", MarshalCSVFromTableData},
		{"TSV", MarshalTSVFromTableData},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			data := output.NewTableData(nil)
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

func TestMarshalTSV_Slice(t *testing.T) {
	t.Parallel()

	b, err := MarshalTSV([]string{"Alice", "30"})
	if err != nil {
		t.Fatalf("MarshalTSV slice: %v", err)
	}

	if !strings.Contains(string(b), "Alice") {
		t.Errorf("expected Alice in output, got %q", string(b))
	}
}

func TestMarshalTSV_UnsupportedType(t *testing.T) {
	t.Parallel()

	_, err := MarshalTSV(42)
	if err == nil {
		t.Fatal("expected error for unsupported type")
	}

	if !errors.Is(err, ErrUnsupportedType) {
		t.Errorf("expected ErrUnsupportedType, got %v", err)
	}
}

package delimited

import (
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
)

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

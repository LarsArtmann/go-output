package serialization

import (
	"bytes"
	"testing"

	"github.com/charmbracelet/x/exp/golden"

	"github.com/larsartmann/go-output"
)

// sampleCQRSTable returns a stable dataset used by every CQRS golden test.
// Using one shared fixture means the formats can be cross-compared and the
// golden files stay human-reviewable.
func sampleCQRSTable() *output.Table {
	data := output.NewTable([]string{"Name", "Status", "Duration"})
	data.AddRow([]string{"Build", "completed", "1.2s"})
	data.AddRow([]string{"Test", "running", "0.5s"})

	return data
}

// TestGolden_CQRS_JSON locks in the exact byte output of WriteJSON.
// Unlike the registry path (json.MarshalIndent, no trailing newline),
// the streaming path uses json.NewEncoder which appends a trailing newline.
// This test makes that behavior difference explicit and regression-proof.
func TestGolden_CQRS_JSON(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := WriteJSON(&buf, sampleCQRSTable()); err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, buf.Bytes())
}

// TestGolden_CQRS_YAML locks in the exact byte output of WriteYAML.
func TestGolden_CQRS_YAML(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := WriteYAML(&buf, sampleCQRSTable()); err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, buf.Bytes())
}

// TestGolden_CQRS_TOML locks in the exact byte output of WriteTOML.
func TestGolden_CQRS_TOML(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := WriteTOML(&buf, sampleCQRSTable()); err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, buf.Bytes())
}

// TestGolden_CQRS_JSONL locks in the exact byte output of WriteJSONL.
func TestGolden_CQRS_JSONL(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := WriteJSONL(&buf, sampleCQRSTable()); err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, buf.Bytes())
}

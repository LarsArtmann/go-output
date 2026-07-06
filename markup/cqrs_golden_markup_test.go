package markup

import (
	"bytes"
	"testing"

	"github.com/charmbracelet/x/exp/golden"
)

func TestGolden_CQRS_HTML(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := WriteHTML(&buf, sampleCQRSTable()); err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, buf.Bytes())
}

func TestGolden_CQRS_AsciiDoc(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := WriteAsciiDoc(&buf, sampleCQRSTable()); err != nil {
		t.Fatal(err)
	}

	golden.RequireEqual(t, buf.Bytes())
}

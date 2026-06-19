package markdown

import (
	"slices"
	"testing"

	"github.com/larsartmann/go-output"
)

// TestRegistersMarkdownFormat verifies the markdown module self-registers its
// TableData renderer via init(), mirroring how every format sub-module works.
func TestRegistersMarkdownFormat(t *testing.T) {
	t.Parallel()

	formats := output.RegisteredTableDataFormats()

	if !slices.Contains(formats, output.FormatMarkdown) {
		t.Errorf("expected markdown module to register %q. Registered: %v", output.FormatMarkdown, formats)
	}
}

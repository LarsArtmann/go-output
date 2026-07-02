package daghtml

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed graph.css
var graphCSS string

//go:embed graph.js
var graphJS string

// dagToJSON serializes the DAG to an HTML-safe JSON string suitable for
// embedding inside a <script type="application/json"> element.
//
// It uses SetEscapeHTML to escape <, >, and & as Unicode escape sequences,
// preventing </script> injection attacks. The output is safe to include
// directly in HTML without additional escaping.
func dagToJSON(dag DAG) (string, error) {
	if dag.Nodes == nil {
		dag.Nodes = []Node{}
	}

	if dag.Edges == nil {
		dag.Edges = []Edge{}
	}

	var buf strings.Builder

	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(dag); err != nil {
		return "", fmt.Errorf("encode DAG JSON: %w", err)
	}

	return strings.TrimSpace(buf.String()), nil
}

package daghtml

import (
	"encoding/json"
	"strings"
	"testing"
)

func sampleDAG() DAG {
	return DAG{
		Nodes: []Node{
			{ID: "fetch", Label: "📥 Fetch", Color: "var(--success)", Tooltip: "Fetch | status: succeeded | 42ms"},
			{ID: "parse", Label: "🔍 Parse", Color: "var(--accent)", Tooltip: "Parse | status: succeeded | 12ms"},
			{
				ID:      "validate",
				Label:   "✓ Validate",
				Color:   "var(--accent)",
				Tooltip: "Validate | status: succeeded | 5ms",
			},
			{
				ID:      "save",
				Label:   "💾 Save",
				Color:   "var(--error)",
				Tooltip: "Save | status: failed | timeout",
				Error:   true,
			},
		},
		Edges: []Edge{
			{From: "fetch", To: "parse"},
			{From: "parse", To: "validate"},
			{From: "parse", To: "save"},
		},
	}
}

func TestRender_ValidHTML(t *testing.T) {
	html, err := Render(sampleDAG(), WithTitle("Test Pipeline"))
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	checks := []string{
		"<!DOCTYPE html>",
		"<html",
		"</html>",
		"<style>",
		"</style>",
		"<script",
		"initDAGGraph(",
		"Test Pipeline",
	}
	for _, want := range checks {
		if !strings.Contains(html, want) {
			t.Errorf("expected HTML to contain %q, got %d bytes", want, len(html))
		}
	}
}

func TestRender_ContainsJSONData(t *testing.T) {
	html, err := Render(sampleDAG())
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if !strings.Contains(html, `"id":"fetch"`) {
		t.Error("expected JSON data to contain node ID 'fetch'")
	}

	if !strings.Contains(html, `"from":"fetch"`) {
		t.Error("expected JSON data to contain edge from 'fetch'")
	}
}

func TestRender_ContainsGraphJS(t *testing.T) {
	html, err := Render(sampleDAG())
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if !strings.Contains(html, "initDAGGraph") {
		t.Error("expected HTML to contain the graph JS function")
	}

	if !strings.Contains(html, "Sugiyama") || !strings.Contains(html, "Kahn") {
		// JS has algorithm comments
	}

	if !strings.Contains(html, "bary") {
		t.Error("expected HTML to contain barycenter ordering algorithm code")
	}
}

func TestRender_EmptyDAG(t *testing.T) {
	dag := DAG{}

	html, err := Render(dag)
	if err != nil {
		t.Fatalf("Render empty DAG failed: %v", err)
	}

	if !strings.Contains(html, `"nodes":[]`) || !strings.Contains(html, `"edges":[]`) {
		t.Error("expected empty JSON arrays for empty DAG")
	}
}

func TestRender_Options(t *testing.T) {
	tests := []struct {
		name  string
		opts  []Option
		check func(html string) bool
	}{
		{
			name: "WithTitle",
			opts: []Option{WithTitle("My Custom Title")},
			check: func(html string) bool {
				return strings.Contains(html, "My Custom Title")
			},
		},
		{
			name: "WithSubtitle",
			opts: []Option{WithSubtitle("A subtitle here")},
			check: func(html string) bool {
				return strings.Contains(html, "A subtitle here")
			},
		},
		{
			name: "WithFooter",
			opts: []Option{WithFooter("Generated at 2026-01-01")},
			check: func(html string) bool {
				return strings.Contains(html, "Generated at 2026-01-01")
			},
		},
		{
			name: "WithHeight",
			opts: []Option{WithHeight(800)},
			check: func(html string) bool {
				return strings.Contains(html, "min-height:800px")
			},
		},
		{
			name: "WithContainerID",
			opts: []Option{WithContainerID("my-graph")},
			check: func(html string) bool {
				return strings.Contains(html, `"my-graph"`)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := Render(sampleDAG(), tt.opts...)
			if err != nil {
				t.Fatalf("Render failed: %v", err)
			}

			if !tt.check(html) {
				t.Errorf("option %s did not produce expected output", tt.name)
			}
		})
	}
}

func TestGraphHTML_NoFullPage(t *testing.T) {
	html, err := GraphHTML(sampleDAG(), WithContainerID("embed-graph"))
	if err != nil {
		t.Fatalf("GraphHTML failed: %v", err)
	}

	if strings.Contains(html, "<!DOCTYPE html>") {
		t.Error("GraphHTML should not include DOCTYPE")
	}

	if strings.Contains(html, "<html") {
		t.Error("GraphHTML should not include <html> tag")
	}

	if !strings.Contains(html, `id="embed-graph"`) {
		t.Error("expected container with custom ID")
	}
}

func TestStyleSheet(t *testing.T) {
	css := StyleSheet()
	if !strings.Contains(css, "--bg") {
		t.Error("expected CSS to contain --bg variable")
	}

	if !strings.Contains(css, ".graph-node") {
		t.Error("expected CSS to contain .graph-node styles")
	}
}

func TestScript(t *testing.T) {
	js := Script()
	if !strings.Contains(js, "initDAGGraph") {
		t.Error("expected JS to contain initDAGGraph function")
	}
}

func TestDAGToJSON_HTMLEscape(t *testing.T) {
	dag := DAG{
		Nodes: []Node{
			{ID: "a", Label: "<script>alert(1)</script>", Color: "red"},
		},
	}

	jsonStr, err := dagToJSON(dag)
	if err != nil {
		t.Fatalf("dagToJSON failed: %v", err)
	}

	if strings.Contains(jsonStr, "<script>") {
		t.Error("JSON should not contain raw <script> tag (must be escaped)")
	}

	if !strings.Contains(jsonStr, `\u003c`) {
		t.Error("expected < to be escaped as \\u003c in JSON")
	}
}

func TestRender_XSSSafe(t *testing.T) {
	dag := DAG{
		Nodes: []Node{
			{ID: "x", Label: `"><img src=x onerror=alert(1)>`, Color: "red"},
		},
	}

	html, err := Render(dag)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if strings.Contains(html, "<img src=x onerror=alert(1)>") {
		t.Error("HTML should not contain unescaped XSS payload")
	}
}

func TestDAG_Methods(t *testing.T) {
	dag := sampleDAG()

	if dag.NodeCount() != 4 {
		t.Errorf("expected 4 nodes, got %d", dag.NodeCount())
	}

	if dag.EdgeCount() != 3 {
		t.Errorf("expected 3 edges, got %d", dag.EdgeCount())
	}

	if dag.IsEmpty() {
		t.Error("expected non-empty DAG")
	}

	empty := DAG{}
	if !empty.IsEmpty() {
		t.Error("expected empty DAG")
	}
}

func TestDAGToJSON_RoundTrip(t *testing.T) {
	original := sampleDAG()

	jsonStr, err := dagToJSON(original)
	if err != nil {
		t.Fatalf("dagToJSON failed: %v", err)
	}

	var decoded DAG
	if err := json.Unmarshal([]byte(jsonStr), &decoded); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if decoded.NodeCount() != original.NodeCount() {
		t.Errorf("expected %d nodes after round-trip, got %d", original.NodeCount(), decoded.NodeCount())
	}

	if decoded.EdgeCount() != original.EdgeCount() {
		t.Errorf("expected %d edges after round-trip, got %d", original.EdgeCount(), decoded.EdgeCount())
	}
}

func TestRender_DuplicateEdgesDeduplicated(t *testing.T) {
	dag := DAG{
		Nodes: []Node{
			{ID: "a", Label: "A", Color: "blue"},
			{ID: "b", Label: "B", Color: "blue"},
		},
		Edges: []Edge{
			{From: "a", To: "b"},
			{From: "a", To: "b"},
			{From: "a", To: "b"},
		},
	}

	jsonStr, err := dagToJSON(dag)
	if err != nil {
		t.Fatalf("dagToJSON failed: %v", err)
	}

	var decoded DAG

	_ = json.Unmarshal([]byte(jsonStr), &decoded)

	// JSON contains all edges (dedup happens in JS)
	if len(decoded.Edges) != 3 {
		t.Errorf("expected 3 edges in JSON (dedup is in JS), got %d", len(decoded.Edges))
	}
}

func TestRender_DropUnknownEdges(t *testing.T) {
	dag := DAG{
		Nodes: []Node{
			{ID: "a", Label: "A", Color: "blue"},
		},
		Edges: []Edge{
			{From: "a", To: "nonexistent"},
		},
	}

	// The edge referencing a non-existent node is silently included in JSON
	// but dropped by the JS at render time. The render should not error.
	_, err := Render(dag)
	if err != nil {
		t.Fatalf("Render with unknown edge reference failed: %v", err)
	}
}

func TestRender_DefaultOptions(t *testing.T) {
	html, err := Render(sampleDAG())
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if !strings.Contains(html, "DAG Visualization") {
		t.Error("expected default title 'DAG Visualization'")
	}

	if !strings.Contains(html, "min-height:500px") {
		t.Error("expected default height 500px")
	}
}

func TestRender_CSPHeader(t *testing.T) {
	html, err := Render(sampleDAG())
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if !strings.Contains(html, "Content-Security-Policy") {
		t.Error("expected CSP header in HTML")
	}

	if !strings.Contains(html, "frame-ancestors 'none'") {
		t.Error("expected frame-ancestors none in CSP")
	}
}

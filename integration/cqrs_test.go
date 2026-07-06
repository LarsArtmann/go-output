package integration

import (
	"bytes"
	"strings"
	"testing"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/delimited"
	"github.com/larsartmann/go-output/graph"
	"github.com/larsartmann/go-output/serialization"
	"github.com/larsartmann/go-output/testhelpers"
)

func TestCQRS_Pipeline_TableToGraphToDOT(t *testing.T) {
	t.Parallel()

	tbl := output.NewTableBuilder().
		SetHeaders("Step").
		AddRow("Fetch").
		AddRow("Build").
		AddRow("Deploy").
		Build()

	g := output.TableToGraph(tbl)

	dot, err := graph.RenderDOT(g)
	if err != nil {
		t.Fatalf("RenderDOT failed: %v", err)
	}

	testhelpers.AssertContains(t, dot, "Fetch", "DOT should contain Fetch node")
	testhelpers.AssertContains(t, dot, "Build", "DOT should contain Build node")
	testhelpers.AssertContains(t, dot, "Deploy", "DOT should contain Deploy node")
	testhelpers.AssertContains(t, dot, "->", "DOT should contain edges")
}

func TestCQRS_Pipeline_GraphToTreeToASCII(t *testing.T) {
	t.Parallel()

	b := output.NewGraphBuilder()
	b.AddNode(*output.NewGraphNode("root", "Root"))
	b.AddNode(*output.NewGraphNode("child", "Child"))
	b.AddEdge(*output.NewGraphEdge("root", "child"))

	root := output.GraphToTree(b.Build())
	if root == nil {
		t.Fatal("GraphToTree returned nil")
	}

	if root.Label.Get() != "Root" {
		t.Fatalf("expected root label 'Root', got %q", root.Label.Get())
	}
}

func TestCQRS_Pipeline_TableBuilderToCSV(t *testing.T) {
	t.Parallel()

	tbl := output.NewTableBuilder().
		SetHeaders("Name", "Status").
		AddRow("Compile", "done").
		AddRow("Test", "done").
		Build()

	csv, err := delimited.RenderCSV(tbl)
	if err != nil {
		t.Fatalf("RenderCSV failed: %v", err)
	}

	testhelpers.AssertContains(t, csv, "Compile", "CSV should contain Compile")
	testhelpers.AssertContains(t, csv, "Test", "CSV should contain Test")
	testhelpers.AssertContains(t, csv, "done", "CSV should contain done")
}

func TestCQRS_Pipeline_TableBuilderToJSON(t *testing.T) {
	t.Parallel()

	tbl := output.NewTableBuilder().
		SetHeaders("Name", "Status").
		AddRow("Compile", "done").
		Build()

	jsonOut, err := serialization.RenderJSON(tbl)
	if err != nil {
		t.Fatalf("RenderJSON failed: %v", err)
	}

	testhelpers.AssertContains(t, jsonOut, "Compile", "JSON should contain Compile")
	testhelpers.AssertContains(t, jsonOut, "Name", "JSON should contain Name header")
}

func TestCQRS_StreamVsRegistry_JSON(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"Name", "Status"})
	data.AddRow([]string{"Compile", "done"})
	data.AddRow([]string{"Test", "done"})

	var cqrsBuf bytes.Buffer
	if err := serialization.WriteJSON(&cqrsBuf, data); err != nil {
		t.Fatalf("CQRS WriteJSON failed: %v", err)
	}

	var registryBuf bytes.Buffer

	err := output.RenderTable(data, output.FormatJSON, output.RenderOptions{Writer: &registryBuf})
	if err != nil {
		t.Fatalf("Registry dispatch failed: %v", err)
	}

	cqrsContent := strings.TrimSpace(cqrsBuf.String())
	registryContent := strings.TrimSpace(registryBuf.String())

	testhelpers.AssertContains(t, cqrsContent, "Compile", "CQRS JSON should contain Compile")
	testhelpers.AssertContains(t, registryContent, "Compile", "Registry JSON should contain Compile")

	testhelpers.AssertContains(t, cqrsContent, "done", "CQRS JSON should contain done")
	testhelpers.AssertContains(t, registryContent, "done", "Registry JSON should contain done")
}

func TestCQRS_StreamVsRegistry_CSV(t *testing.T) {
	t.Parallel()

	data := output.NewTable([]string{"A", "B"})
	data.AddRow([]string{"1", "2"})
	data.AddRow([]string{"3", "4"})

	var cqrsBuf bytes.Buffer
	if err := delimited.WriteCSV(&cqrsBuf, data); err != nil {
		t.Fatalf("CQRS WriteCSV failed: %v", err)
	}

	var registryBuf bytes.Buffer

	err := output.RenderTable(data, output.FormatCSV, output.RenderOptions{Writer: &registryBuf})
	if err != nil {
		t.Fatalf("Registry dispatch failed: %v", err)
	}

	cqrsLines := strings.Split(strings.TrimSpace(cqrsBuf.String()), "\n")
	registryLines := strings.Split(strings.TrimSpace(registryBuf.String()), "\n")

	if len(cqrsLines) != len(registryLines) {
		t.Errorf("line count mismatch: CQRS=%d, registry=%d", len(cqrsLines), len(registryLines))
	}
}

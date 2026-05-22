// Package integration provides end-to-end integration tests for go-output.
package integration

import (
	"testing"

	"github.com/larsartmann/go-output"
)

func TestD2FromTableDataIntegration(t *testing.T) {
	t.Parallel()

	projects := SampleProjects()
	result := renderD2FromTableData(projects)
	testhelpers.AssertContains(t, result, "row0", "D2 from table data should contain row nodes")
	testhelpers.AssertContains(t, result, "->", "D2 from table data should contain edges")
	testhelpers.AssertContains(t, result, "Alpha", "D2 from table data should contain project names")
}

func TestD2FromTreeIntegration(t *testing.T) {
	t.Parallel()

	projects := SampleProjects()
	result := renderD2FromTree(projects)
	testhelpers.AssertContains(t, result, "Projects", "D2 from tree should contain root label")
	testhelpers.AssertContains(t, result, "->", "D2 from tree should contain edges")
	testhelpers.AssertContains(t, result, "Alpha", "D2 from tree should contain child labels")
}

func TestD2ConstraintsIntegration(t *testing.T) {
	t.Parallel()

	d2 := output.NewD2Diagram()
	d2.AddTable("users", []output.D2Column{
		{Name: "id", Type: "int", Constraint: output.D2ConstraintPrimary},
		{Name: "email", Type: "string", Constraint: output.D2ConstraintUnique},
		{Name: "org_id", Type: "int", Constraint: output.D2ConstraintForeign},
	})

	result, err := d2.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(t, result, "constraint: primary_key", "should contain primary key")
	testhelpers.AssertContains(t, result, "constraint: unique", "should contain unique constraint")
	testhelpers.AssertContains(t, result, "constraint: foreign_key", "should contain foreign key")
}

func TestD2ClassesIntegration(t *testing.T) {
	t.Parallel()

	d2 := output.NewD2Diagram()
	d2.AddClass("server", output.D2NodeStyle{Fill: "blue", Stroke: "black"})
	d2.AddNode(output.D2Node{
		ID:    output.NewBrandedID[output.D2NodeIDBrand]("api"),
		Label: output.NewBrandedID[output.D2NodeLabelBrand]("API"),
		Class: "server",
	})

	result, err := d2.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(t, result, "classes:", "should contain classes block")
	testhelpers.AssertContains(t, result, "class: server", "should contain class reference")
}

func TestD2ArrowTypesIntegration(t *testing.T) {
	t.Parallel()

	d2 := output.NewD2Diagram()
	d2.AddEdge(output.D2Edge{
		From:        output.NewBrandedID[output.D2NodeIDBrand]("a"),
		To:          output.NewBrandedID[output.D2NodeIDBrand]("b"),
		TargetArrow: output.D2ArrowDiamond,
		SourceArrow: output.D2ArrowCFMany,
	})

	result, err := d2.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(
		t,
		result,
		"source-arrowhead.shape: cf-many",
		"should contain cf-many arrow",
	)
	testhelpers.AssertContains(
		t,
		result,
		"target-arrowhead.shape: diamond",
		"should contain diamond arrow",
	)
}

func TestD2GridAndNearIntegration(t *testing.T) {
	t.Parallel()

	d2 := output.NewD2Diagram()
	d2.AddNode(output.D2Node{
		ID:          output.NewBrandedID[output.D2NodeIDBrand]("grid"),
		Label:       output.NewBrandedID[output.D2NodeLabelBrand]("Grid"),
		GridRows:    2,
		GridColumns: 3,
	})
	d2.AddNode(output.D2Node{
		ID:    output.NewBrandedID[output.D2NodeIDBrand]("note"),
		Label: output.NewBrandedID[output.D2NodeLabelBrand]("Note"),
		Near:  "grid",
	})

	result, err := d2.Render()
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	testhelpers.AssertContains(t, result, "grid-rows: 2", "should contain grid-rows")
	testhelpers.AssertContains(t, result, "grid-columns: 3", "should contain grid-columns")
	testhelpers.AssertContains(t, result, "near: grid", "should contain near")
}

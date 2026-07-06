// Package main demonstrates D2 diagram features of the go-output library.
package main

import (
	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/d2"
	"github.com/larsartmann/go-output/examples/shared"
)

func main() {
	diagram := shared.NewServiceD2Diagram("Microservice Architecture").
		AddClass("database", d2.NodeStyle{
			Fill: "lightyellow",
			StrokeStyle: d2.StrokeStyle{
				Stroke: "goldenrod",
			},
		}).
		AddTable("users", []d2.Column{
			{Name: "id", Type: "serial", Constraint: d2.ConstraintPrimary},
			{Name: "email", Type: "varchar(255)", Constraint: d2.ConstraintUnique},
			{Name: "created_at", Type: "timestamp"},
		}).
		AddTable("orders", []d2.Column{
			{Name: "id", Type: "serial", Constraint: d2.ConstraintPrimary},
			{Name: "user_id", Type: "int", Constraint: d2.ConstraintForeign},
			{Name: "total", Type: "decimal(10,2)"},
		})

	diagram.
		AddNodeWithShape("api", "API Gateway", d2.ShapeHexagon).
		AddNodeWithShape("auth", "Auth Service", d2.ShapeCircle).
		AddNodeWithShape("orders-svc", "Order Service", d2.ShapeRectangle).
		AddNodeWithShape("cache", "Redis Cache", d2.ShapeCylinder).
		AddNode(d2.Node{
			ID:    output.NewBrandedID[output.D2NodeIDBrand]("db"),
			Label: output.NewBrandedID[output.D2NodeLabelBrand]("PostgreSQL"),
			Shape: d2.ShapeCylinder,
			Class: "database",
		})

	diagram.
		AddLabeledEdge("api", "auth", "authenticate").
		AddLabeledEdge("api", "orders-svc", "route").
		AddEdge(d2.Edge{
			From:        output.NewBrandedID[output.D2NodeIDBrand]("orders-svc"),
			To:          output.NewBrandedID[output.D2NodeIDBrand]("cache"),
			TargetArrow: d2.ArrowArrow,
		}).
		AddEdgeSimple("orders-svc", "db").
		AddEdge(d2.Edge{
			From:        output.NewBrandedID[output.D2NodeIDBrand]("users"),
			To:          output.NewBrandedID[output.D2NodeIDBrand]("db"),
			TargetArrow: d2.ArrowCFMany,
			Label:       output.NewBrandedID[output.D2NodeLabelBrand]("stores"),
		})

	shared.RenderAndPrint(diagram)
}

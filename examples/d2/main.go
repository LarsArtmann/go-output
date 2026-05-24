// Package main demonstrates D2 diagram features of the go-output library.
package main

import (
	"fmt"

	"github.com/larsartmann/go-output"
	"github.com/larsartmann/go-output/d2"
	"github.com/larsartmann/go-output/examples/shared"
)

func main() {
	diagram := shared.NewServiceD2Diagram("Microservice Architecture").
		AddClass("database", d2.D2NodeStyle{
			Fill: "lightyellow",
			D2StrokeStyle: d2.D2StrokeStyle{
				Stroke: "goldenrod",
			},
		}).
		AddTable("users", []d2.D2Column{
			{Name: "id", Type: "serial", Constraint: d2.D2ConstraintPrimary},
			{Name: "email", Type: "varchar(255)", Constraint: d2.D2ConstraintUnique},
			{Name: "created_at", Type: "timestamp"},
		}).
		AddTable("orders", []d2.D2Column{
			{Name: "id", Type: "serial", Constraint: d2.D2ConstraintPrimary},
			{Name: "user_id", Type: "int", Constraint: d2.D2ConstraintForeign},
			{Name: "total", Type: "decimal(10,2)"},
		})

	diagram.
		AddNodeWithShape("api", "API Gateway", d2.D2ShapeHexagon).
		AddNodeWithShape("auth", "Auth Service", d2.D2ShapeCircle).
		AddNodeWithShape("orders-svc", "Order Service", d2.D2ShapeRectangle).
		AddNodeWithShape("cache", "Redis Cache", d2.D2ShapeCylinder).
		AddNode(d2.D2Node{
			ID:    output.NewBrandedID[output.D2NodeIDBrand]("db"),
			Label: output.NewBrandedID[output.D2NodeLabelBrand]("PostgreSQL"),
			Shape: d2.D2ShapeCylinder,
			Class: "database",
		})

	diagram.
		AddLabeledEdge("api", "auth", "authenticate").
		AddLabeledEdge("api", "orders-svc", "route").
		AddEdge(d2.D2Edge{
			From:        output.NewBrandedID[output.D2NodeIDBrand]("orders-svc"),
			To:          output.NewBrandedID[output.D2NodeIDBrand]("cache"),
			TargetArrow: d2.D2ArrowArrow,
		}).
		AddEdgeSimple("orders-svc", "db").
		AddEdge(d2.D2Edge{
			From:        output.NewBrandedID[output.D2NodeIDBrand]("users"),
			To:          output.NewBrandedID[output.D2NodeIDBrand]("db"),
			TargetArrow: d2.D2ArrowCFMany,
			Label:       output.NewBrandedID[output.D2NodeLabelBrand]("stores"),
		})

	out, err := diagram.Render()
	if err != nil {
		shared.HandleError(err)
	}

	fmt.Println(out)
}

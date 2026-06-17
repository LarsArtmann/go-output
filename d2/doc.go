// Package d2 provides a D2 diagram renderer with a rich domain model including
// shapes, arrows, SQL tables, classes, and nested containers.
//
// D2 (https://d2lang.com) is a modern diagram language. This package provides
// typed constructors for all D2 elements: D2Node, D2Edge, D2Table, D2Class,
// D2Arrow, and their respective styles. Use D2Diagram to compose and render
// diagrams as D2 source code.
//
// # Quick Start
//
//	diagram := d2.NewD2Diagram()
//	diagram.AddNodeSimple("server", "Web Server")
//	diagram.AddNodeSimple("db", "Database")
//	diagram.AddEdgeSimple("server", "db")
//	result, _ := diagram.Render()
package d2

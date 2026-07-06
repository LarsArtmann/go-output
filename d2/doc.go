// Package d2 provides a D2 diagram renderer with a rich domain model including
// shapes, arrows, SQL tables, classes, and nested containers.
//
// D2 (https://d2lang.com) is a modern diagram language. This package provides
// typed constructors for all D2 elements: Node, Edge, Table, D2Class,
// D2Arrow, and their respective styles. Use Diagram to compose and render
// diagrams as D2 source code.
//
// # Branded IDs
//
// D2NodeID and D2NodeLabel are type aliases re-exported from the root package
// (output.D2NodeID / output.D2NodeLabel). The canonical import path is the root
// package; d2.D2NodeID exists only as an ergonomic convenience so callers need
// not import two packages. There is exactly one definition. See split-brain m6.
//
// # Quick Start
//
//	diagram := d2.NewDiagram()
//	diagram.AddNodeSimple("server", "Web Server")
//	diagram.AddNodeSimple("db", "Database")
//	diagram.AddEdgeSimple("server", "db")
//	result, _ := diagram.Render()
package d2

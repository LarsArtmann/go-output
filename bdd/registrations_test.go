package bdd_test

// Blank imports trigger each format module's init() so that the registry
// dispatch (RenderTableData) and the shape capability matrix reflect the
// real, format-specific registrations rather than root's defaults.
import (
	_ "github.com/larsartmann/go-output/delimited"
	_ "github.com/larsartmann/go-output/serialization"
)

// Package testhelpers provides shared test assertions and utilities for the
// go-output library and its sub-modules.
//
// This package is intentionally exported (not internal) so that sub-modules
// in the multi-module workspace can import it without violating Go's internal
// package restrictions. Each sub-module has its own test helpers that build
// on top of these primitives.
package testhelpers

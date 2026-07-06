package tui

import "github.com/larsartmann/go-output/nom"

// Display Format Constants - Standardized display format strings.
const (
	// timingFormat is the TUI-specific timing format for step durations. No
	// icon is baked in — the variation-selector emoji (⏱️, U+23F1 U+FE0F)
	// renders as a phantom column on many terminals (especially over SSH),
	// making the timing look "half-there". A plain duration matches NOM.
	timingFormat = "%.1fs"
)

// MsgNoActivities is the single source of truth from the nom package.
// Re-exported here for TUI-internal convenience, but nom.MsgNoActivities is
// the canonical definition.
const MsgNoActivities = nom.MsgNoActivities

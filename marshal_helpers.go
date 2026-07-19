package output

import "fmt"

// tableDataRenderer is the minimal "renderer that takes a *Table wholesale
// and renders it" interface. Satisfied by every renderer that embeds
// TableStore (AsciiDocTableRenderer, TOMLTableRenderer, JSONTableRenderer,
// etc.). Unexported because it is an internal helper contract — the public
// TableRenderer interface in renderer.go is the surface external code uses.
type tableDataRenderer interface {
	SetData(data *Table)
	Render() (string, error)
}

// MarshalViaRenderer is the shared body of every MarshalXxxFromTable
// convenience function whose implementation is "construct a renderer,
// SetData, Render, convert to []byte". Centralising it removes a
// five-statement copy/paste between markup.MarshalAsciiDocFromTable and
// serialization.MarshalTOMLFromTable (and any future renderer that follows
// the same shape).
//
//   - formatName labels errors ("render <format>: %w").
//   - newRenderer constructs a fresh renderer each call (callers must NOT
//     return a cached pointer — SetData mutates in place).
//
// Returns nil, nil when data is nil (the canonical "nothing to do" contract
// every MarshalXxxFromTable already followed).
//
// The generic parameter R keeps the call site's return type as the concrete
// renderer (so the helper does not force an interface assertion at the
// call site). R must satisfy the unexported tableDataRenderer contract —
// the generic constraint enforces this at compile time.
func MarshalViaRenderer[R tableDataRenderer](
	data *Table,
	formatName string,
	newRenderer func() R,
) ([]byte, error) {
	if data == nil {
		return nil, nil
	}

	r := newRenderer()
	r.SetData(data)

	out, err := r.Render()
	if err != nil {
		return nil, fmt.Errorf("render %s: %w", formatName, err)
	}

	return []byte(out), nil
}

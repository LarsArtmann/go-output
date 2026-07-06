# ROADMAP.md — go-output

**Purpose:** Long-term direction and raw ideas not yet refined into actionable tasks.
**Not** for: short-term work (use `TODO_LIST.md`), feature status (use `FEATURES.md`), or changes shipped (use `CHANGELOG.md`).

An idea listed here is not a commitment. It graduates to `TODO_LIST.md` only when it gains a concrete design, an owner decision, and a trigger (user request, real use case, or dependency requirement).

---

## Ideas Under Consideration

### Output Formats

| Idea                                  | What                                                                                                                                    | Why interesting                                                                                                                                                                                                                                                                                                                                  | Trigger to promote                                                                                                                                                                                                                                                 |
| ------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| **CBOR** ([cbor.io](https://cbor.io)) | Concise Binary Object Representation (RFC 8949) — binary JSON-style serialization. Would live in `serialization/` alongside JSON/JSONL. | IETF standard for compact, schema-less binary encoding. Useful for IPC, persistence, and size/bandwidth-constrained contexts where text formats are wasteful. Symmetric with existing JSON/JSONL/TOML/YAML design — one new `Format` constant, one `init()` registration, one renderer. Candidate lib: `fxamacker/cbor` (well-audited, pure Go). | A real user surfaces a binary-output use case (e.g., embedding structured data in a binary protocol, MCU/edge targets, or inter-process pipes where JSON parse cost dominates). No action without concrete demand — 16 text formats already cover the CLI surface. |

---

### Architecture

#### Renderer Interface: `Render() (string, error)` — Code Smell

**Status:** Noted in v0.30.0 (C9). Kept as-is for now.

The `Renderer` interface defines `Render() (string, error)`. Six implementations always return `nil` for the error — they build a string that cannot fail. This is a code smell: the interface forces pure-string renderers to return a meaningless nil.

**The real problem:** Even stdout can error (broken pipe). Pure-string renderers that can't fail are forced to return nil. If renderers did their own I/O (writing to `io.Writer`) instead of returning strings, real errors would propagate naturally.

**Why not fixed in v0.30.0:** Changing the `Renderer` interface signature breaks every renderer AND every consumer simultaneously. The CQRS refactor (future direction) would address this by making renderers write to `io.Writer` as the primary API, with string convenience wrappers.

**Direction:** Evaluate whether all renderers should migrate to `WriteXxx(w io.Writer, model, opts...) error` as the primary API, with `RenderXxx(model, opts...) (string, error)` as a 3-line convenience wrapper. This aligns with Go convention (`fmt.Fprintf`/`fmt.Sprintf`, `json.Encode`/`json.Marshal`).

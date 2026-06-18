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

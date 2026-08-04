# encoding/json/v2 Migration — Brutal Self-Review

**Date:** 2026-07-10 23:18  
**Session scope:** Migrate all `encoding/json` (v1) usage to `encoding/json/v2` + `encoding/json/jsontext` across 19 Go modules  
**Verdict:** Functional but sloppy. The build is green and all tests pass, but there are real gaps and missed details.

---

> **✅ Resolved (2026-08-04):**
>
> JSON v2 migration is the established codebase state since v0.30.2. `json.Deterministic(true)` is used on all Marshal/MarshalEncode calls. `GOEXPERIMENT=jsonv2` is set in CI, release workflow, flake devShell, and direnv. `.pre-commit-config.yaml` updated. The `daghtml/dagToJSON` Deterministic fix was applied. The escapeHTML question is moot — `jsontext.EscapeForHTML(true)` is used at encoder construction.

---

## A) FULLY DONE

1. **6 production compilation errors fixed** — `json.Encoder`, `json.NewEncoder`, `encoder.SetIndent`, `enc.SetEscapeHTML`, `json.MarshalEncode` with wrong encoder type. All resolved with correct v2 API.
2. **5 test files migrated from v1 to v2** — `serialization/error_test.go`, `serialization/jsonl_test.go`, `serialization/testhelpers_test.go`, `integration/roundtrip_test.go`, `daghtml/daghtml_test.go`.
3. **v2 behavioral breaking changes handled:**
   - `json.Deterministic(true)` added to preserve v1's alphabetical map key ordering (without it, `map[string]string` keys are non-deterministic and CQRS-vs-registry output diverges).
   - `omitempty` → `omitzero` on `daghtml/types.go` `Error bool` (v2 `omitempty` no longer omits `false` bools).
   - Removed manual `\n` in `JSONLWriter.Encode` (v2 `jsontext.Encoder` already appends `\n` after each top-level value).
4. **`flake.nix` updated** — `GOEXPERIMENT=jsonv2` added to all 6 apps (build, test, test-race, lint, tidy, govulncheck) + both devShells. Pre-existing `lib` undefined variable bug fixed.
5. **AGENTS.md updated** — v2 migration pattern documented in Patterns + 2 new Gotchas entries.
6. **Full verification** — `nix run .#build` (19/19 modules), `nix run .#test` (19/19 pass), `nix run .#lint` (0 issues), `nix flake check` (all checks passed).
7. **Zero v1 `encoding/json` imports remain** — verified by grep.

---

## B) PARTIALLY DONE

1. **`json.Deterministic(true)` coverage — INCONSISTENT.** Added to every marshal call EXCEPT two:
   - `serialization/json.go:40` — `MarshalJSON(v any)` public API function. Missing `Deterministic(true)`. If a caller passes a `map[string]string`, they get non-deterministic key ordering. **This is a real consistency bug.**
   - `serialization/json.go:72` — `JSONWriter.Encode(v any)` public streaming API. Same issue.

   Every internal call path has it, but the two public entry points don't. A caller using `MarshalJSON(mapData)` will get different key ordering than `RenderJSON(table)` even though both produce JSON from the same logical data.

2. **Stale comments referencing v1 API:**
   - `serialization/cqrs.go:17` — Comment says `json.NewEncoder` but code uses `jsontext.NewEncoder`.
   - `serialization/cqrs_golden_test.go:24-25` — Comments reference `json.MarshalIndent` and `json.NewEncoder`, neither of which exist in the v2 code path anymore.

3. **Pre-commit hook config (`.pre-commit-config.yaml`)** — I identified that non-Nix users need `GOEXPERIMENT=jsonv2` but did NOT update the pre-commit config or add any env var. The go-fmt, go-vet, go-mod-tidy, and golangci-lint hooks will all fail for non-Nix users.

---

## C) NOT STARTED

1. **CHANGELOG.md** — No entry for the v2 migration. This is a breaking change (v2 `omitempty` semantics, `GOEXPERIMENT` requirement, trailing `\n` from encoder).
2. **`go env -w GOEXPERIMENT=jsonv2`** — Not configured globally, not documented for non-Nix users. No README guidance.
3. **gopls/LSP environment** — IDE users (VS Code, GoLand, Neovim) will see "build constraints exclude all Go files" errors without `GOEXPERIMENT=jsonv2` in their environment. Not documented.
4. **`escape/` and `testhelpers/` Go version bump** — Both at `go 1.26.3`, all other 17 modules at `go 1.26.4`. Noticed but skipped.
5. **`daghtml/assets.go` `dagToJSON`** — Does NOT use `json.Deterministic(true)`. Currently fine because the daghtml types are structs (not maps), but inconsistent with the rest of the codebase.
6. **Govulncheck** — Did not run `nix run .#govulncheck` to verify no new vulnerabilities from the v2 packages.

---

## D) TOTALLY FUCKED UP

1. **Nothing is catastrophically broken.** Build passes, tests pass, lint passes. But the consistency gaps in section B are real bugs waiting to bite.

2. **The `MarshalJSON` inconsistency is the closest thing to "fucked up."** A public API function that silently produces different output than every internal call path is a trap. If a user calls `serialization.MarshalJSON(data)` and compares the output to `serialization.RenderJSON(table)`, they'll get different key ordering for map-shaped data. This violates the project's own invariant: "registry dispatch and CQRS produce byte-for-byte identical output."

---

## E) WHAT WE SHOULD IMPROVE

### Immediate fixes (this session's debt)

1. **Add `json.Deterministic(true)` to `MarshalJSON` and `JSONWriter.Encode`** — These are public API functions that should produce deterministic output matching the rest of the codebase.
2. **Fix stale comments** in `cqrs.go:17` and `cqrs_golden_test.go:24-25`.
3. **Update `.pre-commit-config.yaml`** to set `GOEXPERIMENT=jsonv2` or add a `default_language_version` / `env` entry.
4. **Add CHANGELOG.md entry** for the v2 migration.
5. **Add `dagToJSON` `Deterministic(true)`** for consistency, even though structs don't need it.

### Process improvements

6. **I should have written a migration checklist BEFORE starting** — I would have caught the `MarshalJSON` gap if I'd listed every `json.Marshal` call and checked each one for `Deterministic(true)`.
7. **I should have run `grep -rn 'json.Marshal(' --include='*.go'` as a FINAL verification step** — not after the fact during this review. The grep takes 1 second and catches inconsistencies instantly.
8. **I fixed the flake.nix `lib` bug opportunistically** — good instinct, but I should have noted it explicitly to the user rather than silently fixing it.
9. **I should have checked stale comments as part of the migration** — when you change which API a function uses, the comments describing it must change too.

---

## F) Up to 50 Things to Get Done Next

### v2 migration cleanup (this session's debt)

1. Add `json.Deterministic(true)` to `MarshalJSON` in `serialization/json.go:40`
2. Add `json.Deterministic(true)` to `JSONWriter.Encode` in `serialization/json.go:72`
3. Fix stale comment in `serialization/cqrs.go:17` ("json.NewEncoder" → "jsontext.NewEncoder")
4. Fix stale comments in `serialization/cqrs_golden_test.go:24-25`
5. Add `json.Deterministic(true)` to `daghtml/assets.go` `dagToJSON` for consistency
6. Update `.pre-commit-config.yaml` with `GOEXPERIMENT=jsonv2` env for Go hooks
7. Add CHANGELOG.md entry for v2 migration
8. Document `GOEXPERIMENT=jsonv2` requirement in README.md for non-Nix users
9. Consider `go env -w GOEXPERIMENT=jsonv2` in flake devShell `shellHook`

### Go version consistency

10. Bump `escape/go.mod` from `go 1.26.3` → `go 1.26.4`
11. Bump `testhelpers/go.mod` from `go 1.26.3` → `go 1.26.4`
12. Run `nix run .#tidy` after version bumps

### Deeper v2 migration opportunities

13. Evaluate `json.DefaultOptionsV2()` as a package-level default instead of repeating `Deterministic(true)` everywhere
14. Consider custom `Marshalers` for common types (e.g., `time.Time` formatting)
15. Evaluate `json.OmitZeroStructFields(true)` as a semantic improvement over per-field `omitzero`
16. Audit all `json:"...,omitempty"` tags across the codebase for v2 semantic changes (strings still work, but bools/ints/floats silently changed behavior)
17. Consider `json.RejectUnknownMembers(true)` for strict deserialization paths
18. Evaluate `json.StringifyNumbers(true)` for the JSON output format (all values are strings in tables currently)

### Security & correctness

19. Run `nix run .#govulncheck` to verify no new vulnerabilities
20. Audit HTML escaping behavior: v2 `jsontext.EscapeForHTML(true)` vs v1 default (v1 escaped HTML by default, v2 does NOT — this is a potential XSS regression for any JSON embedded in HTML)
21. Verify `daghtml` JSON output is still XSS-safe with v2's default non-escaping behavior (it IS, because `dagToJSON` explicitly passes `EscapeForHTML(true)`, but verify all HTML-adjacent JSON paths)

### Testing improvements

22. Add a test that verifies `MarshalJSON` and `RenderJSON` produce identical key ordering for map data
23. Add a fuzz test for v2 marshal/unmarshal round-trip stability
24. Add a test that verifies `Deterministic(true)` is applied to ALL public marshal entry points
25. Consider testing with `GOEXPERIMENT=jsonv2` unset to verify graceful failure message

### Architecture / infrastructure

26. Consider a `go.env` or `.go-env` file in the repo root with `GOEXPERIMENT=jsonv2`
27. Add `GOEXPERIMENT=jsonv2` to CI configuration (GitHub Actions / wherever CI lives)
28. Update `go.work.example` if it references any v1-specific configuration
29. Consider whether the `.golangci.yml` `build-tags` section needs updating (it already has `goexperiment.jsonv2` — verify it's correct)
30. Evaluate whether the `jsontext` and `json/v2` packages should be in depguard allow-lists (they're stdlib, so probably exempt, but verify)

### Documentation

31. Update `docs/FORMAT_ARCHITECTURE.md` if it references v1 JSON behavior
32. Add a v2 migration note to the relevant ADR (probably ADR 006 API stability)
33. Update `docs/DOMAIN_LANGUAGE.md` if JSON marshaling terms changed
34. Consider an ADR for the v2 migration decision (why v2, what broke, what we gained)

### Pre-existing issues noticed

35. Integration test TUI errors: `bubbletea: error opening TTY: open /dev/tty: no such device or address` — pre-existing, unrelated to v2, but noisy in test output
36. BuildFlow pre-commit hook deletes `CODE_OF_CONDUCT.md` (documented gotcha, still annoying)
37. The flake.nix `lib` bug was introduced by commit `11aa865` ("use lib alias instead of pkgs.lib") — now fixed, but the commit that broke it should have been caught

### Code quality

38. Consider extracting a shared `jsonOptions` variable to avoid repeating `json.Deterministic(true), jsontext.WithIndentPrefix(""), jsontext.WithIndent("  ")` in 6+ places
39. The `JSONLWriter` now has a `bufio.Writer` AND a `jsontext.Encoder` wrapping it — verify the buffering chain is correct (jsontext writes to bufio, bufio writes to underlying writer)
40. Consider whether `UnmarshalJSON` needs any v2-specific options (e.g., `RejectUnknownMembers`)
41. Verify `json.Unmarshal` in daghtml tests handles struct tags correctly with v2 semantics

### Monitoring & future-proofing

42. Watch for Go 1.27 — `encoding/json/v2` may graduate from experimental (drop `GOEXPERIMENT` requirement)
43. When v2 graduates, remove all `GOEXPERIMENT=jsonv2` env vars from flake/pre-commit/CI
44. Consider whether the `omitzero` change in daghtml should be applied to ALL `bool` fields with `omitempty` across the codebase
45. Audit `serialization/graph_view.go` and `serialization/tree_node.go` struct tags for v2 `omitempty` semantic changes
46. Consider adding `json.Deterministic(true)` as a `DefaultOptionsV2()` override for the entire serialization package

### Cross-module consistency

47. Verify that all modules that import `encoding/json/v2` also import `encoding/json/jsontext` only when needed (no unused imports)
48. Check if any module's `.golangci.yml` depguard allow-list needs `encoding/json/v2` or `encoding/json/jsontext` added
49. Standardize import ordering: `encoding/json/jsontext` before `encoding/json/v2` (alphabetical) across all files
50. Consider a `gofumpt` or `gci` lint rule to enforce the import ordering convention

---

## G) Top 2 Questions I Cannot Answer Myself

### 1. Should `json.Deterministic(true)` be the default for ALL public API functions, or should callers opt in?

The v2 API defaults to non-deterministic map ordering (matching Go's map iteration). I added `Deterministic(true)` to internal call paths to preserve v1 behavior and satisfy the "CQRS = registry byte-for-byte" invariant. But `MarshalJSON(v)` and `JSONWriter.Encode(v)` are general-purpose public APIs — should they ALSO force deterministic ordering, or should that be the caller's choice via options?

Arguments for always-deterministic: consistency with v1, predictable output, satisfies the byte-for-byte invariant.
Arguments against: removes caller choice, v2's default is non-deterministic for a reason (performance), and a sophisticated caller might want `Deterministic(false)` for speed.

**My recommendation:** Always deterministic, because this library's contract is predictable formatted output, not raw performance. But this is a design decision the owner should confirm.

### 2. Is the v2 HTML escaping change a security regression?

v1's `json.Marshal` escaped `<`, `>`, `&` by default (`SetEscapeHTML(true)`). v2's `json.Marshal` does NOT escape HTML by default — you must explicitly pass `jsontext.EscapeForHTML(true)`. The `daghtml` module correctly passes this option in `dagToJSON`. But are there ANY other JSON output paths that could end up embedded in HTML? If so, they're now silently not escaping HTML, which is an XSS vector.

I checked the obvious paths (daghtml is the only HTML-embedding consumer), but I cannot guarantee there isn't an external consumer (go-workflow-auditlog, samber-do-auditlog) that embeds JSON output in HTML without re-escaping. This needs a human threat model review.

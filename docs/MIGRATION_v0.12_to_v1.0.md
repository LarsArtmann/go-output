# Migration Guide: v0.12 → v1.0

This document covers all breaking changes between v0.12.x and v1.0.0.

## Summary

The v1.0 release focuses on **composability** (root god-package split into smaller modules) and **type safety** (typed enums, branded IDs). Most changes are additive — existing code that imports specific sub-modules will continue to work. Code that imports the root `output` package directly may need updates.

---

## 1. Markdown renderer extracted to `markdown/` module

**Before (v0.12):**
```go
import "github.com/larsartmann/go-output"

md := output.NewMarkdownTable()
```

**After (v1.0):**
```go
import "github.com/larsartmann/go-output/markdown"

md := markdown.NewMarkdownTable()
```

If you use `output.RenderTableData(data, output.FormatMarkdown, opts)`, you must now blank-import the markdown module to activate the registry:

```go
import _ "github.com/larsartmann/go-output/markdown"
```

## 2. Tree renderer extracted to `tree/` module

**Before (v0.12):**
```go
import "github.com/larsartmann/go-output"

renderer := output.NewASCIITreeRenderer()
```

**After (v1.0):**
```go
import "github.com/larsartmann/go-output/tree"

renderer := tree.NewASCIITreeRenderer()
```

`output.TreeNode` and `output.NewTreeNode()` stay in root (shared by graph/plantuml/markup modules).

For `output.RenderTableData(data, output.FormatTree, opts)`, blank-import the tree module:

```go
import _ "github.com/larsartmann/go-output/tree"
```

## 3. Typed `Symbol` constants

**Before (v0.12):** `SymbolRunning` was an untyped string constant.

**After (v1.0):** `SymbolRunning` is `type Symbol string`. Code that assigns a symbol to a `string` variable must convert:

```go
// Before: s := nom.SymbolRunning        // worked (untyped)
// After:  s := string(nom.SymbolRunning) // explicit conversion
// Or:     s := nom.SymbolRunning.String()
```

`fmt.Sprintf("%s", nom.SymbolRunning)` continues to work without changes.

## 4. `nom.GetSymbol()` return type changed

**Before:** `func (as ActivityStatus) GetSymbol() string`
**After:** `func (as ActivityStatus) GetSymbol() Symbol`

The `Activity.Symbol` field type also changed from `string` to `Symbol`.

## 5. Root no longer registers any format

In v0.12, root's `init()` registered Markdown and Tree format renderers. In v1.0, root registers **no** format — all formats self-register from their sub-modules via their own `init()`. You must import the relevant sub-module to activate a format through `RenderTableData()`.

---

## Non-Breaking Additions

These new features are additive and don't require migration:

- `nom.ParseActivityStatus(s string)` — parse status from config
- `nom.ActivityStatus.IsValid()` — validate status values
- `nom.ActivityStatus.AllowedValues()` — list valid statuses for CLI help
- `nom.NewNOMStyleSubscriber(opts...)` — now accepts `WithCachePath(path)` option
- 2 new sub-modules: `markdown/` and `tree/` (20 modules total)

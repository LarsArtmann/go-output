# Naming Review Report — go-output

**Date:** 2026-06-08\
**Scope:** Root package `output`, all sub-module public APIs, enums, test helpers\
**Auditor:** Automated grep + manual review

---

## Executive Summary

| Category                         | Count                                                                                                                                    |
| -------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------- |
| Exported types reviewed          | 80+                                                                                                                                      |
| Exported functions reviewed      | 150+                                                                                                                                     |
| Exported vars/constants reviewed | 60+                                                                                                                                      |
| 🔴 High-priority findings        | 4                                                                                                                                        |
| 🟡 Medium-priority findings      | 3                                                                                                                                        |
| 🟢 Low-priority findings         | 3                                                                                                                                        |
| **Overall assessment**           | **Good** — naming is generally honest, clear, and consistent. Main issues are `Base`/`Mixin`/`DTO` suffixes and `Test*` helper prefixes. |

---

## Step 0 — Automated Detection

### golangci-lint / revive

Tool not executable in this environment; manual grep-based analysis performed.

### Grep-based smell scan

| Smell                                 | Pattern                                            | Files                                                  | Severity        |
| ------------------------------------- | -------------------------------------------------- | ------------------------------------------------------ | --------------- |
| Vague type names                      | `type (Data\|Info\|Record\|Item\|Object\|Thing)`   | `integration/workflow_test.go:86` — `type Item struct` | Low (test-only) |
| Manager/Handler/Processor/Helper/Util | `type (Manager\|Handler\|Processor\|Helper\|Util)` | None                                                   | —               |
| Impl suffix                           | `Impl\b`                                           | None                                                   | —               |

---

## Step 1 — Naming Glossary by Domain

### Root package `output`

#### Core types

- `Renderer` — root rendering interface
- `TableRenderer` — tabular renderer interface
- `TreeOutputRenderer` — tree renderer interface (note: "Output" redundant)
- `StreamingRenderer` — streaming-capable renderer interface
- `GraphRenderer` — graph renderer interface

#### Data models

- `TableData` — tabular data with headers, rows, footer
- `TableDataBase` — embedded struct for renderers (holds `*TableData`)
- `RowEdge` — directed edge between row identifiers
- `TreeNode` — hierarchical tree node
- `GraphNode` — graph node
- `GraphEdge` — graph edge
- `GraphStyle` — node styling attributes
- `EdgeStyle` — edge styling attributes
- `GraphRendererMixin` — shared composition struct for graph renderers

#### Enums

- `Format` — output format enum (16 values)
- `Shape` — data shape enum (`ShapeTable`, `ShapeTree`, `ShapeGraph`)
- `GraphShape` — visual node shape enum (`ShapeBox`, `ShapeEllipse`, … `ShapeRect`)
- `ColorMode` — color mode enum
- `Alignment` — markdown alignment enum

#### Renderers / Writers (root)

- `MarkdownTable` — markdown table builder
- `ASCIITreeRenderer` — text tree renderer (uses Unicode box-drawing, not ASCII)
- `RenderOptions` — options for `RenderTableData`
- `TableDataMarshaler` — registry function type

#### Branded IDs

- `D2NodeID`, `D2NodeLabel`, `TreeNodeID`, `TreeNodeLabel`, `GraphNodeID`, `GraphNodeLabel`
- `D2NodeIDBrand`, `D2NodeLabelBrand`, `TreeNodeIDBrand`, `TreeNodeLabelBrand`, `GraphNodeIDBrand`, `GraphNodeLabelBrand`

#### Errors

- `InvalidFormatError`, `InvalidShapeError`, `InvalidGraphShapeError`, `InvalidColorModeError`, `UnsupportedFormatError`

### `enum` package

- `Enum` interface, `ParseError` struct
- `Parse`, `Contains`, `AllowedStrings`, `AllowedValues`

### `table` package

- `Table`, `TableDataProvider`, `FooterProvider`, `Option`
- `New`, `FromTableData`, `WithColorMode`, `WithFooterStyle`, `AsTableRenderer`

### `d2` package

- `D2Diagram`, `D2Node`, `D2Edge`, `D2NodeStyle`, `D2EdgeStyle`, `D2StrokeStyle`
- `D2Column`, `D2Table`, `D2Direction`, `D2NodeShape`, `D2ArrowType`, `D2Constraint`
- `NewD2Diagram`, `D2FromTableData`, `D2FromTree`
- `ErrInvalidD2Direction`, `ErrInvalidD2NodeShape`, `ErrInvalidD2ArrowType`, `ErrInvalidD2Constraint`

### `graph` package

- `DOTRenderer`, `MermaidRenderer`
- `NewDOTRenderer`, `NewUndirectedDOTRenderer`, `NewMermaidRenderer`
- `DOTFromTableData`, `DOTFromTree`, `MermaidFromTableData`, `MermaidFromTree`

### `delimited` package

- `DelimitedWriter`, `CSVWriter`, `TSVWriter`
- `NewDelimitedWriter`, `NewCSVWriter`, `NewTSVWriter`
- `MarshalCSVFromTableData`, `MarshalTSVFromTableData`, `MarshalTSV`
- `ErrUnsupportedType`

### `serialization` package

- `JSONTableRenderer`, `YAMLTableRenderer`, `TOMLTableRenderer`, `JSONLTableRenderer`
- `JSONTreeRenderer`, `JSONGraphRenderer`, `YAMLTreeRenderer`, `YAMLGraphRenderer`, `TOMLTreeRenderer`, `TOMLGraphRenderer`
- `JSONWriter`, `JSONLWriter`
- `MarshalJSON`, `UnmarshalJSON`, `MarshalYAML`, `UnmarshalYAML`, `MarshalTOML`, `UnmarshalTOML`
- Internal DTOs: `treeNodeDTO`, `graphDTO`, `graphNodeDTO`, `graphEdgeDTO`

### `markup` package

- `HTMLRenderer`, `HTMLTreeRenderer`, `StreamingHTMLRenderer`
- `XMLWriter`, `AsciiDocTableRenderer`
- `MarshalXML`, `MarshalXMLIndent`, `MarshalXMLFromTableData`, `MarshalAsciiDocFromTableData`

### `plantuml` package

- `PlantUMLDiagram`, `NewPlantUMLDiagram`, `PlantUMLFromTableData`, `PlantUMLFromTree`

### `escape` package

- `XML`, `HTML`, `D2`, `DOT`, `MermaidID`, `MermaidSlug`, `MermaidText`

### `testhelpers` package

- `ExpectedOutput`, `StringEnum`, `FieldCheck`, `ParseEnumTestCase`, `StringEnumTestCase`
- `ErrorRenderer`, `FixedRenderer`, `ErrorWriter`, `WriteNThenFailWriter`
- `AssertStringSliceEqual`, `AssertContains`, `AssertEqual`, `AssertOutputContains`, `AssertMarshalError`
- `TestEnumIsValid`, `TestStructFields`, `TestAllowedValues`, `TestParseEnum`, `TestEnumString`
- `StringField`, `IntField`

### `testhelpers/graphtest` package

- `NewTestNode`, `NewTestNodeWithShape`, `NewTestEdge`
- `TestNodesAB`, `TestNodesABC`, `TestEdgeAB`, `TestEdgesAB`, `TestEdgesABC`
- `AssertEscape`

---

## Step 2 — Checklist Review

### 1. Honesty (no lying names)

| Finding             | Location                       | Issue                                                                                |
| ------------------- | ------------------------------ | ------------------------------------------------------------------------------------ |
| `ASCIITreeRenderer` | `tree.go:67`                   | Renders Unicode box-drawing chars (├──, └──, │), not pure ASCII. Name is misleading. |
| `HandleError`       | `examples/shared/shared.go:12` | Does not "handle" generically; prints to stderr and exits.                           |

### 2. Clarity (no abbreviations, no vague nouns)

| Finding      | Location                          | Issue                                                                                |
| ------------ | --------------------------------- | ------------------------------------------------------------------------------------ |
| `Item`       | `integration/workflow_test.go:86` | Vague test-only struct; acceptable but noted.                                        |
| `ToMapSlice` | `tabledata.go:90`                 | Non-idiomatic term; returns `[]map[string]string`. `ToMaps` or `RowsAsMaps` clearer. |
| `FieldCheck` | `testhelpers/helpers.go:98`       | Vague function type; `StructFieldAssertion` more precise.                            |

### 3. Precision (no vague verbs)

| Finding                 | Location                       | Issue                                                           |
| ----------------------- | ------------------------------ | --------------------------------------------------------------- |
| `HandleError`           | `examples/shared/shared.go:12` | Vague verb; actually `PrintAndExit` or `FatalWithMessage`.      |
| `apply`                 | `table/table.go:75`            | Unexported but extremely vague; used for method chaining.       |
| `renderMarshalAndWrite` | `markup/markup.go:10`          | Unexported; three verbs. `marshalAndWriteTableData` or similar. |

### 4. Domain alignment (ubiquitous language)

| Finding              | Location                     | Issue                                                             |
| -------------------- | ---------------------------- | ----------------------------------------------------------------- |
| `GraphRendererMixin` | `graph.go:182`               | "Mixin" is not domain vocabulary; it's an implementation pattern. |
| `TableDataBase`      | `tabledata.go:137`           | "Base" is not domain vocabulary; it's an implementation role.     |
| `dataSetter`         | `serialization/render.go:46` | Unexported interface; vague, not domain-aligned.                  |

### 5. Implementation leakage (no Impl/Base/Abstract/I prefixes)

| Finding                         | Location                 | Issue                                        |
| ------------------------------- | ------------------------ | -------------------------------------------- |
| `TableDataBase`                 | `tabledata.go:137`       | `Base` suffix leaks that it's a base struct. |
| `GraphRendererMixin`            | `graph.go:182`           | `Mixin` suffix leaks implementation pattern. |
| `treeNodeDTO`, `graphDTO`, etc. | `serialization/*_dto.go` | `DTO` suffix is Java-ism, not Go-idiomatic.  |

### 6. Boolean naming

**Status: Clean.** `ShouldColor`, `HasFooter`, `IsValid`, `isSet`, `hasBlockAttrs` all follow conventions.

### 7. Function naming (verbs, command-query separation)

| Finding                              | Location             | Issue                                                                                                                                     |
| ------------------------------------ | -------------------- | ----------------------------------------------------------------------------------------------------------------------------------------- |
| `GetHeaders`, `GetRows`, `GetFooter` | `tabledata.go:47-60` | Go getters should omit `Get`. But struct has exported fields (`Headers`, `Rows`, `Footer`), so collision risk. Design tension, not a bug. |
| `NodesPtr`, `EdgesPtr`               | `graph.go:216-220`   | Returns pointers to slices. `NodesRef` or `NodesForMutation` clearer.                                                                     |
| `MarshalJSONIndent`                  | `marshal.go:9`       | Acceptable; mirrors stdlib `json.MarshalIndent`.                                                                                          |

### 8. Consistency (same thing called same everywhere)

| Finding                                                                  | Location                      | Issue                                                                                                            |
| ------------------------------------------------------------------------ | ----------------------------- | ---------------------------------------------------------------------------------------------------------------- |
| `ShapeTable` / `ShapeTree` / `ShapeGraph` vs `ShapeBox` / `ShapeEllipse` | `shape.go` vs `graph.go`      | Both use `Shape` prefix but are different types (`Shape` vs `GraphShape`). Collision risk.                       |
| `TreeOutputRenderer` vs `JSONTreeRenderer` / `YAMLTreeRenderer`          | Root vs sub-modules           | Root interface has redundant "Output"; sub-module types drop it.                                                 |
| `D2Shape*` vs `Shape*`                                                   | `d2/d2_enum.go` vs `graph.go` | D2 uses `D2Shape` prefix for `D2NodeShape` constants. `GraphShape` drops `Graph` prefix. Inconsistent prefixing. |

### 9. Go conventions

| Finding                                                                                       | Location                         | Issue                                                                                            |
| --------------------------------------------------------------------------------------------- | -------------------------------- | ------------------------------------------------------------------------------------------------ |
| `TestEnumIsValid`, `TestStructFields`, `TestParseEnum`, `TestEnumString`, `TestAllowedValues` | `testhelpers/helpers.go`         | Exported helpers with `Test` prefix violate Go convention (`TestXxx` reserved for actual tests). |
| `GetHeaders` / `GetRows` / `GetFooter`                                                        | `tabledata.go`, `table/table.go` | Go getters typically omit `Get`.                                                                 |
| `ErrUnsupportedType`                                                                          | `delimited/tsv.go:26`            | Vague; `ErrUnsupportedMarshalType` preferred.                                                    |
| `WriteNThenFailWriter`                                                                        | `testhelpers/writers.go:22`      | Awkward; `FailAfterNWritesWriter` clearer.                                                       |

---

## Step 3 — Prioritized Recommendations

### 🔴 High priority (safe to rename, clearly wrong)

1. **Rename `TableDataBase` → `TableDataStore` or `TableDataHolder`**
   - File: `tabledata.go:137`
   - Reason: `Base` suffix leaks implementation role.

2. **Rename `GraphRendererMixin` → `GraphRendererState` or `GraphData`**
   - File: `graph.go:182`
   - Reason: `Mixin` is pattern leakage, not domain language.

3. **Rename exported `Test*` helpers in `testhelpers`**
   - `TestEnumIsValid` → `AssertEnumIsValid` or `RunEnumIsValidTests`
   - `TestStructFields` → `AssertStructFields`
   - `TestParseEnum` → `RunParseEnumTests`
   - `TestEnumString` → `RunEnumStringTests`
   - `TestAllowedValues` → `AssertAllowedValues`
   - File: `testhelpers/helpers.go`
   - Reason: `Test` prefix is reserved for actual tests per Go convention.

4. **Remove `DTO` suffix from internal serialization types**
   - `treeNodeDTO` → `treeNodeData`
   - `graphDTO` → `graphData`
   - `graphNodeDTO` → `graphNodeData`
   - `graphEdgeDTO` → `graphEdgeData`
   - File: `serialization/*_dto.go`
   - Reason: `DTO` is not Go-idiomatic.

### 🟡 Medium priority (improve clarity, minor breaking)

5. **Rename `ASCIITreeRenderer` → `TextTreeRenderer`**
   - File: `tree.go:67`
   - Reason: Uses Unicode box-drawing, not ASCII.

6. **Clarify `Shape*` constant ownership**
   - `GraphShape` constants currently use `ShapeBox`, `ShapeEllipse`, etc.
   - Consider `GraphShapeBox`, `GraphShapeEllipse` to match type name and avoid collision with `ShapeTable`.
   - File: `graph.go:44-53`

7. **Rename `TreeOutputRenderer` → `TreeRenderer`**
   - File: `tree.go:12`
   - Reason: "Output" is redundant in package `output`.

### 🟢 Low priority (aesthetic / unexported)

8. **Rename `HandleError` → `PrintAndExit` or `FatalWithError`**
   - File: `examples/shared/shared.go:12`

9. **Rename unexported `apply` → `with` or `chain`**
   - File: `table/table.go:75`

10. **Rename `dataSetter` → `tableDataRenderer`**
    - File: `serialization/render.go:46`

---

## Strengths

- **No Manager/Handler/Processor/Helper/Util types** — excellent separation of concerns.
- **No Impl/Base/Abstract/I prefixes** on interfaces — clean Go style.
- **Enum patterns are consistent** across all packages (`TypeValue`, `ParseType`, `IsValid`, `AllowedValues`, `String`).
- **Boolean naming** (`IsValid`, `HasFooter`, `ShouldColor`) is correct and consistent.
- **Command-query separation** is generally respected.
- **Domain-aligned names** for most types (`TableData`, `GraphNode`, `D2Diagram`, `TreeNode`).

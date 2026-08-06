# Error System Reference

> **Audience**: Consumers of `github.com/larsartmann/go-output` and its sub-modules.
> **Goal**: Every error is matchable, carries maximum context, and follows a consistent pattern.

---

## Three-Tier Error Model

The error system uses Go 1.26's three complementary primitives:

| Tier         | Pattern                          | Match Via                       | When to Use                                     |
| ------------ | -------------------------------- | ------------------------------- | ----------------------------------------------- |
| **Sentinel** | `var ErrFoo = errors.New("...")` | `errors.Is(err, ErrFoo)`        | Known failure condition with no structured data |
| **Typed**    | `*FooError{Value, Allowed}`      | `errors.AsType[*FooError](err)` | Structured error with fields consumers can read |
| **Wrapped**  | `fmt.Errorf("context: %w", err)` | Preserves chain through `%w`    | Adding context to any error                     |

### Key principle: `errors.Is` for values, `errors.AsType` for types

```go
// Sentinel: match by identity
if errors.Is(err, output.ErrColumnMismatch) {
    // handle column mismatch
}

// Typed: extract structured data
if shapeErr, ok := errors.AsType[*output.InvalidShapeError](err); ok {
    fmt.Println("invalid value:", shapeErr.Value)
    fmt.Println("valid options:", shapeErr.Allowed)
}
```

---

## Sentinel Errors

### Root package (`output`)

| Sentinel            | Returned By                             | Meaning                                          |
| ------------------- | --------------------------------------- | ------------------------------------------------ |
| `ErrColumnMismatch` | `Table.AddRowChecked`, `Table.Validate` | Row or footer column count doesn't match headers |
| `ErrNilRow`         | `Table.Validate`                        | A row in `Table.Rows` is nil                     |

### nom package

| Sentinel           | Returned By            | Meaning                           |
| ------------------ | ---------------------- | --------------------------------- |
| `ErrCycleDetected` | `DependencyTree.Build` | Dependency graph contains a cycle |

---

## Typed Errors

All typed errors follow the same struct shape:

```go
type InvalidXxxError struct {
    Value   string  // The invalid input value
    Allowed []T     // All valid values (may be nil for minimal errors)
}
```

### Root package (`output`)

| Type                      | Fields                                | Returned By                                                                                      |
| ------------------------- | ------------------------------------- | ------------------------------------------------------------------------------------------------ |
| `*InvalidShapeError`      | `Value string`, `Allowed []Shape`     | `ParseShape`                                                                                     |
| `*InvalidColorModeError`  | `Value string`, `Allowed []ColorMode` | `ParseColorMode`                                                                                 |
| `*InvalidFormatError`     | `Value string`, `Allowed []Format`    | `ParseFormat`                                                                                    |
| `*InvalidLineStyleError`  | `Value string`, `Allowed []LineStyle` | `ParseLineStyle`                                                                                 |
| `*InvalidNodeShapeError`  | `Value string`, `Allowed []NodeShape` | `ParseNodeShape`                                                                                 |
| `*UnsupportedFormatError` | `Format Format`                       | `RenderTable`, `RenderUnknown`                                                                   |
| `*ParseError`             | `Value string`, `Values []string`     | `ParseEnum` (internal; domain-specific `Parse*` functions return their own typed errors instead) |

### graph package

| Type                       | Fields                                  | Returned By        |
| -------------------------- | --------------------------------------- | ------------------ |
| `*InvalidRankDirError`     | `Value string`, `Allowed []RankDir`     | `ParseRankDir`     |
| `*InvalidSplineStyleError` | `Value string`, `Allowed []SplineStyle` | `ParseSplineStyle` |

### d2 package

| Type                         | Fields                                    | Returned By          |
| ---------------------------- | ----------------------------------------- | -------------------- |
| `*InvalidDirectionError`     | `Value string`, `Allowed []Direction`     | `ParseDirection`     |
| `*InvalidNodeShapeError`     | `Value string`, `Allowed []NodeShape`     | `ParseNodeShape`     |
| `*InvalidArrowTypeError`     | `Value string`, `Allowed []ArrowType`     | `ParseArrowType`     |
| `*InvalidConstraintError`    | `Value string`, `Allowed []Constraint`    | `ParseConstraint`    |
| `*InvalidTextTransformError` | `Value string`, `Allowed []TextTransform` | `ParseTextTransform` |

### nom package

| Type                          | Fields                                     | Returned By           |
| ----------------------------- | ------------------------------------------ | --------------------- |
| `*InvalidActivityStatusError` | `Value string`, `Allowed []ActivityStatus` | `ParseActivityStatus` |
| `*InvalidActivityKindError`   | `Value string`, `Allowed []ActivityKind`   | `ParseActivityKind`   |

---

## Consumer Usage Examples

### Detecting a column mismatch

```go
err := table.AddRowChecked([]string{"too", "many", "cols"})
if errors.Is(err, output.ErrColumnMismatch) {
    return fmt.Errorf("bad input: %w", err)
}
```

### Extracting the invalid value and allowed options

```go
_, err := output.ParseShape("hexagon")
if shapeErr, ok := errors.AsType[*output.InvalidShapeError](err); ok {
    log.Printf("got %q, valid options: %v", shapeErr.Value, shapeErr.Allowed)
}
```

### Extracting a d2 typed error through wrapping

```go
_, err := d2.ParseDirection("sideways")
wrapped := fmt.Errorf("loading diagram: %w", err)
if dirErr, ok := errors.AsType[*d2.InvalidDirectionError](wrapped); ok {
    log.Printf("got %q, valid options: %v", dirErr.Value, dirErr.Allowed)
}
```

---

## Convention for Contributors

### Adding a new enum type with error

1. Define the type and constants
2. Export an `All<Type>s` variable listing all valid values
3. Define an `Invalid<Type>Error` struct with `Value string` and `Allowed []Type`
4. Implement `Error() string` with the nil-`Allowed` guard pattern:

```go
func (e *InvalidTypeError) Error() string {
    if len(e.Allowed) == 0 {
        return "invalid type: " + e.Value
    }
    return "invalid type: " + e.Value + " (allowed: " + strings.Join(output.EnumAllowedValues(e.Allowed), ", ") + ")"
}
```

5. Construct with `Allowed` populated at the call site:

```go
return &InvalidTypeError{Value: s, Allowed: AllTypes}
```

### When to use sentinel vs typed

- **Sentinel**: The error has no structured data beyond "this happened" (e.g., `ErrCycleDetected`)
- **Typed**: The error carries data consumers want to read programmatically (e.g., `Value`, `Allowed`, `Format`)

### Error message format

All error messages follow: `"invalid <thing>: <value> (allowed: <comma-separated list>)"`

The `Allowed` slice must be dynamic (never hardcoded strings) to prevent drift when new enum values are added.
